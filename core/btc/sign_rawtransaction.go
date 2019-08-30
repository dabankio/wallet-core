package btc

import (
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/lomocoin/wallet-core/core/btc/internal"
	"github.com/pkg/errors"
)

// Sign signs raw tx with wif privateKey
func Sign(tx *BTCTransaction, privateKeyWif string, chainFlag int) (*SignRawTransactionResult, error) {
	data, err := tx.Encode()
	if err != nil {
		return nil, err
	}

	var inputs []internal.RawTxInput
	for _, ipt := range *tx.rawTxInput {
		inputs = append(inputs, internal.RawTxInput{
			Txid: ipt.Txid, Vout: ipt.Vout, ScriptPubKey: ipt.ScriptPubKey, RedeemScript: ipt.RedeemScript,
		})
	}
	cmd := internal.NewSignRawTransactionCmd(data, &inputs, &[]string{privateKeyWif}, btcjson.String("ALL"))
	chainParam, err := internal.ChainFlag2ChainParams(chainFlag)
	if err != nil {
		return nil, err
	}
	result, err := internal.SignRawTransaction(cmd, chainParam)
	if err != nil {
		err = errors.Wrap(err, "btc.sign.signRawTransaction")
		return nil, err
	}
	var errs string
	if len(result.Errors) > 0 {
		errs = fmt.Sprintf("%v", result.Errors)
	}

	return &SignRawTransactionResult{
		Hex:      result.Hex,
		Complete: result.Complete,
		Errors:   btcjson.String(errs),
	}, nil
}
