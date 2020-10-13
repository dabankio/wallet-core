package trx

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	r "github.com/stretchr/testify/require"
)

/*
docker run -d -it -p 9097:9090 -p 50051:50051 -p 50052:50052 --rm --name tron_ci \
	-e "mnemonic=mirror increase slot auto memory bicycle flip latin correct humble private online" \
	-e "defaultBalance=1000000" -e "formatJson=true" \
	trontools/quickstart

curl http://127.0.0.1:9097/admin/accounts?format=all
*/
func TestTRX_TRC20(t *testing.T) {
	fn := func(addr, privk, apiAddress, _ string) {
		grpc := client.NewGrpcClientWithTimeout(apiAddress, 5*time.Second)
		r.NoError(t, grpc.Start())
		defer grpc.Stop()

		nodeInfo, _ := grpc.GetNodeInfo()
		fmt.Println("node block", nodeInfo.Block)

		acct, err := grpc.GetAccount(addr)
		r.NoError(t, err)
		r.NotNil(t, acct)
		fmt.Println("TRX balance", acct.Balance, acct.Balance/1e6)
		fmt.Println("eng use ", acct.AccountResource.EnergyUsage)

		t.Run("TRX转账", func(t *testing.T) {
			fmt.Println("-------------TRX 转账")
			ext, err := grpc.Transfer(addr, "TUU9bqm9CCA1dAU9iaa6HcJF4twMBi5N86", 12000)
			r.NotNil(t, acct)

			_ = SignBroadcastTx(t, ext, privk, grpc)

			time.Sleep(2 * time.Second)
			acct2, err := grpc.GetAccount(addr)
			r.NoError(t, err)
			fmt.Println(acct2.Balance)
			r.Greater(t, acct.Balance, acct2.Balance)
		})

		if acct.AccountResource.EnergyUsage < 100 {
			fmt.Println("-------------will freeze balance")
			ext, err := grpc.FreezeBalance(addr, "", core.ResourceCode_ENERGY, 500*1e6)
			r.NoError(t, err)
			_ = SignBroadcastTx(t, ext, privk, grpc)
			fmt.Println("will sleep, wait freeze balance")
			time.Sleep(2 * time.Second)
		}

		var trc20Addr string
		var baseTrc20Supply *big.Int
		{ //deploy trc20 contract
			fmt.Println("-------------will deploy trc20 token")
			dptx, err := grpc.DeployContract(
				addr, "trc20xxx", nil, TestTRC20ContractBytecode, 1e9, 20, 9e10,
			)
			r.NoError(t, err)

			dpTxid := SignBroadcastTx(t, dptx, privk, grpc)

			fmt.Println("will sleep, wait for token contract")
			time.Sleep(3 * time.Second)

			txinfo, err := grpc.GetTransactionInfoByID(dpTxid)
			r.NoError(t, err)

			trc20Addr = address.Address(txinfo.GetContractAddress()).String()
			fmt.Println("token addr", trc20Addr)

			baseTrc20Supply, err = grpc.TRC20ContractBalance(addr, trc20Addr)
			r.NoError(t, err)
			fmt.Println("baseTrc20Supply", baseTrc20Supply.String())
		}

		sendAmount := big.NewInt(50)
		t.Run("trc20 transfer", func(t *testing.T) { //
			fmt.Println("-------------will send trc20 token")
			to := "TSGYZ3VAVsa2SgoYhV7mfqnde89zTU7zNh"
			ext, err := grpc.TRC20Send(addr, to, trc20Addr, sendAmount, 999999)
			r.NoError(t, err)

			txid := SignBroadcastTx(t, ext, privk, grpc)
			fmt.Println("will sleep, wait for trc20 send tx")
			time.Sleep(3 * time.Second)

			// curl --request GET --url  http://127.0.0.1:9097/v1/transactions/txid/events
			txinfo, err := grpc.GetTransactionInfoByID(txid)
			r.NoError(t, err)

			fmt.Println("result:", txinfo.Result)
			fmt.Println("block number", txinfo.BlockNumber)
			fmt.Println("fee", txinfo.Fee)
			fmt.Println("eng", txinfo.Receipt.EnergyUsageTotal)
			fmt.Println("rec result:", txinfo.Receipt.Result)

			bal, err := grpc.TRC20ContractBalance(addr, trc20Addr)
			r.NoError(t, err)
			fmt.Println("bal", bal.String())
			r.Equal(t, big.NewInt(0).Sub(baseTrc20Supply, sendAmount).Text(10), bal.Text(10))
		})
	}
	testIn(fn)
}

func testIn(fn func(address, privk, apiAddress, httpAPI string)) {
	const api = "192.168.50.5:50051"
	const httpAPI = "http://192.168.50.5:9097"

	const ( //bind to docker env now , see test doc
		addr  = "TMJTFYx6oQVKLkn5pMmsegsvgqKUoYnaEB"
		privk = "1f8caad99870f73ed96a777b05d79c230197f9a9f42a5c3511a5d0f5101c47b7"
	)

	fn(addr, privk, api, httpAPI)
}
