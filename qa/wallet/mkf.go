package wallet

import (
	"testing"

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
}
