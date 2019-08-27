package internal

import (
	"fmt"
	"math/big"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lomocoin/wallet-core/core/eth/internal/contracts"
	"github.com/lomocoin/wallet-core/core/eth/internal/testtool"
)

// some const
const (
	LocalRPCHost7545 = "http://localhost:7545"
)

// 本测试要求本地7545端口运行有ganache,并且至少有一个账号余额超过5eth
// 测试简单的erc20 (并不使用erc20,而是用来测试多签支持)
func TestErc20(t *testing.T) {
	var (
		rpcClient      *ethclient.Client
		contract       *contracts.FixedSupplyToken
		a0, a1, a2, a3 *testtool.AddrInfo
		addrs          []*testtool.AddrInfo
		err            error
	)
	{
		_, _, _ = a1, a2, a3
	}

	{ //生成4个地址，并排序
		for i := 0; i < 4; i++ {
			addr := testtool.GenAddr()
			addrs = append(addrs, addr)
		}
		sort.Slice(addrs, func(i, j int) bool {
			return addrs[i].Address < addrs[j].Address
		})
		a0, a1, a2, a3 = addrs[0], addrs[1], addrs[2], addrs[3]
	}

	{ // init vars
		rpcClient, err = ethclient.Dial(LocalRPCHost7545)
		testtool.FailOnErr(t, err, "Failed to dial rpc")
	}

	{ // 为操作的账号准备些手续费
		for _, add := range addrs {
			testtool.PrepareFunds4address(t, LocalRPCHost7545, add.Address, 1)
		}
	}

	{ // 部署erc20 合约, owner 为 a0
		auth := bind.NewKeyedTransactor(a0.ToECDSAKey())
		_, _, contrakt, err := contracts.DeployFixedSupplyToken(auth, rpcClient)
		testtool.FailOnErr(t, err, "Failed to deploy erc20 contract")
		contract = contrakt
	}

	{ // 首先测试查询数据
		bal, e := contract.BalanceOf(&bind.CallOpts{}, a0.ToAddress())
		testtool.FailOnErr(t, e, "Failed to get balance of owner")
		fmt.Println("balance of erc20 owner", bal)
		shouldBe, flag := big.NewInt(1).SetString("1000000000000000000000000", 10)
		testtool.FailOnFlag(t, !flag, "set int failed")
		testtool.FailOnFlag(t, bal.Cmp(shouldBe) != 0, "余额异常")
	}

	{ // erc20转账, a0 -> a1 , testtool.E18
		auth := bind.NewKeyedTransactor(a0.ToECDSAKey())
		tx, err := contract.Transfer(auth, a1.ToAddress(), big.NewInt(testtool.E18))
		testtool.FailOnErr(t, err, "Failed to do erc20 transfer")
		fmt.Println("transfer tx", tx)
	}

	{ // 查询转账后a1的余额
		bal, err := contract.BalanceOf(&bind.CallOpts{}, a1.ToAddress())
		testtool.FailOnErr(t, err, "Failed to get balance of erc20(a1)")
		fmt.Println("erc20 balance of a1", bal)
		shouldBe := big.NewInt(testtool.E18)
		testtool.FailOnFlag(t, shouldBe.Cmp(bal) != 0, "wrong erc20 balance")
	}

}
