package wallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dabankio/devtools4chains"
	"github.com/dabankio/wallet-core/core/btc"
	"github.com/dabankio/wallet-core/core/eth/internalized/testtool"
	"github.com/dabankio/wallet-core/qa/omni"
	"github.com/dabankio/wallet-core/wallet"
	"github.com/stretchr/testify/require"
)

func testBTCPubkSign(t *testing.T, w *wallet.Wallet, c ctx) {
	rq := require.New(t)

	killFunc, bitcoinInfo, err := devtools4chains.DockerRunBitcoin(devtools4chains.DockerRunOptions{
		AutoRemove: true, Image: &bbcImage,
	})
	require.NoError(t, err)
	t.Cleanup(killFunc)

	rpcInfo := devtools4chains.RPCInfo{
		Host:     fmt.Sprintf("http://127.0.0.1:%d", bitcoinInfo.RPCPort),
		User:     bitcoinInfo.RPCUser,
		Password: bitcoinInfo.RPCPwd,
	}

	testtool.WaitSomething(t, time.Minute, func() error {
		b, err := devtools4chains.RPCCallJSON(rpcInfo, "getblockcount", nil, nil)
		if b != nil && strings.Contains(string(b), "Loading wallet") {
			return fmt.Errorf("Loading wallet")
		}
		return err
	})

	fmt.Println("BTC addr:", c.address)
	_, err = devtools4chains.RPCCallJSON(rpcInfo, "importaddress", []string{c.address}, nil)
	rq.Nil(err)

	var coinbaseAddress string
	_, err = devtools4chains.RPCCallJSON(rpcInfo, "getnewaddress", []interface{}{}, &coinbaseAddress)
	rq.Nil(err)
	fmt.Println("coinbase address", coinbaseAddress)

	var sendtoAddress string
	_, err = devtools4chains.RPCCallJSON(rpcInfo, "getnewaddress", []interface{}{}, &sendtoAddress)
	rq.Nil(err)
	fmt.Println("sendto address", sendtoAddress)

	_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{101, coinbaseAddress}, nil)
	rq.Nil(err)

	_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendtoaddress", []interface{}{c.address, 1.1}, nil)
	rq.Nil(err)

	_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{101, coinbaseAddress}, nil)
	rq.Nil(err)

	var unspents []omni.ListUnspentResult
	_, err = devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{c.address}}, &unspents)
	rq.Nil(err)
	utxo := unspents[0]
	fmt.Printf("%#v\n", utxo)

	{ //RPC创建交易， SDK签名
		var createdTx string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "createrawtransaction", []interface{}{
			// []interface{}{map[string]interface{}{
			// 	"txid": utxo.TxID,
			// 	"vout": utxo.Vout,
			// }},
			// []interface{}{map[string]interface{}{
			// 	sendtoAddress: 1.09999,
			// }},
			json.RawMessage(fmt.Sprintf(`[{"txid":"%s","vout":%d}]`, utxo.TxID, utxo.Vout)),
			json.RawMessage(fmt.Sprintf(`[{"%s":%f}]`, sendtoAddress, 1.09999)),
		}, &createdTx)
		rq.Nil(err)
		fmt.Println("created tx:", createdTx)

		m := map[string]interface{}{
			"RawTx": createdTx,
			"Inputs": []map[string]interface{}{
				{
					"txid":         utxo.TxID,
					"vout":         utxo.Vout,
					"scriptPubKey": utxo.ScriptPubKey,
				},
			},
		}
		msgB, err := json.Marshal(&m)
		rq.NoError(err)
		fmt.Println("msgB", string(msgB))

		sig, err := w.Sign("BTC", hex.EncodeToString(msgB))
		rq.NoError(err)

		fmt.Println("sig:", sig)
		var txid string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{sig}, &txid)
		rq.Nil(err)
		fmt.Println("txid:", txid)

		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{1, coinbaseAddress}, nil)
		rq.Nil(err)

		{ // validate utxo for receiver
			resp, err := devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{sendtoAddress}}, &unspents)
			rq.Nil(err)

			fmt.Println("utxo for sendto address", string(resp))
			rq.Len(unspents, 1, "需要有1个UTXO")
		}
	}

	if 2 < 1 { //这种办法暂时不能顺利的构建tx
		var tx *btc.BTCTransaction
		unspent := new(btc.BTCUnspent) //java: new btc.BTCUnspent()
		unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, "")

		amount, err := btc.NewBTCAmount(0.6)
		rq.Nil(err)

		toAddressA1, err := btc.NewBTCAddressFromString(coinbaseAddress, btc.ChainRegtest)
		rq.Nil(err)

		outputAmount := btc.BTCOutputAmount{} //java: new btc.BTCOutputAmount()
		outputAmount.Add(toAddressA1, amount)

		feeRate := int64(80)

		changeAddress, err := btc.NewBTCAddressFromString(c.address, btc.ChainRegtest)
		rq.Nil(err)

		tx, err = btc.NewBTCTransaction(unspent, &outputAmount, changeAddress, feeRate, btc.ChainRegtest)
		rq.Nil(err)
		_ = tx
	}

	// var signedHex string
	// transferAmount := 2.3
	// chainID := btc.ChainRegtest

}
