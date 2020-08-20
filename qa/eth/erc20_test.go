package eth

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/dabankio/devtools4chains"
	"github.com/dabankio/wallet-core/core/eth/internalized/testtool"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 本测试要求本地7545端口运行有ganache,并且至少有一个账号余额超过5eth
// 测试简单的erc20 (并不使用erc20,而是用来测试多签支持)
func TestErc20(t *testing.T) {
	killFunc, port, err := devtools4chains.DockerRunGanacheCli(&devtools4chains.DockerRunOptions{
		AutoRemove: true,
	})
	t.Cleanup(killFunc)

	var ganacheRPCHost = fmt.Sprintf("http://localhost:%d", port)
	rq := require.New(t)

	var (
		rpcClient *ethclient.Client
		contract  *FixedSupplyToken
		a0, a1    *AddrInfo
		addrs     []*AddrInfo
	)

	{ //生成4个地址，并排序
		for i := 0; i < 4; i++ {
			addr := GenAddr()
			addrs = append(addrs, addr)
		}
		sort.Slice(addrs, func(i, j int) bool {
			return addrs[i].Address < addrs[j].Address
		})
		a0, a1 = addrs[0], addrs[1]
	}

	{ // init vars
		rpcClient, err = ethclient.Dial(ganacheRPCHost)
		rq.Nil(err, "Failed to dial rpc")
		testtool.WaitSomething(t, time.Minute, func() error { _, e := rpcClient.NetworkID(context.Background()); return e })
	}

	// 为操作的账号准备些手续费
	for _, add := range addrs {
		PrepareFunds4address(t, ganacheRPCHost, add.Address, 1)
	}

	{ // 部署erc20 合约, owner 为 a0
		auth := bind.NewKeyedTransactor(a0.ToECDSAKey())
		_, _, contrakt, err := DeployFixedSupplyToken(auth, rpcClient)
		rq.Nil(err, "Failed to deploy erc20 contract")
		contract = contrakt
	}

	{ // 首先测试查询数据
		bal, e := contract.BalanceOf(&bind.CallOpts{}, a0.ToAddress())
		rq.Nil(e, "Failed to get balance of owner")
		fmt.Println("balance of erc20 owner", bal)
		shouldBe, flag := big.NewInt(1).SetString("1000000000000000000000000", 10)
		rq.False(!flag, "set int failed")
		rq.False(bal.Cmp(shouldBe) != 0, "余额异常")
	}

	{ // erc20转账, a0 -> a1 , E18
		auth := bind.NewKeyedTransactor(a0.ToECDSAKey())
		tx, err := contract.Transfer(auth, a1.ToAddress(), big.NewInt(E18))
		rq.Nil(err, "Failed to do erc20 transfer")
		fmt.Println("transfer tx", tx)
	}

	{ // 查询转账后a1的余额
		bal, err := contract.BalanceOf(&bind.CallOpts{}, a1.ToAddress())
		rq.Nil(err, "Failed to get balance of erc20(a1)")
		fmt.Println("erc20 balance of a1", bal)
		shouldBe := big.NewInt(E18)
		rq.False(shouldBe.Cmp(bal) != 0, "wrong erc20 balance")
	}

}
