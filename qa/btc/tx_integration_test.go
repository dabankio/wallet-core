// +build integration

package btc

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcjson"

	clibtcjson "github.com/lomocoin/btccli/btcjson"
	"github.com/lomocoin/wallet-core/core/btc"

	"github.com/lomocoin/btccli"
	"github.com/stretchr/testify/require"
)

func TestSimpleTX(t *testing.T) {
	rq := require.New(t)

	// killbitcoind, err := btccli.BitcoindWithOptions(btccli.StartOptions{NewTmpDir: true})
	cli, killBitcoind, err := btccli.RunBitcoind(&btccli.RunOptions{NewTmpDir: true})
	rq.Nil(err)
	defer killBitcoind()

	// 实施过程：
	// 使用账号a0,a1进行测试
	// 生成101个块给a0(产生utxo)
	// a0转账给a1
	// a1检查余额(utxo)

	{
		// err = cli.Importprivkey(clibtcjson.ImportPrivKeyCmd{PrivKey: a0.Privkey})
		// rq.Nil(err)
		for _, ad := range []string{a0.Address, a1.Address} {
			err = cli.Importaddress(clibtcjson.ImportAddressCmd{Address: ad})
			rq.Nil(err)
		}
	}

	var utxo clibtcjson.ListUnspentResult
	{
		_, err = cli.Generatetoaddress(101, a0.Address, nil)
		rq.Nil(err)

		cliUnspents, err := cli.Listunspent(1, 999, []string{a0.Address}, nil, nil)
		rq.Nil(err)
		rq.Equal(1, len(cliUnspents), "")

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
			// decAddr, err := btcutil.DecodeAddress(a0.Address, &chaincfg.RegressionNetParams)
			// rq.Nil(err)
			// _ = decAddr

			info, err := cli.GetAddressInfo(a0.Address)
			rq.Nil(err)
			fmt.Printf("info: %#v\n", info)

			// pkScript, err := txscript.PayToAddrScript(decAddr)
			// rq.Nil(err)

			// b, err := hex.DecodeString(a0.Pubkey)
			// rq.Nil(err)
			// pk, err := btcutil.NewAddressPubKey(b, &chaincfg.RegressionNetParams)
			// rq.Nil(err)

			// pk, err := btcec.ParsePubKey(b, btcec.S256())
			// pk := (*btcec.PublicKey)(&key.PublicKey).
			// 	SerializeUncompressed()
			// tmpAddress, err := btcutil.NewAddressPubKeyHash(
			// 	btcutil.Hash160(b), &chaincfg.RegressionNetParams)
			// rq.Nil(err)
			// fmt.Println("tmpAddress:", tmpAddress.EncodeAddress())

			// pkScript, err := txscript.PayToAddrScript(tmpAddress)
			// rq.Nil(err)
			// _ = pkScript

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

		{ //sdk sign
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

			rs, err := btc.Sign(tx, a0.Privkey, chainID)
			rq.Nil(err)
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
		rs, err := cli.Decoderawtransaction(clibtcjson.DecodeRawTransactionCmd{
			HexTx: signedHex,
		})
		rq.Nil(err)
		fmt.Println("decoded tx:", btccli.ToJSONIndent(rs))
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
