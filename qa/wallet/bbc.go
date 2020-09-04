package wallet

import (
	"testing"

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
}
