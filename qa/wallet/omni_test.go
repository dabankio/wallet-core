package wallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dabankio/devtools4chains"
	"github.com/dabankio/wallet-core/core/eth/internalized/testtool"
	"github.com/dabankio/wallet-core/qa/omni"
	"github.com/dabankio/wallet-core/wallet"
	"github.com/stretchr/testify/require"
)

func testOmniPubkSign(t *testing.T, w *wallet.Wallet, c ctx) {
	rq := require.New(t)

	killFunc, bitcoinInfo, err := devtools4chains.DockerRunOmnicored(devtools4chains.DockerRunOptions{
		AutoRemove: true,
		Image:      &omniImage,
	})
	rq.NoError(err)
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

	fmt.Println("OMNI addr:", c.address)
	_, err = devtools4chains.RPCCallJSON(rpcInfo, "importaddress", []string{c.address}, nil)
	rq.Nil(err)

	var coinbaseAddress string
	_, err = devtools4chains.RPCCallJSON(rpcInfo, "getnewaddress", []interface{}{}, &coinbaseAddress)
	rq.Nil(err)
	fmt.Println("coinbase address", coinbaseAddress)

	_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{101, coinbaseAddress}, nil)
	rq.Nil(err)
	_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendtoaddress", []interface{}{c.address, 5.1}, nil)
	rq.Nil(err)

	var toAddress string
	_, err = devtools4chains.RPCCallJSON(rpcInfo, "getnewaddress", []interface{}{}, &toAddress)
	rq.Nil(err)
	fmt.Println("sendto address", toAddress)

	_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{102, coinbaseAddress}, nil)
	rq.Nil(err)

	propertyID := omni.CreateToken(t, rpcInfo, omni.OmniSenddissuancefixedCmd{
		Fromaddress: coinbaseAddress,
		Ecosystem:   2, //2 fot test
		Typ:         1, // 1 for indivisible
		Previousid:  0, // 0 for new tokens
		Category:    "test_omni",
		Subcategory: "unit_test",
		Name:        "FakeUSDT",
		Amount:      "10000",
	})

	_, err = devtools4chains.RPCCallJSON(rpcInfo, "omni_send", []interface{}{coinbaseAddress, c.address, propertyID, "998"}, nil)
	rq.Nil(err)
	_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{102, coinbaseAddress}, nil)
	rq.Nil(err)

	t.Run("使用RPC创建交易", func(t *testing.T) {

		var unspents []omni.ListUnspentResult
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{c.address}}, &unspents)
		rq.Nil(err)
		rq.Greater(len(unspents), 0)
		var utxo omni.ListUnspentResult
		for _, out := range unspents {
			if out.Amount > 0.01 {
				utxo = out
				break
			}
		}
		fmt.Printf("utxo %#v\n", utxo)

		{ //RPC创建交易， SDK签名, OMNI构造交易过程参考 https://github.com/OmniLayer/omnicore/wiki/Use-the-raw-transaction-API-to-create-a-Simple-Send-transaction
			var payload string
			_, err = devtools4chains.RPCCallJSON(rpcInfo, "omni_createpayload_simplesend", []interface{}{propertyID, "101"}, &payload)
			rq.NoError(err)

			var createdTx string
			_, err = devtools4chains.RPCCallJSON(rpcInfo, "createrawtransaction", []interface{}{
				json.RawMessage(fmt.Sprintf(`[{"txid":"%s","vout":%d}]`, utxo.TxID, utxo.Vout)),
				json.RawMessage(fmt.Sprintf(`{}`)),
			}, &createdTx)
			rq.Nil(err)
			fmt.Println("created tx:", createdTx)

			var payloadAttached string
			_, err = devtools4chains.RPCCallJSON(rpcInfo, "omni_createrawtx_opreturn", []interface{}{createdTx, payload}, &payloadAttached)
			rq.NoError(err)

			var receiverAttached string
			_, err = devtools4chains.RPCCallJSON(rpcInfo, "omni_createrawtx_reference", []interface{}{payloadAttached, toAddress}, &receiverAttached)
			rq.NoError(err)

			var minerFeeAndChangeAttached string
			_, err = devtools4chains.RPCCallJSON(rpcInfo, "omni_createrawtx_change", []interface{}{
				receiverAttached,
				json.RawMessage(fmt.Sprintf(`[{"txid":"%s","vout":%d, "scriptPubKey": "%s", "value": %f}]`, utxo.TxID, utxo.Vout, utxo.ScriptPubKey, utxo.Amount)),
				c.address,
				0.000035,
			}, &minerFeeAndChangeAttached)
			rq.NoError(err)

			fmt.Println("tx to sign:", minerFeeAndChangeAttached)

			var decodeTx map[string]interface{}
			_, err = devtools4chains.RPCCallJSON(rpcInfo, "decoderawtransaction", []string{minerFeeAndChangeAttached}, &decodeTx)
			rq.NoError(err)
			b, _ := json.MarshalIndent(decodeTx, "", "  ")
			fmt.Println("decodeTx:", string(b))

			m := map[string]interface{}{
				"RawTx": minerFeeAndChangeAttached,
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

			sig, err := w.Sign("OMNI", hex.EncodeToString(msgB))
			rq.NoError(err)

			fmt.Println("sig:", sig)

			{ //omnicored sign
				// var coreSig string
				// resp, err := devtools4chains.RPCCallJSON(rpcInfo, "signrawtransactionwithkey", []interface{}{
				// 	minerFeeAndChangeAttached,
				// 	[]string{privk},
				// 	json.RawMessage(fmt.Sprintf(`[{"txid":"%s","vout":%d,"scriptPubKey":"%s","amount":%f}]`, utxo.TxID, utxo.Vout, utxo.ScriptPubKey, utxo.Amount)),
				// }, nil)
				// rq.Nil(err)
				// fmt.Println("core sig resp", string(resp))
				// fmt.Println("core sig:", coreSig)
			}

			var txid string
			_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{sig}, &txid)
			rq.Nil(err)
			fmt.Println("txid:", txid)

			_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{1, c.address}, nil)
			rq.Nil(err)

			{ // validate utxo for receiver
				resp, err := devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{toAddress}}, &unspents)
				rq.Nil(err)

				fmt.Println("utxo for sendto address", string(resp))
				rq.Len(unspents, 1, "需要有1个UTXO")

				for _, x := range []struct {
					address string
					balance float64
				}{
					{c.address, 998 - 101},
					{toAddress, 101},
				} {
					resp, err = devtools4chains.RPCCallJSON(rpcInfo, "omni_getbalance", []interface{}{x.address, propertyID}, nil)
					rq.Nil(err)
					rq.Contains(string(resp), fmt.Sprintf("%v", x.balance))
					fmt.Println("balance", string(resp))
				}
			}
		}
	})

	t.Run("使用SDK创建交易", func(t *testing.T) {
		t.Skip()
		// TBD
	})

}
