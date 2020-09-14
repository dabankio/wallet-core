package eth

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/dabankio/devtools4chains"
	"github.com/dabankio/wallet-core/bip39"
	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/core/eth"
	"github.com/dabankio/wallet-core/core/eth/internalized/testtool"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
)

// 简单的ETH转账测试（示例）
func TestSimpleTransfer(t *testing.T) {
	rq := require.New(t)

	type addressInfo struct {
		address, privateKey, publicKey, mnemonic string
	}

	addrs := new([2]addressInfo)

	//bip44推导私钥
	for i := range addrs {
		var seed []byte
		var mnemonic string
		{ //bip39 生成种子
			bip39.SetWordListLang(bip39.LangChineseSimplified)
			ent, err := bip39.NewEntropy(128)
			rq.Nil(err)

			mnemonic, err = bip39.NewMnemonic(ent)
			rq.Nil(err)
			fmt.Println("mnemonic:", mnemonic)

			seed, err = bip39.NewSeedWithErrorChecking(mnemonic, "")
			rq.Nil(err)
		}
		deriver, err := eth.NewBip44Deriver(bip44.FullPathFormat, seed)
		rq.Nil(err)

		address, err := deriver.DeriveAddress()
		rq.Nil(err)
		privateKey, err := deriver.DerivePrivateKey()
		rq.Nil(err)
		publicKey, err := deriver.DerivePublicKey()
		rq.Nil(err)
		fmt.Printf("address/private/public: \n%s\n%s\n%s\n", address, privateKey, publicKey)
		addrs[i] = addressInfo{
			address:    address,
			privateKey: "0x" + privateKey,
			publicKey:  publicKey,
			mnemonic:   mnemonic,
		}
	}

	//启动本地测试节点
	killFunc, port, err := devtools4chains.DockerRunGanacheCli(&devtools4chains.DockerRunOptions{
		AutoRemove: true,
	})
	t.Cleanup(killFunc)
	var rpcHost = fmt.Sprintf("http://localhost:%d", port)
	var rpcClient *ethclient.Client
	{ //rpc client
		rpcClient, err = ethclient.Dial(rpcHost)
		rq.Nil(err, "dial failed")
		testtool.WaitSomething(t, time.Minute, func() error { _, e := rpcClient.NetworkID(context.Background()); return e })
	}

	PrepareFunds4address(t, rpcHost, addrs[0].address, 50)
	//首先确定addrs[0]里确实有50ETH, addrs[1]里为0
	rightBalances := []struct {
		address string
		balance *big.Int
	}{
		{address: addrs[0].address, balance: big.NewInt(0).Mul(big.NewInt(E18), big.NewInt(50))},
		{address: addrs[1].address, balance: big.NewInt(0)},
	}

	for _, rightBal := range rightBalances {
		bal, err := rpcClient.BalanceAt(context.Background(), common.HexToAddress(rightBal.address), nil)
		rq.Nil(err, "无法获取地址的余额")
		fmt.Println("addr余额", rightBal.address, bal)
		rq.Equal(bal.String(), rightBal.balance.String(), "Wrong balance")
		// rq.True(bal.Cmp(rightBal.balance) == 0, "Wrong balance")
	}

	gasLimit := 6721975
	gasPrice := 20000000000
	transferAmount := 6 * E18
	{ //现在使用sdk,从0转账{transferAmount}eth到1
		a0Nonce, err := rpcClient.PendingNonceAt(context.Background(), common.HexToAddress(addrs[0].address))
		rq.Nil(err)

		toA1Address, err := eth.NewETHAddressFromHex(addrs[1].address)
		rq.Nil(err)
		tx := eth.NewETHTransaction(int64(a0Nonce), toA1Address, eth.NewBigInt(int64(transferAmount)), int64(gasLimit), eth.NewBigInt(int64(gasPrice)), nil)

		rawTx, err := tx.EncodeRLP()
		rq.Nil(err)
		signedHex, err := eth.SignRawTransaction(rawTx, addrs[0].privateKey)
		rq.Nil(err)

		{ //广播交易， 一般步骤下，应该用签名好的数据调用广播api,这里用了geth,处理下数据然后调用geth jsonrpc 进行广播
			sigBytes, err := hexutil.Decode(signedHex)
			rq.Nil(err, "Failed to decode hexed sig ")
			var rawTx types.Transaction
			err = rlp.DecodeBytes(sigBytes, &rawTx)
			rq.Nil(err, "Failed to rlp decode tx")

			err = rpcClient.SendTransaction(context.Background(), &rawTx)
			rq.Nil(err, "Failed to send tx")
		}

	}

	{ //确认到账, 首先确定addrs[0]里确实有100ETH, addrs[1]里为0
		fmt.Println("Get balance --------")
		rightBalances := []struct {
			address string
			balance *big.Int
		}{
			{address: addrs[0].address, balance: nil},
			{address: addrs[1].address, balance: big.NewInt(int64(transferAmount))},
		}

		for _, rightBal := range rightBalances {
			bal, err := rpcClient.BalanceAt(context.Background(), common.HexToAddress(rightBal.address), nil)
			rq.Nil(err, "无法获取地址的余额")
			fmt.Println("addr余额", rightBal.address, bal)
			if rightBal.balance != nil {
				rq.Equal(bal.String(), rightBal.balance.String(), "Wrong balance")
				// rq.True(bal.Cmp(rightBal.balance) == 0, "Wrong balance")
			}
		}
	}

}
