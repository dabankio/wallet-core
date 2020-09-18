package wallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/dabankio/devtools4chains"
	"github.com/dabankio/wallet-core/core/btc"
	"github.com/dabankio/wallet-core/qa/omni"
	"github.com/dabankio/wallet-core/wallet"
	r "github.com/stretchr/testify/require"
)

func testBTCPubkSign(t *testing.T, w *wallet.Wallet, c ctx) {
	T := t
	rq := r.New(t)

	var err error
	var rpcInfo devtools4chains.RPCInfo
	var coinbaseAddress string
	var sendtoAddress string

	t.Run("准备节点", func(t *testing.T) {
		killFunc, bitcoinInfo, err := devtools4chains.DockerRunBitcoin(devtools4chains.DockerRunOptions{
			AutoRemove: true, Image: &bbcImage,
		})
		r.NoError(t, err)
		T.Cleanup(killFunc)

		rpcInfo = devtools4chains.RPCInfo{
			Host:     fmt.Sprintf("http://127.0.0.1:%d", bitcoinInfo.RPCPort),
			User:     bitcoinInfo.RPCUser,
			Password: bitcoinInfo.RPCPwd,
		}

		_, err = devtools4chains.RPCCallJSON(rpcInfo, "importaddress", []string{c.address}, nil)
		rq.Nil(err)
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "getnewaddress", []interface{}{}, &coinbaseAddress)
		rq.Nil(err)
		fmt.Println("coinbase address", coinbaseAddress)
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "getnewaddress", []interface{}{}, &sendtoAddress)
		rq.Nil(err)
		fmt.Println("sendto address", sendtoAddress)

		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{101, coinbaseAddress}, nil)
		rq.Nil(err)
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendtoaddress", []interface{}{c.address, 1.1}, nil)
		rq.Nil(err)
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{101, coinbaseAddress}, nil)
		rq.Nil(err)
	})

	t.Run("使用RPC创建交易", func(t *testing.T) {
		var unspents []omni.ListUnspentResult
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{c.address}}, &unspents)
		rq.Nil(err)
		utxo := unspents[0]
		fmt.Printf("utxo %#v\n", utxo)

		var createdTx string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "createrawtransaction", []interface{}{
			// []interface{}{map[string]interface{}{"txid": utxo.TxID, "vout": utxo.Vout}},
			// []interface{}{map[string]interface{}{sendtoAddress: 1.09999}},
			json.RawMessage(fmt.Sprintf(`[{"txid":"%s","vout":%d}]`, utxo.TxID, utxo.Vout)),
			json.RawMessage(fmt.Sprintf(`[{"%s":%f}]`, sendtoAddress, 1.09999)),
		}, &createdTx)
		rq.Nil(err)
		fmt.Println("created tx:", createdTx)

		m := map[string]interface{}{
			"RawTx": createdTx,
			"Inputs": []map[string]interface{}{{
				"txid":         utxo.TxID,
				"vout":         utxo.Vout,
				"scriptPubKey": utxo.ScriptPubKey,
			}},
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

	})

	t.Run("使用SDK创建交易", func(t *testing.T) {
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendtoaddress", []interface{}{c.address, 1.1}, nil)
		rq.Nil(err)

		var unspents []omni.ListUnspentResult
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{c.address}}, &unspents)
		rq.Nil(err)
		utxo := unspents[0]

		var tx *btc.BTCTransaction
		unspent := new(btc.BTCUnspent) //java: new btc.BTCUnspent()
		unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, "")

		amount, err := btc.NewBTCAmount(0.0021)
		rq.Nil(err)

		toAddress, err := btc.NewBTCAddressFromString(coinbaseAddress, btc.ChainRegtest)
		rq.Nil(err)

		outputAmount := btc.BTCOutputAmount{} //java: new btc.BTCOutputAmount()
		outputAmount.Add(toAddress, amount)

		feeRate := int64(80)

		changeAddress, err := btc.NewBTCAddressFromString(c.address, btc.ChainRegtest) //找零地址
		rq.Nil(err)

		tx, err = btc.NewBTCTransaction(unspent, &outputAmount, changeAddress, feeRate, btc.ChainRegtest)
		rq.Nil(err)

		toSignTx, err := tx.EncodeToSignCmd() //编码为可签名的格式
		rq.NoError(err)

		sig, err := w.Sign("BTC", toSignTx) //签名
		rq.NoError(err)

		fmt.Println("sig:", sig)
		var txid string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{sig}, &txid)
		rq.Nil(err)
		fmt.Println("txid:", txid)
	})

	t.Run("使用SDK创建交易,但不提供scriptPubKey", func(t *testing.T) {
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendtoaddress", []interface{}{c.address, 1.1}, nil)
		rq.Nil(err)

		var unspents []omni.ListUnspentResult
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{c.address}}, &unspents)
		rq.Nil(err)
		utxo := unspents[0]

		var tx *btc.BTCTransaction
		unspent := new(btc.BTCUnspent) //java: new btc.BTCUnspent()
		// unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, "")
		unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, "", "")

		amount, err := btc.NewBTCAmount(0.0021)
		rq.Nil(err)

		toAddress, err := btc.NewBTCAddressFromString(coinbaseAddress, btc.ChainRegtest)
		rq.Nil(err)

		outputAmount := btc.BTCOutputAmount{} //java: new btc.BTCOutputAmount()
		outputAmount.Add(toAddress, amount)

		feeRate := int64(80)

		changeAddress, err := btc.NewBTCAddressFromString(c.address, btc.ChainRegtest) //找零地址
		rq.Nil(err)

		tx, err = btc.NewBTCTransaction(unspent, &outputAmount, changeAddress, feeRate, btc.ChainRegtest)
		rq.Nil(err)

		toSignTx, err := tx.EncodeToSignCmd() //编码为可签名的格式
		rq.NoError(err)

		sig, err := w.Sign("BTC", toSignTx) //签名
		rq.NoError(err)

		fmt.Println("sig:", sig)
		var txid string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{sig}, &txid)
		rq.Nil(err)
		fmt.Println("txid:", txid)
	})

	t.Run("隔离见证地址:使用SDK创建交易", func(t *testing.T) {

		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendtoaddress", []interface{}{c.address, 1.1}, nil)
		rq.Nil(err)

		var unspents []omni.ListUnspentResult
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{c.address}}, &unspents)
		rq.Nil(err)
		utxo := unspents[0]

		var tx *btc.BTCTransaction
		unspent := new(btc.BTCUnspent) //java: new btc.BTCUnspent()
		// unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, "")
		unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, "", "")

		amount, err := btc.NewBTCAmount(0.0021)
		rq.Nil(err)

		toAddress, err := btc.NewBTCAddressFromString(coinbaseAddress, btc.ChainRegtest)
		rq.Nil(err)

		outputAmount := btc.BTCOutputAmount{} //java: new btc.BTCOutputAmount()
		outputAmount.Add(toAddress, amount)

		feeRate := int64(80)

		changeAddress, err := btc.NewBTCAddressFromString(c.address, btc.ChainRegtest) //找零地址
		rq.Nil(err)

		tx, err = btc.NewBTCTransaction(unspent, &outputAmount, changeAddress, feeRate, btc.ChainRegtest)
		rq.Nil(err)

		toSignTx, err := tx.EncodeToSignCmd() //编码为可签名的格式
		rq.NoError(err)

		sig, err := w.Sign("BTC", toSignTx) //签名
		rq.NoError(err)

		fmt.Println("sig:", sig)
		var txid string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{sig}, &txid)
		rq.Nil(err)
		fmt.Println("txid:", txid)
	})

}
func testBTCPubkSignSegwit(t *testing.T, w *wallet.Wallet, c ctx) {
	T := t
	rq := r.New(t)

	// w.AddFlag(wallet.FlagBTCUseSegWitFormat)
	var err error
	var rpcInfo devtools4chains.RPCInfo
	var coinbaseAddress string
	var sendtoAddress string

	c.address, err = w.DeriveAddress("BTC")
	r.NoError(t, err)
	r.True(t, strings.HasPrefix(c.address, "2"), "隔离见证地址应该以2开头")
	privk, err := w.DerivePrivateKey("BTC")
	r.NoError(t, err)

	t.Run("准备节点", func(t *testing.T) {
		killFunc, bitcoinInfo, err := devtools4chains.DockerRunBitcoin(devtools4chains.DockerRunOptions{
			AutoRemove: true, Image: &bbcImage,
		})
		r.NoError(t, err)
		T.Cleanup(killFunc)

		rpcInfo = devtools4chains.RPCInfo{
			Host:     fmt.Sprintf("http://127.0.0.1:%d", bitcoinInfo.RPCPort),
			User:     bitcoinInfo.RPCUser,
			Password: bitcoinInfo.RPCPwd,
		}

		// _, err = devtools4chains.RPCCallJSON(rpcInfo, "importaddress", []string{c.address}, nil)
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "importprivkey", []string{privk}, nil)
		rq.Nil(err)
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "getnewaddress", []interface{}{}, &coinbaseAddress)
		rq.Nil(err)
		fmt.Println("coinbase address", coinbaseAddress)
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "getnewaddress", []interface{}{}, &sendtoAddress)
		rq.Nil(err)
		fmt.Println("sendto address", sendtoAddress)

		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{101, coinbaseAddress}, nil)
		rq.Nil(err)
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendtoaddress", []interface{}{c.address, 1.1}, nil)
		rq.Nil(err)
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{101, coinbaseAddress}, nil)
		rq.Nil(err)
	})

	t.Run("使用RPC创建交易", func(t *testing.T) {
		var unspents []omni.ListUnspentResult
		unB, err := devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{c.address}}, &unspents)
		rq.Nil(err)
		utxo := unspents[0]
		fmt.Printf("utxo %#v\n", utxo)
		fmt.Println("json utxo:", string(unB))

		var createdTx string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "createrawtransaction", []interface{}{
			json.RawMessage(fmt.Sprintf(`[{"txid":"%s","vout":%d}]`, utxo.TxID, utxo.Vout)), // []interface{}{map[string]interface{}{"txid": utxo.TxID, "vout": utxo.Vout}},
			json.RawMessage(fmt.Sprintf(`[{"%s":%f}]`, sendtoAddress, 1.09999)),             // []interface{}{map[string]interface{}{sendtoAddress: 1.09999}},
		}, &createdTx)
		rq.Nil(err)
		fmt.Println("created tx:", createdTx)

		m := map[string]interface{}{
			"RawTx": createdTx,
			"Inputs": []map[string]interface{}{{
				"txid":         utxo.TxID,
				"vout":         utxo.Vout,
				"scriptPubKey": utxo.ScriptPubKey,
				"redeemScript": utxo.RedeemScript,
				"amount":       utxo.Amount,
			}},
		}
		msgB, err := json.Marshal(&m)
		rq.NoError(err)
		fmt.Println("msgB", string(msgB))

		sig, err := w.Sign("BTC", hex.EncodeToString(msgB))
		fmt.Println("signErr", err)
		rq.NoError(err)

		fmt.Println("sig:", sig)

		walletSig := map[string]interface{}{}
		wsResp, err := devtools4chains.RPCCallJSON(rpcInfo, "signrawtransactionwithwallet", []interface{}{createdTx}, &walletSig)
		fmt.Println("wsResp:", string(wsResp))
		rq.Nil(err)
		fmt.Println("wsig", walletSig["hex"])

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

	})

	t.Run("隔离见证地址:使用SDK创建交易", func(t *testing.T) {
		t.Skip("暂时不测试用sdk创建")

		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendtoaddress", []interface{}{c.address, 1.1}, nil)
		rq.Nil(err)

		var unspents []omni.ListUnspentResult
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{c.address}}, &unspents)
		rq.Nil(err)
		utxo := unspents[0]

		var tx *btc.BTCTransaction
		unspent := new(btc.BTCUnspent) //java: new btc.BTCUnspent()
		// unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, "")
		unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, "", "")

		amount, err := btc.NewBTCAmount(0.0021)
		rq.Nil(err)

		toAddress, err := btc.NewBTCAddressFromString(coinbaseAddress, btc.ChainRegtest)
		rq.Nil(err)

		outputAmount := btc.BTCOutputAmount{} //java: new btc.BTCOutputAmount()
		outputAmount.Add(toAddress, amount)

		feeRate := int64(80)

		changeAddress, err := btc.NewBTCAddressFromString(c.address, btc.ChainRegtest) //找零地址
		rq.Nil(err)

		tx, err = btc.NewBTCTransaction(unspent, &outputAmount, changeAddress, feeRate, btc.ChainRegtest)
		rq.Nil(err)

		toSignTx, err := tx.EncodeToSignCmd() //编码为可签名的格式
		rq.NoError(err)

		sig, err := w.Sign("BTC", toSignTx) //签名
		rq.NoError(err)

		fmt.Println("sig:", sig)
		var txid string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{sig}, &txid)
		rq.Nil(err)
		fmt.Println("txid:", txid)
	})

}
