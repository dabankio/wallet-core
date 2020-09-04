package omni

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dabankio/devtools4chains"

	"github.com/dabankio/wallet-core/core/btc"
	"github.com/dabankio/wallet-core/core/eth/internalized/testtool"
	"github.com/dabankio/wallet-core/core/omni"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 多重签名转账
// 本地起一条全新的链，发布一个omni代币，在此基础上进行多签测试
// 5个地址(01234)，0为矿工同时也是代币的拥有者，1-2-3创建了多签地址，0转账代币给多签地址，1-3签名转账给4
func TestMultisig(t *testing.T) {
	rq := require.New(t)
	killFunc, nodeInfo, err := devtools4chains.DockerRunOmnicored(devtools4chains.DockerRunOptions{
		AutoRemove: true,
		Image:      &TestImage,
	})
	rq.NoError(err)
	t.Cleanup(killFunc)

	rpcInfo := devtools4chains.RPCInfo{
		Host:     fmt.Sprintf("http://127.0.0.1:%d", nodeInfo.RPCPort),
		User:     nodeInfo.RPCUser,
		Password: nodeInfo.RPCPwd,
		// Debug: true,
	}

	testtool.WaitSomething(t, time.Minute, func() error {
		_, err := devtools4chains.RPCCallJSON(rpcInfo, "getblockcount", nil, nil)
		return err
	})
	for _, add := range presetAddrs {
		_, err := devtools4chains.RPCCallJSON(rpcInfo, "importprivkey", []string{add.Privkey}, nil)
		rq.Nil(err)
	}
	a0, a1, a2, a3, a4 := presetAddrs[0], presetAddrs[1], presetAddrs[2], presetAddrs[3], presetAddrs[4]

	{ // 生成多个块，获取utxo
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{105, a0.Address}, nil)
		rq.Nil(err)
	}

	var multisigAddress, redeemScript string
	{ //a1/a2/a3 生成多签地址,往多签地址转入btc
		keys := strings.Join([]string{a1.Pubkey, a2.Pubkey, a3.Pubkey}, ",")
		ret, err := btc.NewMultiSigAddress(2, btc.ChainRegtest, keys)
		assert.Nil(t, err)
		arr := strings.Split(ret, ",")
		multisigAddress, redeemScript = arr[0], arr[1]

		//导入到钱包
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "importaddress", []string{multisigAddress}, nil)
		rq.Nil(err)

		// 给多签地址转账一部分btc，以产生utxo，同时支持dust费用
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendtoaddress", []interface{}{multisigAddress, 23.3}, nil)
		rq.Nil(err)

		{ // 确认多签地址上有足够的btc
			// received, err := cli.Getreceivedbyaddress(multisigAddress, 0)
			// rq.Nil(err, "failed to get received")
			// fmt.Println("btc on multisig address", received)
			// f64, err := strconv.ParseFloat(received, 64)
			// rq.Nil(err, "parse bal failed")
			// rq.False(f64 <= 0, "wrong bal of multisig", f64)
		}

	}

	propertyID := CreateToken(t, rpcInfo, OmniSenddissuancefixedCmd{
		Fromaddress: a0.Address,
		Ecosystem:   2, //2 fot test
		Typ:         1, // 1 for indivisible
		Previousid:  0, // 0 for new tokens
		Category:    "test_omniii",
		Subcategory: "unit_test",
		Name:        "FakeUSDT",
		Amount:      "10000",
	})

	chainID := btc.ChainRegtest

	{ // simple send from a0 to multisig address
		transferAmount := float64(233)
		{
			var unspents []ListUnspentResult
			_, err = devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{a0.Address}}, &unspents)
			rq.Nil(err)
			utxo := unspents[0]

			unspent := new(btc.BTCUnspent)
			unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, utxo.RedeemScript)
			toAddr, err := btc.NewBTCAddressFromString(multisigAddress, chainID)
			rq.Nil(err, "failed to create addr")
			changeAddr, err := btc.NewBTCAddressFromString(a0.Address, chainID)
			rq.Nil(err, "failed to create addr")
			feeRate := int64(80)
			btctx, err := omni.CreateSimpleSendTransaction(propertyID, false, unspent, toAddr, transferAmount, changeAddr, feeRate, chainID)
			rq.Nil(err, "Failed to crate btctx")

			toSignMsg, err := btctx.EncodeToSignCmd()
			rq.Nil(err, "failed to encode to sign")

			btcCoin, _ := btc.New(nil, chainID)
			signedRawHex, err := btcCoin.Sign(toSignMsg, a0.Privkey)
			rq.Nil(err, "failed to sign")

			// decodeTx, err := btc.DecodeRawTransaction(&btcdbtcjson.DecodeRawTransactionCmd{HexTx: signedRawHex}, &chaincfg.RegressionNetParams)
			// rq.Nil(err, "failed to decode signed raw tx")
			// b, _ := json.MarshalIndent(&decodeTx, "", " ")
			// fmt.Println("signed tx(token: a0 > multisig address)", string(b))

			// 广播交易
			_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{signedRawHex}, nil)
			rq.Nil(err)

			// fmt.Println("broadcasted txid", txid)
		}
		{ // 生成一个块确认代币转账
			_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{1, a0.Address}, nil) //生成几个块，确认刚才的交易
			rq.Nil(err)
		}

		{ // 确认代币转账成功
			// bal, err := cli.OmniGetbalance(multisigAddress, propertyID)
			// rq.Nil(err, "Failed to get omni balance")
			// expectedBal := strconv.FormatFloat(transferAmount, 'f', 0, 32)
			// rq.False(bal.Balance != expectedBal, "wrong balance, not ", expectedBal)
		}
	}
	{ // simple send from multisig address to a4
		transferAmount := float64(23)
		{

			var unspents []ListUnspentResult
			_, err = devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{multisigAddress}}, &unspents)
			rq.Nil(err)

			var utxo ListUnspentResult
			for _, u := range unspents {
				if u.Amount > 0.001 {
					utxo = u
					break
				}
			}

			unspent := new(btc.BTCUnspent)
			unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, redeemScript)
			changeAddr, err := btc.NewBTCAddressFromString(multisigAddress, chainID)
			rq.Nil(err, "failed to create addr")
			toAddr, err := btc.NewBTCAddressFromString(a4.Address, chainID)
			rq.Nil(err, "failed to create addr")
			feeRate := int64(80)
			btctx, err := omni.CreateSimpleSendTransaction(propertyID, false, unspent, toAddr, transferAmount, changeAddr, feeRate, chainID)
			rq.Nil(err, "Failed to crate btctx")

			btcCoin, _ := btc.New(nil, chainID)
			var nextSignData string
			{ // a1签名
				toSignMsg, err := btctx.EncodeToSignCmd()
				rq.Nil(err, "failed to encode to sign")

				signedRawHex, err := btcCoin.Sign(toSignMsg, a1.Privkey)
				assert.Nil(t, err)

				// 下一个人的签名消息
				nextSignData, _ = btctx.EncodeToSignCmdForNextSigner(signedRawHex)
			}

			{ // a3签名，并广播交易
				signedRawHex, err := btcCoin.Sign(nextSignData, a3.Privkey)
				assert.Nil(t, err)

				// 广播交易
				_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{signedRawHex}, nil)
				rq.Nil(err)
				assert.Nil(t, err)
				// fmt.Println("broadcasted txid", txid)
			}

		}
		{ // 生成一个块确认代币转账
			_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{1, a0.Address}, nil) //生成几个块，确认刚才的交易
			rq.Nil(err)
		}

		{ // 确认代币转账成功
			var balanceMap map[string]string
			_, err = devtools4chains.RPCCallJSON(rpcInfo, "omni_getbalance", []interface{}{a4.Address, propertyID}, &balanceMap) //生成几个块，确认刚才的交易
			rq.Nil(err, "Failed to get omni balance")
			rq.Equal(balanceMap["balance"], "23")
		}
		{ // 确认代币转账成功
			// bal, err := cli.OmniGetbalance(a4.Address, propertyID)
			// rq.Nil(err, "Failed to get omni balance")
			// expectedBal := strconv.FormatFloat(23, 'f', 0, 32)
			// rq.False(bal.Balance != expectedBal, "wrong balance, not ", expectedBal)
		}
	}

}
