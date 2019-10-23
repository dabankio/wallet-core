package btc

import (
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/dabankio/wallet-core/core/btc/internal"
	"github.com/pkg/errors"
)

// SignTransaction signs tx with wif privateKey
// tx: transaction
// privateKeyWif: private key in wallet import format(wif)
// chainFlag: chain id
// return: {
//     changed: (bool, 签名后的内容是否发生变化，),
// 	   complete: (bool， 含不同地址的输入或多重签名情况下即使提供正确的私钥也可能存在签名未完成的情况),
//     hex: (string, 前后的rawTransactionHex)
//     errors: (string 可作为调试参考的错误信息)
// }
func SignTransaction(tx *BTCTransaction, privateKeyWif string, chainFlag int) (*SignRawTransactionResult, error) {
	data, err := tx.Encode()
	if err != nil {
		return nil, err
	}
	var inputs []btcjson.ListUnspentResult
	for _, ipt := range *tx.rawTxInput {
		inputs = append(inputs, btcjson.ListUnspentResult{
			TxID: ipt.Txid, Vout: ipt.Vout, ScriptPubKey: ipt.ScriptPubKey, RedeemScript: ipt.RedeemScript,
		})
	}
	return SignRawTransactionWithKey(data, privateKeyWif, &BTCUnspent{unspent: inputs}, chainFlag)
}

// SignRawTransactionWithKey Refer to https://bitcoin.org/en/developer-reference#signrawtransactionwithkey
// Diff from rpc:
//     sighashtype is set to "ALL"
//     single private key
func SignRawTransactionWithKey(hexstring string, privateKeyWif string, unspents *BTCUnspent, chainFlag int) (*SignRawTransactionResult, error) {
	var inputs []internal.RawTxInput
	for _, ipt := range unspents.unspent {
		inputs = append(inputs, internal.RawTxInput{
			Txid: ipt.TxID, Vout: ipt.Vout, ScriptPubKey: ipt.ScriptPubKey, RedeemScript: ipt.RedeemScript,
		})
	}
	cmd := internal.NewSignRawTransactionCmd(hexstring, &inputs, &[]string{privateKeyWif}, btcjson.String("ALL"))
	chainParam, err := internal.ChainFlag2ChainParams(chainFlag)
	if err != nil {
		return nil, err
	}
	result, err := internal.SignRawTransaction(cmd, chainParam)
	if err != nil {
		err = errors.Wrap(err, "btc.sign.signRawTransaction")
		return nil, err
	}

	return &SignRawTransactionResult{
		Changed:  result.Hex != hexstring,
		Hex:      result.Hex,
		Complete: result.Complete,
		Errors:   fmt.Sprintf("%v", result.Errors),
	}, nil
}
