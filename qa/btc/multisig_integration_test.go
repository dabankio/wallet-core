// +build integration

package btc

import (
	"fmt"
	"github.com/dabankio/wallet-core/core/btc"
	"strings"
	"testing"

	clibtcjson "github.com/dabankio/btccli/btcjson"

	"github.com/dabankio/btccli"
	"github.com/stretchr/testify/require"
)

// 多重签名测试
func TestMultisig(t *testing.T) {
	rq := require.New(t)

	cli, killBitcoind, err := btccli.RunBitcoind(&btccli.RunOptions{NewTmpDir: true})
	rq.Nil(err)
	defer killBitcoind()

	// 导入a0 private Key, a1 a2 a3 address
	// 首先为a0生成 101 个块
	// 用a1 a2 a3 生成多签地址 (2-3)
	// 往多签地址转入btc
	// a1 a2 签名，转出到 a3
	// 查询a3 utxo, 应该不为0

	{ // import addresses
		err = cli.Importprivkey(clibtcjson.ImportPrivKeyCmd{
			PrivKey: a0.Privkey,
		})
		rq.Nil(err)

		for _, add := range []string{a1.Address, a2.Address, a3.Address} {
			err = cli.Importaddress(clibtcjson.ImportAddressCmd{Address: add})
			rq.Nil(err)
		}
	}

	{
		_, err = cli.Generatetoaddress(101, a0.Address, nil)
		rq.Nil(err)
	}

	var multisigAddress, redeemScript string
	{ // create multisig address,and import to bitcoind
		rs, err := btc.NewMultiSigAddress(2, btc.ChainRegtest, strings.Join([]string{a1.Pubkey, a2.Pubkey, a3.Pubkey}, ","))
		rq.Nil(err)
		arr := strings.Split(rs, ",")
		rq.Len(arr, 2, "")
		multisigAddress, redeemScript = arr[0], arr[1]

		err = cli.Importaddress(clibtcjson.ImportAddressCmd{
			Address: multisigAddress,
		})
		rq.Nil(err)
	}

	{ //send to multisig address for next step
		txid, err := cli.Sendtoaddress(&clibtcjson.SendToAddressCmd{
			Address: multisigAddress, Amount: 23.3,
		})
		rq.Nil(err)
		rq.NotContains(txid, "error", "")
	}

	unspents, err := cli.Listunspent(0, 999, []string{multisigAddress}, nil, nil)
	rq.Nil(err)
	rq.Len(unspents, 1, "")

	utxo := unspents[0]
	fmt.Printf("%#v", utxo)

	var signedHex string
	transferAmount := 2.3
	chainID := btc.ChainRegtest
	{ //SDK 的使用可以参考这里
		var tx *btc.BTCTransaction
		unspent := new(btc.BTCUnspent) //java: new btc.BTCUnspent()
		{                              // build tx
			unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, redeemScript)

			amount, err := btc.NewBTCAmount(transferAmount)
			rq.Nil(err)

			toAddressA1, err := btc.NewBTCAddressFromString(a1.Address, chainID)
			rq.Nil(err)

			outputAmount := btc.BTCOutputAmount{} //java: new btc.BTCOutputAmount()
			outputAmount.Add(toAddressA1, amount)

			feeRate := int64(80)

			changeAddressMultisig, err := btc.NewBTCAddressFromString(multisigAddress, chainID)
			rq.Nil(err)

			tx, err = btc.NewBTCTransaction(unspent, &outputAmount, changeAddressMultisig, feeRate, chainID)
			rq.Nil(err)
		}

		{ //sdk sign
			// member a1 sign
			signRs, err := btc.SignTransaction(tx, a1.Privkey, chainID)
			rq.Nil(err)
			rq.True(signRs.Changed, "")
			rq.False(signRs.Complete, "")
			signedHex = signRs.Hex

			// member a3 sign
			signRs, err = btc.SignRawTransactionWithKey(signedHex, a3.Privkey, unspent, chainID)
			rq.Nil(err)
			rq.True(signRs.Changed)
			rq.True(signRs.Complete)
			signedHex = signRs.Hex
		}
	}

	{ // relay tx
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

	{ // validate utxo for receiver
		unspentsForA1, err := cli.Listunspent(0, 999, []string{a1.Address}, clibtcjson.Bool(true), nil)
		rq.Nil(err)
		rq.True(len(unspentsForA1) > 0, "No utxo of a1 fund!")
		fmt.Println("utxo for a1", btccli.ToJSONIndent(unspentsForA1))
		rq.Equal(transferAmount, unspentsForA1[0].Amount, "Wrong amount")
	}

}
