// +build integration

package omni

import (
	"github.com/lomocoin/wallet-core/core/omni"
	"github.com/stretchr/testify/require"
	// "time"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lomocoin/omnicli"
	"github.com/lomocoin/omnicli/btcjson"

	"github.com/lomocoin/wallet-core/core/btc"
)

func TestSimpleSend(t *testing.T) {
	rq := require.New(t)

	cli, killomnicored, err := omnicli.RunOmnicored(&omnicli.RunOptions{NewTmpDir: true})
	rq.Nil(err)
	defer killomnicored()

	rq.Nil(importAddrs(cli, presetAddrs))
	a0, a1 := presetAddrs[0], presetAddrs[1]
	// time.Sleep(time.Second)

	{ // 生成多个块，获取utxo
		omnicli.NoPrintCmd(func() {
			_, err := cli.Generatetoaddress(107, a0.Address, nil)
			rq.False(err != nil, "Failed to generate to addresss", err)
		})
		// cli.Importaddress(btcjson.ImportAddressCmd{Address: "2MuNFG2wDexZUJseuxn9kXpqmHq4pBe9Fzi"}) //这里，不清楚为什么有时候generate后不能获取utxo，这样导入随便一个地址后就可以
	}

	propertyID := createToken(t, cli, &omnicli.OmniSenddissuancefixedCmd{
		Fromaddress: a0.Address,
		Ecosystem:   2, //2 fot test
		Typ:         1, // 1 for indivisible
		Previousid:  0, // 0 for new tokens
		Category:    "test_omniii",
		Subcategory: "unit_test",
		Name:        "FakeUSDT",
		Amount:      "10000",
	})

	var utxo btcjson.ListUnspentResult
	{
		var unspents []btcjson.ListUnspentResult
		unspents, err = cli.Listunspent(0, 999, []string{a0.Address})
		rq.Nil(err, "Failed to list unspent")
		rq.False(len(unspents) == 0, "no unspent find")

		utxo = unspents[0]
	}

	chainID := btc.ChainRegtest
	transferAmount := float64(233)
	{ // simple send from a0 to a1,
		unspent := new(btc.BTCUnspent)
		unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, utxo.RedeemScript)
		toAddr, err := btc.NewBTCAddressFromString(a1.Address, chainID)
		rq.Nil(err, "failed to create addr")
		changeAddr, err := btc.NewBTCAddressFromString(a0.Address, chainID)
		rq.Nil(err, "failed to create addr")
		feeRate := int64(80)
		btctx, err := omni.CreateTransactionForOmni(propertyID, false, unspent, toAddr, transferAmount, changeAddr, feeRate, chainID)
		rq.Nil(err, "Failed to crate btctx")

		toSignMsg, err := btctx.EncodeToSignCmd()
		rq.Nil(err, "failed to encode to sign")

		btcCoin, _ := btc.New(nil, chainID)
		signedRawHex, err := btcCoin.Sign(toSignMsg, a0.Privkey)
		assert.Nil(t, err)

		cli.DecodeAndPrintTX("signed tx", signedRawHex)
		// detx, _ := btc.DecodeRawTransaction(&btcdbtcjson.DecodeRawTransactionCmd{
		// 	HexTx: signedRawHex,
		// }, &chaincfg.RegressionNetParams)
		// b, _ := json.MarshalIndent(&detx, "", " ")
		// fmt.Println("to send tx", string(b))

		// 广播交易
		txid, err := cli.Sendrawtransaction(btcjson.SendRawTransactionCmd{
			HexTx: signedRawHex,
		})
		rq.Nil(err, "failed to send tx")
		fmt.Println("broadcasted txid", txid)
	}

	{ // 生成一个块确认代币转账
		_, err = cli.Generatetoaddress(1, a0.Address, nil)
		rq.Nil(err, "Failed to generate to address")
	}

	{ // 确认代币转账成功
		bal, err := cli.OmniGetbalance(a1.Address, propertyID)
		rq.Nil(err, "Failed to get omni balance")
		expectedBal := strconv.FormatFloat(transferAmount, 'f', 0, 32)
		rq.False(bal.Balance != expectedBal, "wrong balance, not ", expectedBal)
	}
}

func createToken(t *testing.T, cli *omnicli.Cli, cmd *omnicli.OmniSenddissuancefixedCmd) (propertyID int) {
	rq := require.New(t)
	{ // create token
		txHash, err := cli.OmniSendissuancefixed(cmd)
		rq.Nil(err, "Failed to create omni coin")
		{ //生成几个块，确认刚才的交易
			_, err = cli.Generatetoaddress(1, cmd.Fromaddress, nil)
			rq.Nil(err, "Failed to generate to address")
		}
		tx, err := cli.OmniGettransaction(txHash)
		rq.Nil(err, "Failed to get tx")
		propertyID = tx.Propertyid

		rq.False(propertyID == 0, "Got property id error", propertyID)
	}

	{ // 代币创建完成后查询代币持有人的余额，应该等于总的发行量
		fmt.Println("-------then balance of new created property-----")
		bal, err := cli.OmniGetbalance(cmd.Fromaddress, propertyID)
		rq.Nil(err, "Failed to get balance of owner")
		rq.False(bal.Balance != cmd.Amount, "余额不符合预期")
	}
	return
}
