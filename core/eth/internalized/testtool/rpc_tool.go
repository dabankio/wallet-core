package testtool

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

// PrepareFunds4address 为addr准备一定量的eth,
func PrepareFunds4address(t *testing.T, rpcHost, addr string, funds int64) {
	rpcClient, err := rpc.DialContext(context.Background(), rpcHost)
	FailOnErr(t, err, "Failed to dial rpc")
	defer rpcClient.Close()

	client := ethclient.NewClient(rpcClient)

	var hexedAccounts []string
	err = rpcClient.Call(&hexedAccounts, "eth_accounts")
	// fmt.Println(hexedAccounts)
	FailOnErr(t, err, "Fail on get accounts")

	fromAccount := ""
	for _, acc := range hexedAccounts {
		bal, err := client.BalanceAt(context.Background(), common.HexToAddress(acc), nil)
		FailOnErr(t, err, "Failed to get balance of account")
		if bal.Cmp(big.NewInt(E18*funds+E18)) > 0 {
			fromAccount = acc
		}
	}

	if fromAccount == "" {
		t.Fatal("余额不足，无法充值")
	}

	tx := map[string]interface{}{
		"from": fromAccount,
		"to":   addr,
		// "gas": "0x76c0", // 30400
		"gasPrice": "0x9184e72a000", // 10000000000000
		"value":    E18 * funds,
	}
	var txHash string
	err = rpcClient.Call(&txHash, "eth_sendTransaction", tx)
	FailOnErr(t, err, "转账失败")
	fmt.Println("PrepareFunds4addr done", addr, funds)
}

func WaitSomething(t *testing.T, timeout time.Duration, fn func() error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			t.Fatalf("wait something timeout, %s", timeout)
		default:
			if e := fn(); e == nil {
				return
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

}
