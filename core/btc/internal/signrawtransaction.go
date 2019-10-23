// Based on the https://github.com/btcsuite/btcwallet/blob/master/rpc/legacyrpc/methods.go
// todo: need more testing

package internal

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

// SignRawTransaction handles the signrawtransaction command.
func SignRawTransaction(cmd *SignRawTransactionCmd, chainCfg *chaincfg.Params) (*btcjson.SignRawTransactionResult, error) {
	serializedTx, err := DecodeHexStr(cmd.RawTx)
	if err != nil {
		return nil, err
	}
	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewBuffer(serializedTx))
	if err != nil {
		e := errors.New("TX decode failed")
		return nil, e
	}

	var hashType txscript.SigHashType
	switch *cmd.Flags {
	case "ALL":
		hashType = txscript.SigHashAll
	case "NONE":
		hashType = txscript.SigHashNone
	case "SINGLE":
		hashType = txscript.SigHashSingle
	case "ALL|ANYONECANPAY":
		hashType = txscript.SigHashAll | txscript.SigHashAnyOneCanPay
	case "NONE|ANYONECANPAY":
		hashType = txscript.SigHashNone | txscript.SigHashAnyOneCanPay
	case "SINGLE|ANYONECANPAY":
		hashType = txscript.SigHashSingle | txscript.SigHashAnyOneCanPay
	default:
		e := errors.New("Invalid sighash parameter")
		return nil, e
	}

	// TODO: really we probably should look these up with btcd anyway to
	// make sure that they match the blockchain if present.
	inputs := make(map[wire.OutPoint][]byte)
	scripts := make(map[string][]byte)
	var cmdInputs []RawTxInput
	if cmd.Inputs != nil {
		cmdInputs = *cmd.Inputs
	}
	for _, rti := range cmdInputs {
		inputHash, err := chainhash.NewHashFromStr(rti.Txid)
		if err != nil {
			return nil, err
		}

		script, err := DecodeHexStr(rti.ScriptPubKey)
		if err != nil {
			return nil, err
		}

		// redeemScript is only actually used iff the user provided
		// private keys. In which case, it is used to get the scripts
		// for signing. If the user did not provide keys then we always
		// get scripts from the wallet.
		// Empty strings are ok for this one and hex.DecodeString will
		// DTRT.
		if cmd.PrivKeys != nil && len(*cmd.PrivKeys) != 0 {
			redeemScript, err := DecodeHexStr(rti.RedeemScript)
			if err != nil {
				return nil, err
			}

			addr, err := btcutil.NewAddressScriptHash(redeemScript, chainCfg)
			if err != nil {
				return nil, err
			}
			scripts[addr.String()] = redeemScript
		}
		inputs[wire.OutPoint{
			Hash:  *inputHash,
			Index: rti.Vout,
		}] = script
	}

	// Parse list of private keys, if present. If there are any keys here
	// they are the keys that we may use for signing. If empty we will
	// use any keys known to us already.
	var keys map[string]*btcutil.WIF
	if cmd.PrivKeys != nil {
		keys = make(map[string]*btcutil.WIF)

		for _, key := range *cmd.PrivKeys {
			wif, err := btcutil.DecodeWIF(key)
			if err != nil {
				return nil, err
			}

			if !wif.IsForNet(chainCfg) {
				s := "key network doesn't match wallet's"
				return nil, errors.New(s)
			}

			addr, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), chainCfg)
			if err != nil {
				return nil, err
			}
			keys[addr.EncodeAddress()] = wif
		}
	}

	getScriptPubKey := func(txID string) (scriptPubKey []byte, err error) {
		for i := range *cmd.Inputs {
			if (*cmd.Inputs)[i].Txid == txID {
				return hex.DecodeString((*cmd.Inputs)[i].ScriptPubKey)
			}
		}
		err = errors.Errorf("not fund scriptPubKey: %s", txID)
		return
	}
	for i := range tx.TxIn {
		var scriptPubKey []byte
		scriptPubKey, err = getScriptPubKey(tx.TxIn[i].PreviousOutPoint.Hash.String())
		if err != nil {
			return nil, err
		}
		inputs[tx.TxIn[i].PreviousOutPoint] = scriptPubKey
	}

	// All args collected. Now we can sign all the inputs that we can.
	// `complete' denotes that we successfully signed all outputs and that
	// all scripts will run to completion. This is returned as part of the
	// reply.
	signErrs, err := signTransaction(&tx, hashType, inputs, keys, scripts, chainCfg)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Grow(tx.SerializeSize())

	// All returned errors (not OOM, which panics) encounted during
	// bytes.Buffer writes are unexpected.
	if err = tx.Serialize(&buf); err != nil {
		return nil, err
	}

	signErrors := make([]btcjson.SignRawTransactionError, 0, len(signErrs))
	for _, e := range signErrs {
		input := tx.TxIn[e.InputIndex]
		signErrors = append(signErrors, btcjson.SignRawTransactionError{
			TxID:      input.PreviousOutPoint.Hash.String(),
			Vout:      input.PreviousOutPoint.Index,
			ScriptSig: hex.EncodeToString(input.SignatureScript),
			Sequence:  input.Sequence,
			Error:     e.Error.Error(),
		})
	}

	return &btcjson.SignRawTransactionResult{
		Hex:      hex.EncodeToString(buf.Bytes()),
		Complete: len(signErrors) == 0,
		Errors:   signErrors,
	}, nil
}

// SignatureError records the underlying error when validating a transaction
// input signature.
type SignatureError struct {
	InputIndex uint32
	Error      error
}

// SignTransaction uses secrets of the wallet, as well as additional secrets
// passed in by the caller, to create and add input signatures to a transaction.
//
// Transaction input script validation is used to confirm that all signatures
// are valid.  For any invalid input, a SignatureError is added to the returns.
// The final error return is reserved for unexpected or fatal errors, such as
// being unable to determine a previous output script to redeem.
//
// The transaction pointed to by tx is modified by this function.
func signTransaction(tx *wire.MsgTx, hashType txscript.SigHashType,
	additionalPrevScripts map[wire.OutPoint][]byte,
	additionalKeysByAddress map[string]*btcutil.WIF,
	p2shRedeemScriptsByAddress map[string][]byte,
	chainCfg *chaincfg.Params) ([]SignatureError, error) {

	var signErrors []SignatureError
	for i, txIn := range tx.TxIn {
		prevOutScript, ok := additionalPrevScripts[txIn.PreviousOutPoint]
		if !ok {
			err := errors.Errorf("not fund prevOutScript: %s", txIn.PreviousOutPoint.Hash.String())
			return signErrors, err
		}

		// Set up our callbacks that we pass to txscript so it can
		// look up the appropriate keys and scripts by address.
		getKey := txscript.KeyClosure(func(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
			if len(additionalKeysByAddress) != 0 {
				addrStr := addr.EncodeAddress()
				wif, ok := additionalKeysByAddress[addrStr]
				if !ok {
					return nil, false, errors.New("no key for address")
				}
				return wif.PrivKey, wif.CompressPubKey, nil
			}
			return nil, false, errors.New("no pk")
		})
		getScript := txscript.ScriptClosure(func(addr btcutil.Address) ([]byte, error) {
			// If keys were provided then we can only use the
			// redeem scripts provided with our inputs, too.
			if len(additionalKeysByAddress) != 0 {
				addrStr := addr.EncodeAddress()
				script, ok := p2shRedeemScriptsByAddress[addrStr]
				if !ok {
					return nil, errors.New("no script for address")
				}
				return script, nil
			}
			return nil, errors.New("no key")
		})

		// SigHashSingle inputs can only be signed if there's a
		// corresponding output. However this could be already signed,
		// so we always verify the output.
		if (hashType&txscript.SigHashSingle) != txscript.SigHashSingle || i < len(tx.TxOut) {
			script, err := txscript.SignTxOutput(chainCfg, tx, i, prevOutScript, hashType, getKey, getScript, txIn.SignatureScript)
			// Failure to sign isn't an error, it just means that
			// the tx isn't complete.
			if err != nil {
				signErrors = append(signErrors, SignatureError{
					InputIndex: uint32(i),
					Error:      err,
				})
				continue
			}
			txIn.SignatureScript = script
		}

		// Either it was already signed or we just signed it.
		// Find out if it is completely satisfied or still needs more.
		vm, err := txscript.NewEngine(prevOutScript, tx, i, txscript.StandardVerifyFlags, nil, nil, 0)
		if err == nil {
			err = vm.Execute()
		}
		if err != nil {
			signErrors = append(signErrors, SignatureError{
				InputIndex: uint32(i),
				Error:      err,
			})
		}
	}
	return signErrors, nil
}
