package wallet

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/dabankio/devtools4chains"
	"github.com/dabankio/wallet-core/core/eth"
	"github.com/dabankio/wallet-core/core/eth/internalized/testtool"
	qaEth "github.com/dabankio/wallet-core/qa/eth"
	"github.com/dabankio/wallet-core/wallet"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const E18 = 1e18

func testETHPubkSign(t *testing.T, w *wallet.Wallet, c ctx) {
	const toAddress = "0x9461d8c5D4d7567E889eAB278851ce6556Ce05d9" //随便找的地址

	rq := require.New(t)
	killFunc, port, err := devtools4chains.DockerRunGanacheCli(&devtools4chains.DockerRunOptions{
		AutoRemove: true,
	})
	t.Cleanup(killFunc)
	var rpcHost = fmt.Sprintf("http://localhost:%d", port)
	var rpcClient *ethclient.Client

	rpcInfo := devtools4chains.RPCInfo{
		Host: fmt.Sprintf("localhost:%d", port),
	}
	{ //rpc client
		rpcClient, err = ethclient.Dial(rpcHost)
		rq.Nil(err, "dial failed")
		testtool.WaitSomething(t, time.Minute, func() error { _, e := rpcClient.NetworkID(context.Background()); return e })
	}

	qaEth.PrepareFunds4address(t, rpcHost, c.address, 50)

	outAmountInt := int64(E18 + 1e17)
	outAmount := big.NewInt(outAmountInt)
	{ //prepare tx
		var (
			nonce    uint64
			gasLimit uint64 = 21000 //ETH fixed gasLimit
			gasPrice *big.Int
		)

		nonce, err = rpcClient.PendingNonceAt(context.Background(), common.HexToAddress(c.address))
		rq.NoError(err)
		gasPrice, err = rpcClient.SuggestGasPrice(context.Background())
		rq.NoError(err)

		ethToAddress, err := eth.NewETHAddressFromHex(toAddress)
		rq.NoError(err)

		tx := eth.NewETHTransaction(int64(nonce), ethToAddress, eth.NewBigInt(outAmountInt), int64(gasLimit), eth.NewBigInt(gasPrice.Int64()), nil)

		{ //等同于
			tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), outAmount, gasLimit, gasPrice, nil)

			buf := bytes.NewBuffer(nil)
			rq.NoError(tx.EncodeRLP(buf))
			toSignData := "0x" + hex.EncodeToString(buf.Bytes())
			_ = toSignData
		}
		toSignData, err := tx.EncodeRLP() //编码为可签名/传输数据
		rq.NoError(err)

		sig, err := w.Sign("ETH", toSignData) //签名
		rq.NoError(err)

		resp, err := devtools4chains.RPCCallJSON(rpcInfo, "eth_sendRawTransaction", []string{sig}, nil)
		rq.NoError(err)
		fmt.Println("send tx resp:", string(resp))
	}

	balance, err := rpcClient.BalanceAt(context.Background(), common.HexToAddress(toAddress), nil)
	rq.NoError(err)
	fmt.Println("balance fo to address", balance)
	rq.Equal(outAmount.String(), balance.String(), "wrong balance")
}

func testERC20PubkSign(t *testing.T, w *wallet.Wallet, c ctx) {
	const toAddress = "0x9461d8c5D4d7567E889eAB278851ce6556Ce05d9" //随便找的地址

	rq := require.New(t)
	killFunc, port, err := devtools4chains.DockerRunGanacheCli(&devtools4chains.DockerRunOptions{
		AutoRemove: true,
	})
	_ = killFunc
	// t.Cleanup(killFunc)
	var rpcHost = fmt.Sprintf("http://localhost:%d", port)
	var rpcClient *ethclient.Client

	rpcInfo := devtools4chains.RPCInfo{
		Host: fmt.Sprintf("localhost:%d", port),
	}
	{ //rpc client
		rpcClient, err = ethclient.Dial(rpcHost)
		rq.Nil(err, "dial failed")
		testtool.WaitSomething(t, time.Minute, func() error { _, e := rpcClient.NetworkID(context.Background()); return e })
	}

	qaEth.PrepareFunds4address(t, rpcHost, c.address, 50)

	privk, err := w.DerivePrivateKey("ETH")
	rq.NoError(err)

	//部署ERC20
	privateKey, err := crypto.HexToECDSA(privk)
	rq.NoError(err)
	auth := bind.NewKeyedTransactor(privateKey)
	erc20ContractAddress, _, erc20Contrakt, err := qaEth.DeployFixedSupplyToken(auth, rpcClient)
	rq.Nil(err, "Failed to deploy erc20 contract")

	erc20ABI := eth.NewERC20InterfaceABIHelper()
	wrapToEthAddress, err := eth.NewETHAddressFromHex(toAddress)
	rq.NoError(err)

	packedData, err := erc20ABI.PackedTransfer(wrapToEthAddress, eth.NewBigInt(2*E18+1e17))
	rq.NoError(err)

	var (
		nonce    uint64
		gasLimit uint64 = 3_000 //随便写的值，需要根据实际情况取值
		gasPrice *big.Int
	)

	nonce, err = rpcClient.PendingNonceAt(context.Background(), common.HexToAddress(c.address))
	rq.NoError(err)
	gasPrice, err = rpcClient.SuggestGasPrice(context.Background())
	rq.NoError(err)

	ethToAddress := eth.NewETHAddress()
	rq.NoError(ethToAddress.SetHex(erc20ContractAddress.Hex()))

	tx := eth.NewETHTransaction(int64(nonce), ethToAddress, eth.NewBigInt(0), int64(gasLimit), eth.NewBigInt(gasPrice.Int64()), packedData)
	toSignData, err := tx.EncodeRLP()
	rq.NoError(err)

	{ //tx创建等同于
		tx := types.NewTransaction(nonce, erc20ContractAddress, new(big.Int), gasLimit, gasPrice, packedData)

		buf := bytes.NewBuffer(nil)
		rq.NoError(tx.EncodeRLP(buf))
		toSign := "0x" + hex.EncodeToString(buf.Bytes())
		_ = toSign
	}
	sig, err := w.Sign("ETH", toSignData)
	rq.NoError(err)
	fmt.Println("packed", hex.EncodeToString(packedData))
	fmt.Println("toSignData", toSignData)
	fmt.Println("sig   ", sig)

	resp, err := devtools4chains.RPCCallJSON(rpcInfo, "eth_sendRawTransaction", []string{sig}, nil)
	rq.NoError(err)
	fmt.Println("send tx resp:", string(resp))

	trTx, err := erc20Contrakt.Transfer(auth, common.HexToAddress(toAddress), big.NewInt(2*E18+1e17))
	_ = trTx
	rq.NoError(err)

	for _, addr := range []struct {
		address, shouldBe string
	}{
		{c.address, "9999958"},
		{toAddress, "42"},
	} {
		balance, err := erc20Contrakt.BalanceOf(&bind.CallOpts{}, common.HexToAddress(addr.address))
		rq.NoError(err)
		b := balance.Div(balance, big.NewInt(1e17))
		fmt.Println("erc20 balance", b)
		rq.Equal(addr.shouldBe, b.String())
	}

}
