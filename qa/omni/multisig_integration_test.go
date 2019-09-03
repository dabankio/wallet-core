// +build integration

package omni

import (
	"fmt"
	"github.com/lomocoin/omnicli"
	"github.com/lomocoin/omnicli/btcjson"
	"github.com/lomocoin/wallet-core/core/btc"
	"github.com/lomocoin/wallet-core/core/omni"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
)

// 多重签名转账
// 本地起一条全新的链，发布一个omni代币，在此基础上进行多签测试
// 5个地址(01234)，0为矿工同时也是代币的拥有者，1-2-3创建了多签地址，0转账代币给多签地址，1-3签名转账给4
func TestMultisig(t *testing.T) {
	rq := require.New(t)
	cli, killomnicored, err := omnicli.RunOmnicored(&omnicli.RunOptions{NewTmpDir: true})
	rq.Nil(err)
	defer killomnicored()

	rq.Nil(importAddrs(cli, presetAddrs), "Failed to import privkeys")
	a0, a1, a2, a3, a4 := presetAddrs[0], presetAddrs[1], presetAddrs[2], presetAddrs[3], presetAddrs[4]

	{ // 生成多个块，获取utxo
		omnicli.NoPrintCmd(func() {
			_, err := cli.Generatetoaddress(103, a0.Address, nil)
			rq.False(err != nil, "Failed to generate to address", err)
		})
	}

	var multisigAddress, redeemScript string
	{ //a1/a2/a3 生成多签地址,往多签地址转入btc
		keys := strings.Join([]string{a1.Pubkey, a2.Pubkey, a3.Pubkey}, ",")
		ret, err := btc.NewMultiSigAddress(2, btc.ChainRegtest, keys)
		assert.Nil(t, err)
		arr := strings.Split(ret, ",")
		multisigAddress, redeemScript = arr[0], arr[1]

		//导入到钱包
		err = cli.Importaddress(btcjson.ImportAddressCmd{
			Address: multisigAddress,
		})
		assert.Nil(t, err)

		// 给多签地址转账一部分btc，以产生utxo，同时支持dust费用
		txid, err := cli.Sendtoaddress(&btcjson.SendToAddressCmd{
			Address: multisigAddress,
			Amount:  2.33,
		})
		assert.Nil(t, err)
		_ = txid

		{ // 确认多签地址上有足够的btc
			received, err := cli.Getreceivedbyaddress(multisigAddress, 0)
			rq.Nil(err, "failed to get received")
			fmt.Println("btc on multisig address", received)
			f64, err := strconv.ParseFloat(received, 64)
			rq.Nil(err, "parse bal failed")
			rq.False(f64 <= 0, "wrong bal of multisig", f64)
		}

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

	chainID := btc.ChainRegtest

	{ // simple send from a0 to multisig address
		transferAmount := float64(233)
		{
			var utxo btcjson.ListUnspentResult
			omnicli.NoPrintCmd(func() {
				unspents, err := cli.Listunspent(0, 999, []string{a0.Address})
				rq.Nil(err, "Failed to list unspent")
				rq.False(len(unspents) == 0, "no unspent find")
				utxo = unspents[0]
			})

			unspent := new(btc.BTCUnspent)
			unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, utxo.RedeemScript)
			toAddr, err := btc.NewBTCAddressFromString(multisigAddress, chainID)
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
			rq.Nil(err, "failed to sign")

			// decodeTx, err := btc.DecodeRawTransaction(&btcdbtcjson.DecodeRawTransactionCmd{HexTx: signedRawHex}, &chaincfg.RegressionNetParams)
			// rq.Nil(err, "failed to decode signed raw tx")
			// b, _ := json.MarshalIndent(&decodeTx, "", " ")
			// fmt.Println("signed tx(token: a0 > multisig address)", string(b))

			// 广播交易
			txid, err := cli.Sendrawtransaction(btcjson.SendRawTransactionCmd{
				HexTx: signedRawHex,
			})
			assert.Nil(t, err)
			fmt.Println("broadcasted txid", txid)
		}
		{ // 生成一个块确认代币转账
			_, err = cli.Generatetoaddress(1, a0.Address, nil)
			rq.Nil(err, "Failed to generate to address")
		}

		{ // 确认代币转账成功
			bal, err := cli.OmniGetbalance(multisigAddress, propertyID)
			rq.Nil(err, "Failed to get omni balance")
			expectedBal := strconv.FormatFloat(transferAmount, 'f', 0, 32)
			rq.False(bal.Balance != expectedBal, "wrong balance, not ", expectedBal)
		}
	}
	{ // simple send from multisig address to a4
		transferAmount := float64(23)
		{
			var utxo btcjson.ListUnspentResult
			omnicli.NoPrintCmd(func() {
				unspents, err := cli.Listunspent(0, 999, []string{multisigAddress})
				rq.Nil(err, "Failed to list unspent")
				rq.False(len(unspents) == 0, "no unspent find")
				for _, u := range unspents {
					if u.Amount > 0.001 {
						utxo = u
						break
					}
				}
			})

			unspent := new(btc.BTCUnspent)
			unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, redeemScript)
			changeAddr, err := btc.NewBTCAddressFromString(multisigAddress, chainID)
			rq.Nil(err, "failed to create addr")
			toAddr, err := btc.NewBTCAddressFromString(a4.Address, chainID)
			rq.Nil(err, "failed to create addr")
			feeRate := int64(80)
			btctx, err := omni.CreateTransactionForOmni(propertyID, false, unspent, toAddr, transferAmount, changeAddr, feeRate, chainID)
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
				txid, err := cli.Sendrawtransaction(btcjson.SendRawTransactionCmd{
					HexTx: signedRawHex,
				})
				assert.Nil(t, err)
				fmt.Println("broadcasted txid", txid)
			}

		}
		{ // 生成一个块确认代币转账
			_, err = cli.Generatetoaddress(1, a0.Address, nil)
			rq.Nil(err, "Failed to generate to address")
		}

		{ // 确认代币转账成功
			bal, err := cli.OmniGetbalance(multisigAddress, propertyID)
			rq.Nil(err, "Failed to get omni balance")
			expectedBal := strconv.FormatFloat(233-23, 'f', 0, 32)
			rq.False(bal.Balance != expectedBal, "wrong balance, not ", expectedBal)
		}
		{ // 确认代币转账成功
			bal, err := cli.OmniGetbalance(a4.Address, propertyID)
			rq.Nil(err, "Failed to get omni balance")
			expectedBal := strconv.FormatFloat(23, 'f', 0, 32)
			rq.False(bal.Balance != expectedBal, "wrong balance, not ", expectedBal)
		}
	}

}
