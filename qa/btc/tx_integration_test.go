// +build integration

package btc

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcjson"

	devtools "github.com/dabankio/devtools4chains"

	"github.com/dabankio/wallet-core/core/btc"

	"github.com/stretchr/testify/require"
)

// 简单转账签名测试
func TestSimpleTX(t *testing.T) {
	rq := require.New(t)

	// killbitcoind, err := btccli.BitcoindWithOptions(btccli.StartOptions{NewTmpDir: true})
	// cli, killBitcoind, err := btccli.RunBitcoind(&btccli.RunOptions{NewTmpDir: true})
	// rq.Nil(err)
	// defer killBitcoind()

	killFunc, nodeInfo, err := devtools.DockerRunBitcoin(devtools.DockerRunOptions{
		AutoRemove: true,
	})
	rq.NoError(err)
	defer killFunc()

	rpcInfo := devtools.RPCInfo{
		Host:     fmt.Sprintf("http://127.0.0.1:%d", nodeInfo.RPCPort),
		User:     nodeInfo.RPCUser,
		Password: nodeInfo.RPCPwd,
	}

	// 实施过程：
	// 使用账号a0,a1进行测试
	// 生成101个块给a0(产生utxo)
	// a0转账给a1
	// a1检查余额(utxo)

	{
		// err = cli.Importprivkey(clibtcjson.ImportPrivKeyCmd{PrivKey: a0.Privkey})
		// rq.Nil(err)
		for _, ad := range []string{a0.Address, a1.Address} {
			err = devtools.RPCCallJSON(rpcInfo, "importaddress", []interface{}{ad}, nil)
			rq.Nil(err)
		}
	}

	type ListUnspentResult struct {
		TxID string `json:"txid"`
	}
	var utxo clibtcjson.ListUnspentResult
	{
		err = devtools.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{101, a0.Address}, nil)
		rq.Nil(err)

		var unspents []ListUnspentResult
		err = devtools.RPCCallJSON(rpcInfo, "listunspent", []interface{}{1, 999, a0.Address}, unspents)
		rq.Nil(err)
		rq.Equal(1, len(unspents), "")

		utxo = cliUnspents[0]
		rq.Equal(50.0, utxo.Amount, "coinbase amount should be 50")
		fmt.Println(utxo)
	}

	var signedHex string
	transferAmount := 2.3
	chainID := btc.ChainRegtest

	{ //SDK 的使用可以参考这里
		var tx *btc.BTCTransaction
		{ // build tx
			unspent := new(btc.BTCUnspent) //java: new btc.BTCUnspent()
			unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, utxo.RedeemScript)

			amount, err := btc.NewBTCAmount(transferAmount)
			rq.Nil(err)

			toAddressA1, err := btc.NewBTCAddressFromString(a1.Address, chainID)
			rq.Nil(err)

			outputAmount := btc.BTCOutputAmount{} //java: new btc.BTCOutputAmount()
			outputAmount.Add(toAddressA1, amount)

			feeRate := int64(80)

			changeAddressA0, err := btc.NewBTCAddressFromString(a0.Address, chainID)
			rq.Nil(err)

			tx, err = btc.NewBTCTransaction(unspent, &outputAmount, changeAddressA0, feeRate, chainID)
			rq.Nil(err)
		}

		{ //sdk sign,这个不能在移动端使用，btc.New 未导出
			// coin, _ := btc.New(nil, chainID)
			// notSignedHex, _ := tx.Encode()
			// cmd, err := tx.EncodeToSignCmd()
			// rq.Nil(err)
			// signedHex, err = coin.Sign(cmd, a0.Privkey)
			// rq.Nil(err)
			// rq.NotEqual(notSignedHex, signedHex, "raw tx not changed")
		}

		{ // sdk sign 2
			notSignedHex, _ := tx.Encode()

			rs, err := btc.SignTransaction(tx, a0.Privkey, chainID)
			rq.Nil(err)
			rq.True(rs.Complete, "")
			rq.True(rs.Changed, "")
			signedHex = rs.Hex
			rq.NotEqual(notSignedHex, signedHex, "raw tx not changed")
		}

		{ // cli sign
			// cmd, err := tx.EncodeToSignCmd()
			// rq.Nil(err)

			// var msg clibtcjson.SignRawTransactionCmd
			// b, err := hex.DecodeString(cmd)
			// rq.Nil(err)

			// err = json.Unmarshal(b, &msg)
			// rq.Nil(err)
			// msg.PrivKeys = &[]string{a0.Privkey}
			// if msg.Flags == nil {
			// 	var flagALL = "ALL"
			// 	msg.Flags = &flagALL
			// }
			// msg.Prevtxs = []clibtcjson.PreviousDependentTxOutput{{
			// 	TxID:         utxo.TxID,
			// 	Vout:         utxo.Vout,
			// 	ScriptPubKey: utxo.ScriptPubKey,
			// 	RedeemScript: utxo.RedeemScript,
			// 	Amount:       utxo.Amount,
			// }}
			// rs, err := cli.Signrawtransactionwithkey(msg)
			// rq.Nil(err)
			// signedHex = rs.Hex
		}

		{ // 完全用cli签名
			// cmd := clibtcjson.CreateRawTransactionCmd{
			// 	Inputs: []clibtcjson.TransactionInput{
			// 		clibtcjson.TransactionInput{
			// 			Txid: utxo.TxID,
			// 			Vout: utxo.Vout,
			// 		},
			// 	},
			// 	Outputs: []map[string]interface{}{
			// 		map[string]interface{}{
			// 			a1.Address: 45.0,
			// 		},
			// 		map[string]interface{}{
			// 			a0.Address: utxo.Amount - 45.0 - 0.001,
			// 		},
			// 	},
			// }
			// rawHex, err := cli.Createrawtransaction(cmd)
			// rq.Nil(err)

			// fmt.Println("Then decode rawHex")
			// // _, err = CliDecoderawtransaction(clibtcjson.DecodeRawTransactionCmd{
			// // 	HexTx: rawHex,
			// // })
			// // testtool.FailOnFlag(t, err != nil, "Failed to decode raw tx", err)

			// keys := []string{a0.Privkey}
			// signRes, err := cli.Signrawtransactionwithkey(clibtcjson.SignRawTransactionCmd{
			// 	RawTx:    rawHex,
			// 	PrivKeys: &keys,
			// 	Prevtxs: []clibtcjson.PreviousDependentTxOutput{
			// 		clibtcjson.PreviousDependentTxOutput{
			// 			TxID:         utxo.TxID,
			// 			Vout:         utxo.Vout,
			// 			ScriptPubKey: utxo.ScriptPubKey,
			// 			Amount:       utxo.Amount,
			// 		},
			// 	},
			// })
			// rq.Nil(err)
			// signedHex = signRes.Hex
		}

	}

	{ // decode
		// rs, err := cli.Decoderawtransaction(clibtcjson.DecodeRawTransactionCmd{
		// 	HexTx: signedHex,
		// })
		// rq.Nil(err)
		// fmt.Println("decoded tx:", btccli.ToJSONIndent(rs))
	}

	{ // 广播交易

		txid, err := cli.Sendrawtransaction(clibtcjson.SendRawTransactionCmd{
			HexTx: signedHex,
		})
		rq.Nil(err)
		fmt.Println("txid", txid)
		rq.NotContains(txid, "error", "")
	}

	{ //generate 1 block
		_, err := cli.Generatetoaddress(1, a0.Address, nil)
		rq.Nil(err)
	}

	{ //查询a1 的utxo
		utxos, err := cli.Listunspent(0, 999, []string{a1.Address}, btcjson.Bool(true), nil)
		rq.Nil(err)
		rq.True(len(utxos) > 0, "No utxo of a1 fund!")
		fmt.Println("utxo for a1", btccli.ToJSONIndent(utxos))
		// rq.Equal(transferAmount, utxos[0].Amount, "Wrong amount")
	}
}
