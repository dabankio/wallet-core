package wallet

import (
	"testing"
	"time"

	"github.com/dabankio/bbrpc"
	"github.com/dabankio/devtools4chains"
	"github.com/dabankio/wallet-core/core/bbc"
	"github.com/dabankio/wallet-core/wallet"
	"github.com/stretchr/testify/require"
)

func testMKFPubkSign(t *testing.T, w *wallet.Wallet, c ctx) {
	r := require.New(t)
	const pass = "123"
	nodeInfo := devtools4chains.MustRunDockerMKFDev(t, mkfImage, true, true)

	jsonRPC := nodeInfo.Client
	minerAddress := nodeInfo.MinerAddress

	pubk, address := c.pubk, c.address
	var err error

	registeredAssets := 12.34
	{ // 导入公钥
		_, err = jsonRPC.Importpubkey(pubk) // <<=== RPC 导入公钥
		r.NoError(err)
		r.NoError(bbrpc.Wait4balanceReach(minerAddress, 10, jsonRPC))
		jsonRPC.Unlockkey(nodeInfo.MinerOwnerPubk, nodeInfo.UnlockPass, nil)
		_, err = jsonRPC.Sendfrom(bbrpc.CmdSendfrom{
			From: minerAddress, To: address, Amount: registeredAssets,
		})
		r.NoError(err)
		r.NoError(bbrpc.Wait4balanceReach(address, registeredAssets, jsonRPC))
	}

	outAmount := 2.3

	t.Run("用RPC创建交易", func(t *testing.T) {

		//创建交易、签名、广播、检查余额
		rawTX, err := jsonRPC.Createtransaction(bbrpc.CmdCreatetransaction{ // <<=== RPC 创建交易
			From: address, To: minerAddress, Amount: outAmount,
		})
		r.NoError(err)

		// fmt.Println("rawTX:", *rawTX)
		// deTx, err := bbc.DecodeSymbolTX("MKF", *rawTX) // <<=== sdk 反序列化交易
		// r.NoError(err)
		// fmt.Println("decoded tx", deTx) //decoded tx {"Version":1,"Typ":0,"Timestamp":1584952846,"LockUntil":0,"SizeIn":1,"Prefix":2,"Amount":1340000,"TxFee":100,"SizeOut":0,"SizeSign":0,"HashAnchor":"00000000c335f935650a427bf548242eac4e4a444e25691b47351e7945f4a8d4","Address":"10g06z2bmwb71n9xg9zsv4vzay86ab7avt6n97hm6ra2z3rsbrtc2ncer","Sign":""}

		signedTX, err := w.Sign(bbc.SymbolMKF, *rawTX)
		r.NoError(err)

		_, err = jsonRPC.Sendtransaction(signedTX) // <<=== RPC 发送交易
		r.NoError(err)

		r.NoError(bbrpc.Wait4nBlocks(1, jsonRPC))

		bal, err := jsonRPC.Getbalance(nil, &address) // <<=== RPC 查询余额
		r.NoError(err)
		r.Len(bal, 1)
		r.InDelta(bal[0].Avail, registeredAssets-outAmount-0.03, 0.00001)
		// fmt.Println("balance after send", bal[0]) //balance after send {1dmyvkbkbk5zaqvx46zqpy2vzywjz02sv5kdd0gq2c56mwb48925hfhpd 0.9899 0 0}
	})

	t.Run("使用SDK创建交易", func(t *testing.T) {
		toAddress := "1sxs9gnbxs7nfb0m4xrmwkw3ew1dzg9hmv56e09dndd1rt0bbqwy9f6gv"

		forks, err := jsonRPC.Listfork(true)
		r.NoError(err)

		unspents, err := jsonRPC.Listunspent(c.address, nil, 999)
		r.NoError(err)
		utxo := unspents.Addresses[0].Unspents[0]

		tb := bbc.NewTxBuilder()
		tx, err := tb.
			SetAddress(toAddress).
			SetAmount(2.1).
			SetAnchor(forks[0].Fork). // <== MKF 不需要这个(但还是设置下，sdk内部会检查forkID,这里可以随便填个64长度的hex值)
			ExcludeAnchor().          //<== MKF 要这个, 否则会编码成BBC交易
			SetFee(0.03).
			SetVersion(0).
			SetTimestamp(int(time.Now().Unix())).
			AddInput(utxo.Txid, int8(utxo.Out)).
			Build()
		r.NoError(err)

		sig, err := w.Sign("MKF", tx)
		r.NoError(err)

		txid, err := jsonRPC.Sendtransaction(sig)
		r.NoError(err)
		t.Log("txid:", *txid)
	})

}

func testMKFDexTXSign(t *testing.T, w *wallet.Wallet, c ctx) {
	r := require.New(t)
	const pass = "123"
	nodeInfo := devtools4chains.MustRunDockerMKFDev(t, mkfImageDexTest, true, true)

	jsonRPC := nodeInfo.Client
	// minerAddress := nodeInfo.MinerAddress
	minerAddress := "20g003rgxdn4s64r4d0dchvb87p791q4epswkn1txadgv1evjqqwv70e5"

	pubk, address := c.pubk, c.address
	_ = pubk
	var err error

	registeredAssets := 12.34
	{ // 导入公钥
		privk, err := w.DerivePrivateKey("MKF")
		r.NoError(err)
		_, err = jsonRPC.Importprivkey(privk, pass) // <<=== RPC 导入公钥
		r.NoError(err)
		_, err = jsonRPC.Unlockkey(pubk, pass, nil)
		r.NoError(err)

		r.NoError(bbrpc.Wait4balanceReach(minerAddress, registeredAssets+1, jsonRPC))
		jsonRPC.Unlockkey(nodeInfo.MinerOwnerPubk, nodeInfo.UnlockPass, nil)
		f := 0.03
		_, err = jsonRPC.Sendfrom(bbrpc.CmdSendfrom{
			From: minerAddress, To: address, Amount: registeredAssets, Txfee: &f,
		})
		r.NoError(err)
		r.NoError(bbrpc.Wait4balanceReach(address, registeredAssets, jsonRPC))
	}

	t.Run("使用SDK创建交易,签名结果应该与rpc签名一致", func(t *testing.T) {
		// 确保创建的模版id, data一致
		// offline签名结果一致π
		// 向挂单地址转账
		// 从挂单地址向撮合地址转账
		cmd := map[string]interface{}{
			"type": "dexorder",
			"dexorder": map[string]interface{}{
				"seller_address": c.address,
				"coinpair":       "bbc/mkf",
				"price":          10,
				"fee":            0.002,
				"recv_address":   c.address,
				"valid_height":   300,
				"match_address":  "15cx56x0gtv44bkt21yryg4m6nn81wtc7gkf6c9vwpvq1cgmm8jm7m5kd",
				"deal_address":   "1f2b2n3asbm2rb99fk1c4wp069d0z91enxdz8kmqmq7f0w8tzw64hdevb",
			}}

		var orderID string
		_, err = jsonRPC.CallJSONRPC("addnewtemplate", cmd, &orderID)
		require.NoError(t, err)
		t.Log("orderID:", orderID)

		tplID, err := bbc.CreateTemplateDataDexOrder(
			c.address,
			"bbc/mkf",
			10_00000_00000, //10*10e10
			20,             //0.002*4
			c.address,
			300,
			"15cx56x0gtv44bkt21yryg4m6nn81wtc7gkf6c9vwpvq1cgmm8jm7m5kd",
			"1f2b2n3asbm2rb99fk1c4wp069d0z91enxdz8kmqmq7f0w8tzw64hdevb",
		)
		require.NoError(t, err)
		require.Equal(t, orderID, tplID.Address)

		amt := 2.2
		fee := 0.03
		createTime := time.Now()
		raw, err := jsonRPC.Createtransaction(bbrpc.CmdCreatetransaction{
			From:   c.address,
			To:     tplID.Address,
			Amount: amt,
			Txfee:  &fee,
		})
		require.NoError(t, err)

		rpcSig, err := jsonRPC.Signrawtransactionwithwallet(c.address, *raw)
		require.NoError(t, err)

		forks, err := jsonRPC.Listfork(true)
		require.NoError(t, err)

		unspents, err := jsonRPC.Listunspent(c.address, nil, 999)
		require.NoError(t, err)
		utxo := unspents.Addresses[0].Unspents[0]

		tb := bbc.NewTxBuilder()
		tx, err := tb.
			SetAddress(tplID.Address).
			SetAmount(amt).
			SetAnchor(forks[0].Fork). // <== MKF 不需要这个(但还是设置下，sdk内部会检查forkID,这里可以随便填个64长度的hex值)
			ExcludeAnchor().          //<== MKF 要这个, 否则会编码成BBC交易
			SetFee(fee).
			SetVersion(2).
			SetTimestamp(int(createTime.Unix())).
			AddInput(utxo.Txid, int8(utxo.Out)).
			AddTemplateData(tplID.RawHex).
			Build()
		require.NoError(t, err)

		sig, err := w.Sign("MKF", tx)
		require.NoError(t, err)

		require.Equal(t, rpcSig.Hex, sig)

		txid, err := jsonRPC.Sendtransaction(sig)
		require.NoError(t, err)
		t.Log("txid:", *txid)
	})

}
