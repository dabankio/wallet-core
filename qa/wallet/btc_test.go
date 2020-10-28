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
		r.Nil(t, err)
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "getnewaddress", []interface{}{}, &coinbaseAddress)
		r.Nil(t, err)
		t.Log("coinbase address", coinbaseAddress)

		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{101, coinbaseAddress}, nil)
		r.Nil(t, err)
	})

	fnPrepareFundAndGetNewAddress := func(t *testing.T) { //给c.address转账，并为sendtoAddress设置新值
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendtoaddress", []interface{}{c.address, 1.1}, nil)
		r.Nil(t, err)
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{101, coinbaseAddress}, nil)
		r.Nil(t, err)

		_, err = devtools4chains.RPCCallJSON(rpcInfo, "getnewaddress", []interface{}{}, &sendtoAddress)
		r.Nil(t, err)
	}

	fnGetUTXO := func(t *testing.T) omni.ListUnspentResult {
		var unspents []omni.ListUnspentResult
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{c.address}}, &unspents)
		r.Nil(t, err)
		r.Len(t, unspents, 1)
		return unspents[0]
	}

	var unspents []omni.ListUnspentResult
	t.Run("使用RPC创建交易", func(t *testing.T) {
		fnPrepareFundAndGetNewAddress(t)
		utxo := fnGetUTXO(t)

		var createdTx string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "createrawtransaction", []interface{}{
			// []interface{}{map[string]interface{}{"txid": utxo.TxID, "vout": utxo.Vout}},
			// []interface{}{map[string]interface{}{sendtoAddress: 1.09999}},
			json.RawMessage(fmt.Sprintf(`[{"txid":"%s","vout":%d}]`, utxo.TxID, utxo.Vout)),
			json.RawMessage(fmt.Sprintf(`[{"%s":%f}]`, sendtoAddress, 1.09999)),
		}, &createdTx)
		r.Nil(t, err)
		t.Log("created tx:", createdTx)

		m := map[string]interface{}{
			"RawTx": createdTx,
			"Inputs": []map[string]interface{}{{
				"txid":         utxo.TxID,
				"vout":         utxo.Vout,
				"scriptPubKey": utxo.ScriptPubKey,
			}},
		}
		msgB, err := json.Marshal(&m)
		r.NoError(t, err)
		t.Log("msgB", string(msgB))

		sig, err := w.Sign("BTC", hex.EncodeToString(msgB))
		r.NoError(t, err)

		t.Log("sig:", sig)
		var txid string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{sig}, &txid)
		r.Nil(t, err)
		t.Log("txid:", txid)

		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{1, coinbaseAddress}, nil)
		r.Nil(t, err)

		{ // validate utxo for receiver
			resp, err := devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{sendtoAddress}}, &unspents)
			r.Nil(t, err)

			t.Log("utxo for sendto address", string(resp))
			r.Len(t, unspents, 1, "需要有1个UTXO")
		}

	})

	t.Run("使用SDK创建交易", func(t *testing.T) {
		fnPrepareFundAndGetNewAddress(t)
		utxo := fnGetUTXO(t)

		var tx *btc.BTCTransaction
		unspent := new(btc.BTCUnspent) //java: new btc.BTCUnspent()
		unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, "")

		amount, err := btc.NewBTCAmount(0.0021)
		r.Nil(t, err)

		toAddress, err := btc.NewBTCAddressFromString(coinbaseAddress, btc.ChainRegtest)
		r.Nil(t, err)

		outputAmount := btc.BTCOutputAmount{} //java: new btc.BTCOutputAmount()
		outputAmount.Add(toAddress, amount)

		feeRate := int64(80)

		changeAddress, err := btc.NewBTCAddressFromString(c.address, btc.ChainRegtest) //找零地址
		r.Nil(t, err)

		tx, err = btc.NewBTCTransaction(unspent, &outputAmount, changeAddress, feeRate, btc.ChainRegtest)
		r.Nil(t, err)

		toSignTx, err := tx.EncodeToSignCmd() //编码为可签名的格式
		r.NoError(t, err)

		sig, err := w.Sign("BTC", toSignTx) //签名
		r.NoError(t, err)

		t.Log("sig:", sig)
		var txid string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{sig}, &txid)
		r.Nil(t, err)
		t.Log("txid:", txid)
	})

	t.Run("使用SDK创建交易,但不提供scriptPubKey", func(t *testing.T) {
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendtoaddress", []interface{}{c.address, 1.1}, nil)
		r.Nil(t, err)

		var unspents []omni.ListUnspentResult
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{c.address}}, &unspents)
		r.Nil(t, err)
		utxo := unspents[0]

		var tx *btc.BTCTransaction
		unspent := new(btc.BTCUnspent) //java: new btc.BTCUnspent()
		// unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, "")
		unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, "", "")

		amount, err := btc.NewBTCAmount(0.0021)
		r.Nil(t, err)

		toAddress, err := btc.NewBTCAddressFromString(coinbaseAddress, btc.ChainRegtest)
		r.Nil(t, err)

		outputAmount := btc.BTCOutputAmount{} //java: new btc.BTCOutputAmount()
		outputAmount.Add(toAddress, amount)

		feeRate := int64(80)

		changeAddress, err := btc.NewBTCAddressFromString(c.address, btc.ChainRegtest) //找零地址
		r.Nil(t, err)

		tx, err = btc.NewBTCTransaction(unspent, &outputAmount, changeAddress, feeRate, btc.ChainRegtest)
		r.Nil(t, err)

		toSignTx, err := tx.EncodeToSignCmd() //编码为可签名的格式
		r.NoError(t, err)

		sig, err := w.Sign("BTC", toSignTx) //签名
		r.NoError(t, err)

		t.Log("sig:", sig)
		var txid string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{sig}, &txid)
		r.Nil(t, err)
		t.Log("txid:", txid)
	})

	t.Run("使用SDK创建交易", func(t *testing.T) {
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendtoaddress", []interface{}{c.address, 1.1}, nil)
		r.Nil(t, err)

		var unspents []omni.ListUnspentResult
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{c.address}}, &unspents)
		r.Nil(t, err)
		utxo := unspents[0]

		var tx *btc.BTCTransaction
		unspent := new(btc.BTCUnspent) //java: new btc.BTCUnspent()
		// unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, "")
		unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, "", "")

		amount, err := btc.NewBTCAmount(0.0021)
		r.Nil(t, err)

		toAddress, err := btc.NewBTCAddressFromString(coinbaseAddress, btc.ChainRegtest)
		r.Nil(t, err)

		outputAmount := btc.BTCOutputAmount{} //java: new btc.BTCOutputAmount()
		outputAmount.Add(toAddress, amount)

		feeRate := int64(80)

		changeAddress, err := btc.NewBTCAddressFromString(c.address, btc.ChainRegtest) //找零地址
		r.Nil(t, err)

		tx, err = btc.NewBTCTransaction(unspent, &outputAmount, changeAddress, feeRate, btc.ChainRegtest)
		r.Nil(t, err)

		toSignTx, err := tx.EncodeToSignCmd() //编码为可签名的格式
		r.NoError(t, err)

		sig, err := w.Sign("BTC", toSignTx) //签名
		r.NoError(t, err)

		t.Log("sig:", sig)
		var txid string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{sig}, &txid)
		r.Nil(t, err)
		t.Log("txid:", txid)
	})

}
func testBTCPubkSignSegwit(t *testing.T, w *wallet.Wallet, _ ctx) {
	T := t
	rq := r.New(t)

	w.AddFlag(wallet.FlagBTCUseSegWitFormat)

	pubk, err := w.DerivePublicKey("BTC")
	r.NoError(t, err)
	address, err := w.DeriveAddress("BTC")
	r.NoError(t, err)
	c := ctx{pubk: pubk, address: address}
	t.Log("pubk:", pubk)

	var rpcInfo devtools4chains.RPCInfo
	var coinbaseAddress string
	var sendtoAddress string

	r.True(t, strings.HasPrefix(c.address, "2"), "隔离见证地址应该以2开头")
	privk, err := w.DerivePrivateKey("BTC")
	r.NoError(t, err)

	balanceInBTC := 1.1               //btc
	eachFee := 0.0001                 //in btc
	eachOut := balanceInBTC - eachFee //btc
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
		r.Nil(t, err)
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "getnewaddress", []interface{}{}, &coinbaseAddress)
		r.Nil(t, err)
		t.Log("coinbase address", coinbaseAddress)
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "getnewaddress", []interface{}{}, &sendtoAddress)
		r.Nil(t, err)
		t.Log("sendto address", sendtoAddress)

		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{101, coinbaseAddress}, nil)
		r.Nil(t, err)
	})

	fnPrepareFundAndGetNewAddress := func(t *testing.T) { //给c.address转账，并为sendtoAddress设置新值
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendtoaddress", []interface{}{c.address, 1.1}, nil)
		r.Nil(t, err)
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{101, coinbaseAddress}, nil)
		r.Nil(t, err)

		_, err = devtools4chains.RPCCallJSON(rpcInfo, "getnewaddress", []interface{}{}, &sendtoAddress)
		r.Nil(t, err)
	}

	fnGetUTXO := func(t *testing.T) omni.ListUnspentResult {
		var unspents []omni.ListUnspentResult
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{c.address}}, &unspents)
		r.Nil(t, err)
		r.Len(t, unspents, 1)
		return unspents[0]
	}
	var unspents []omni.ListUnspentResult //share var

	t.Run("使用RPC创建交易,提供scriptPubk", func(t *testing.T) {
		t.Skip()
		fnPrepareFundAndGetNewAddress(t)
		utxo := fnGetUTXO(t)
		// t.Logf("utxo %#v\n", utxo)

		var createdTx string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "createrawtransaction", []interface{}{
			json.RawMessage(fmt.Sprintf(`[{"txid":"%s","vout":%d}]`, utxo.TxID, utxo.Vout)), // []interface{}{map[string]interface{}{"txid": utxo.TxID, "vout": utxo.Vout}},
			json.RawMessage(fmt.Sprintf(`[{"%s":%f}]`, sendtoAddress, eachOut)),             // []interface{}{map[string]interface{}{sendtoAddress: 1.09999}},
		}, &createdTx)
		r.Nil(t, err)
		// t.Log("created tx:", createdTx)

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
		r.NoError(t, err)
		// t.Log("msgB", string(msgB))

		sig, err := w.Sign("BTC", hex.EncodeToString(msgB))
		// t.Log("signErr", err)
		r.NoError(t, err)

		walletSig := map[string]interface{}{}
		wsResp, err := devtools4chains.RPCCallJSON(rpcInfo, "signrawtransactionwithwallet", []interface{}{createdTx}, &walletSig)
		_ = wsResp
		// t.Log("wsResp:", string(wsResp))
		r.Nil(t, err)
		// t.Log("wsig", walletSig["hex"])

		var txid string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{sig}, &txid)
		r.Nil(t, err)
		t.Log("txid:", txid)

		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{1, coinbaseAddress}, nil)
		r.Nil(t, err)

		{ // validate utxo for receiver
			resp, err := devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{sendtoAddress}}, &unspents)
			r.Nil(t, err)

			t.Log("utxo for sendto address", string(resp))
			rq.Len(unspents, 1, "需要有1个UTXO")
		}
	})

	t.Run("使用RPC创建交易,签名时不提供scriptPubk(内部自动创建scriptPubk)", func(t *testing.T) {
		fnPrepareFundAndGetNewAddress(t)

		t.Log("privK:", privk)
		//获取新的收币地址
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "getnewaddress", []interface{}{}, &sendtoAddress)
		r.Nil(t, err)
		t.Log("sendto address", sendtoAddress)

		utxo := fnGetUTXO(t)
		t.Logf("utxo: %#v", utxo)

		var createdTx string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "createrawtransaction", []interface{}{
			json.RawMessage(fmt.Sprintf(`[{"txid":"%s","vout":%d}]`, utxo.TxID, utxo.Vout)), // []interface{}{map[string]interface{}{"txid": utxo.TxID, "vout": utxo.Vout}},
			json.RawMessage(fmt.Sprintf(`[{"%s":%f}]`, sendtoAddress, eachOut)),             // []interface{}{map[string]interface{}{sendtoAddress: 1.09999}},
			// json.RawMessage(fmt.Sprintf(`[{"%s":%f}, {"%s":%f}]`, sendtoAddress, eachOut - 0.1, c.address, 0.1)),             // []interface{}{map[string]interface{}{sendtoAddress: 1.09999}},
		}, &createdTx)
		r.Nil(t, err)
		t.Log("created tx:", createdTx)

		m := map[string]interface{}{
			"RawTx": createdTx,
			"Inputs": []map[string]interface{}{{
				"txid": utxo.TxID,
				"vout": utxo.Vout,
				// "scriptPubKey": utxo.ScriptPubKey, //不要传入scriptPubk,RedeemScript,以测试自动生成的解锁脚本是否有效
				// "redeemScript": utxo.RedeemScript,
				"amount": utxo.Amount,
			}},
		}
		msgB, err := json.Marshal(&m)
		r.NoError(t, err)
		// t.Log("msgB", string(msgB))

		sig, err := w.Sign("BTC", hex.EncodeToString(msgB))
		r.NoError(t, err)
		// t.Log("sig:", sig)
		walletSig := map[string]interface{}{}
		wsResp, err := devtools4chains.RPCCallJSON(rpcInfo, "signrawtransactionwithwallet", []interface{}{createdTx}, &walletSig)
		_ = wsResp
		// t.Log("wsResp:", string(wsResp))
		r.Nil(t, err)
		// t.Log("wsig", walletSig["hex"])

		var txid string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{sig}, &txid)
		r.Nil(t, err)
		t.Log("txid:", txid)

		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{1, coinbaseAddress}, nil)
		r.Nil(t, err)

		{ // validate utxo for receiver
			resp, err := devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{sendtoAddress}}, &unspents)
			r.Nil(t, err)
			_ = resp
			// t.Log("utxo for sendto address", string(resp))
			rq.Len(unspents, 1, "需要有1个UTXO")
		}
	})

	t.Run("使用SDK创建交易,不提供解锁脚本", func(t *testing.T) {
		fnPrepareFundAndGetNewAddress(t)
		utxo := fnGetUTXO(t)

		unspent := new(btc.BTCUnspent) //java: new btc.BTCUnspent()
		unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, "", "")

		amount, err := btc.NewBTCAmount(balanceInBTC - 0.2)
		r.Nil(t, err)

		toAddress, err := btc.NewBTCAddressFromString(sendtoAddress, btc.ChainRegtest)
		r.Nil(t, err)

		outputAmount := btc.BTCOutputAmount{} //java: new btc.BTCOutputAmount()
		outputAmount.Add(toAddress, amount)

		feeRate := int64(80)

		changeAddress, err := btc.NewBTCAddressFromString(c.address, btc.ChainRegtest) //找零地址
		r.Nil(t, err)

		tx, err := btc.NewBTCTransaction(unspent, &outputAmount, changeAddress, feeRate, btc.ChainRegtest)
		r.Nil(t, err)

		toSignTx, err := tx.EncodeToSignCmd() //编码为可签名的格式
		r.NoError(t, err)

		h, e := tx.Encode()
		r.NoError(t, e)
		t.Log("toSignTx:", h)
		sig, err := w.Sign("BTC", toSignTx) //签名
		r.NoError(t, err)
		t.Log("sig:", sig)
		var txid string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{sig}, &txid)
		r.Nil(t, err)
		t.Log("txid:", txid)

		{ // validate utxo for receiver
			resp, err := devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{sendtoAddress}}, &unspents)
			r.Nil(t, err)
			_ = resp
			// t.Log("utxo for sendto address", string(resp))
			rq.Len(unspents, 1, "需要有1个UTXO")
		}
	})
}
