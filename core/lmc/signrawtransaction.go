// Copyright (c) 2013-2015 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package lmc

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/lomocoin/HDWallet-Core/core/btc"
	"github.com/lomocoin/HDWallet-Core/core/lmc/ppcd/txscript"
	"github.com/lomocoin/HDWallet-Core/core/lmc/ppcd/wire"
	"github.com/pkg/errors"
)

// SignRawTransaction handles the signrawtransaction command.
func SignRawTransaction(cmd *btc.SignRawTransactionCmd, chainCfg *chaincfg.Params) (*btcjson.SignRawTransactionResult, error) {
	serializedTx, err := btc.DecodeHexStr(cmd.RawTx)
	if err != nil {
		return nil, err
	}
	msgTx := wire.NewMsgTx()
	err = msgTx.Deserialize(bytes.NewBuffer(serializedTx))
	if err != nil {
		e := errors.New("TX decode failed")
		return nil, e
	}

	// First we add the stuff we have been given.
	// TODO(oga) really we probably should look these up with btcd anyway
	// to make sure that they match the blockchain if present.
	inputs := make(map[wire.OutPoint][]byte)
	scripts := make(map[string][]byte)
	var cmdInputs []btc.RawTxInput
	if cmd.Inputs != nil {
		cmdInputs = *cmd.Inputs
	}
	for _, rti := range cmdInputs {
		inputSha, err := wire.NewShaHashFromStr(rti.Txid)
		if err != nil {
			return nil, err
		}

		script, err := btc.DecodeHexStr(rti.ScriptPubKey)
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
			redeemScript, err := btc.DecodeHexStr(rti.RedeemScript)
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
			Hash:  *inputSha,
			Index: rti.Vout,
		}] = script
	}

	// Now we go and look for any inputs that we were not provided by
	// querying btcd with getrawtransaction. We queue up a bunch of async
	// requests and will wait for replies after we have checked the rest of
	// the arguments.
	for _, txIn := range msgTx.TxIn {
		// Did we get this txin from the arguments?
		if _, ok := inputs[txIn.PreviousOutPoint]; ok {
			continue
		}
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

	// All args collected. Now we can sign all the inputs that we can.
	// `complete' denotes that we successfully signed all outputs and that
	// all scripts will run to completion. This is returned as part of the
	// reply.
	var signErrors []btcjson.SignRawTransactionError
	for i, txIn := range msgTx.TxIn {
		input, ok := inputs[txIn.PreviousOutPoint]
		if !ok {
			// failure to find previous is actually an error since
			// we failed above if we don't have all the inputs.
			return nil, errors.Errorf("%s:%d not found", txIn.PreviousOutPoint.Hash, txIn.PreviousOutPoint.Index)
		}

		// Set up our callbacks that we pass to txscript so it can
		// look up the appropriate keys and scripts by address.
		getKey := txscript.KeyClosure(func(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
			if len(keys) != 0 {
				wif, ok := keys[addr.EncodeAddress()]
				if !ok {
					return nil, false,
						errors.New("no key for address")
				}
				return wif.PrivKey, wif.CompressPubKey, nil
			}
			return nil, false, errors.New("no pk")
		})

		getScript := txscript.ScriptClosure(func(addr btcutil.Address) ([]byte, error) {
			// If keys were provided then we can only use the
			// scripts provided with our inputs, too.
			if len(keys) != 0 {
				script, ok := scripts[addr.EncodeAddress()]
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
		if (hashType&txscript.SigHashSingle) != txscript.SigHashSingle || i < len(msgTx.TxOut) {
			script, err := txscript.SignTxOutput(chainCfg, msgTx, i, input, hashType, getKey, getScript, txIn.SignatureScript)
			// Failure to sign isn't an error, it just means that
			// the tx isn't complete.
			if err != nil {
				signErrors = append(signErrors,
					btcjson.SignRawTransactionError{
						TxID:      txIn.PreviousOutPoint.Hash.String(),
						Vout:      txIn.PreviousOutPoint.Index,
						ScriptSig: hex.EncodeToString(txIn.SignatureScript),
						Sequence:  txIn.Sequence,
						Error:     err.Error(),
					})
				continue
			}
			txIn.SignatureScript = script
		}

		// Either it was already signed or we just signed it.
		// Find out if it is completely satisfied or still needs more.
		vm, err := txscript.NewEngine(input, msgTx, i, txscript.StandardVerifyFlags)
		if err == nil {
			err = vm.Execute()
		}
		if err != nil {
			signErrors = append(signErrors,
				btcjson.SignRawTransactionError{
					TxID:      txIn.PreviousOutPoint.Hash.String(),
					Vout:      txIn.PreviousOutPoint.Index,
					ScriptSig: hex.EncodeToString(txIn.SignatureScript),
					Sequence:  txIn.Sequence,
					Error:     err.Error(),
				})
		}
	}

	var buf bytes.Buffer
	buf.Grow(msgTx.SerializeSize())

	// All returned errors (not OOM, which panics) encounted during
	// bytes.Buffer writes are unexpected.
	if err = msgTx.Serialize(&buf); err != nil {
		return nil, err
	}

	return &btcjson.SignRawTransactionResult{
		Hex:      hex.EncodeToString(buf.Bytes()),
		Complete: len(signErrors) == 0,
		Errors:   signErrors,
	}, nil
}
