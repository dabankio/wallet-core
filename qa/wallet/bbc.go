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

func testBBCPubkSign(t *testing.T, w *wallet.Wallet, c ctx) {
	const pass = "123"
	r := require.New(t)
	nodeInfo := devtools4chains.MustRunDockerDevCore(t, bbcCoreImage, true, true)

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

	t.Run("使用RPC创建交易", func(t *testing.T) {
		rawTX, err := jsonRPC.Createtransaction(bbrpc.CmdCreatetransaction{
			From: address, To: minerAddress, Amount: outAmount,
		})
		r.NoError(err)
		signedTX, err := w.Sign(bbc.SymbolBBC, *rawTX)
		r.NoError(err)
		_, err = jsonRPC.Sendtransaction(signedTX)
		r.NoError(err)
		r.NoError(bbrpc.Wait4nBlocks(1, jsonRPC))
		bal, err := jsonRPC.Getbalance(nil, &address)
		r.NoError(err)
		r.Len(bal, 1)
		r.InDelta(bal[0].Avail, registeredAssets-outAmount-0.01, 0.00001)
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
			SetAnchor(forks[0].Fork).
			SetFee(0.01).
			SetVersion(0).
			SetTimestamp(int(time.Now().Unix())).
			AddInput(utxo.Txid, int8(utxo.Out)).
			Build()
		r.NoError(err)

		sig, err := w.Sign("BBC", tx)
		r.NoError(err)

		txid, err := jsonRPC.Sendtransaction(sig)
		r.NoError(err)
		t.Log("txid:", *txid)
	})

}
