// +build integration

package internal

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sort"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"

	"github.com/dabankio/wallet-core/core/eth/internal/contracts"
	"github.com/dabankio/wallet-core/core/eth/internal/testtool"
)

const chainID = 1

//在模拟链上测试多签（建议废弃，在测试的时候似乎合约这块不是很顺利），考虑使用ganache测试
// 使用模拟链时存在的问题：1.最后无法查询到转账后的余额，2.获取的合约日志也查不到execute log
// 如果模拟链可以完全替代ganache顺利测试，那么其实是比较适合作为单元测试的
func TestSimplemultisigSimulated(t *testing.T) {
	// t.SkipNow()

	var a0, a1, a2, a3 *testtool.AddrInfo
	var addrs []*testtool.AddrInfo
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
	alloc := core.GenesisAlloc{
		a0.ToAddress(): core.GenesisAccount{
			Balance:    big.NewInt(1).Mul(big.NewInt(testtool.E18), big.NewInt(100)),
			PrivateKey: crypto.FromECDSA(a0.ToECDSAKey()),
		},
	}
	sim := backends.NewSimulatedBackend(alloc, uint64(testtool.E18/100000))
	{ // 首先确定addr0里的余额
		bal, err := sim.BalanceAt(context.Background(), a0.ToAddress(), nil)
		testtool.FailOnErr(t, err, "无法获取地址的余额")
		fmt.Println("addr0余额", bal)
	}

	type args struct {
		privkHex   string
		hexAddress []string
		mRequired  uint8
	}
	var (
		contractAddressHex string
		contractAddress    common.Address
	)
	{ // 部署合约测试
		arg := args{
			privkHex:   a0.PrivkHex,
			hexAddress: []string{a0.Address, a1.Address, a2.Address},
			mRequired:  2,
		}

		chainID := big.NewInt(1)
		got, err := DeploySimpleMultiSigContract(*chainID, sim, arg.privkHex, arg.hexAddress, arg.mRequired)
		if err != nil {
			t.Errorf("DeployMultiSigWalletContract() error = %v", err)
			t.FailNow()
		}
		fmt.Println("deployMultisigWalletContract got:", got)

		contractAddressHex = got
		contractAddress = common.HexToAddress(contractAddressHex)
		fmt.Println("contractAddressHex", contractAddressHex)
	}

	sim.Commit()

	{ // 部署好后验证合约属性（owners/mRequired）
		owners, mRequired, err := GetContractInfo(sim, contractAddressHex)
		testtool.FailOnErr(t, err, "Failed to get contract info")
		fmt.Println("contract info", owners, mRequired)
		testtool.FailOnFlag(t, len(owners) != 3, "len owners != 3", len(owners))
		testtool.FailOnFlag(t, mRequired != 2, "mRequired != 2", mRequired)
	}

	{ //合约部署后往其中转入资金
		value := big.NewInt(1).Mul(big.NewInt(testtool.E18), big.NewInt(20))
		ctx := context.Background()
		nonce, err := sim.NonceAt(ctx, contractAddress, nil)
		testtool.FailOnErr(t, err, "Failed to get nonce")
		tx := types.NewTransaction(nonce, contractAddress, value, uint64(100000000), nil, nil)
		signer := types.MakeSigner(params.TestChainConfig, nil)
		signature, err := crypto.Sign(signer.Hash(tx).Bytes(), a0.ToECDSAKey())
		testtool.FailOnErr(t, err, "签名交易失败")
		tx, err = tx.WithSignature(signer, signature)
		testtool.FailOnErr(t, err, "为交易附加签名数据错误")
		err = sim.SendTransaction(ctx, tx)
		testtool.FailOnErr(t, err, "充值到合约地址异常")
		sim.Commit()

		bal, err := sim.BalanceAt(ctx, contractAddress, nil)
		testtool.FailOnErr(t, err, "无法获取合约地址的余额")
		testtool.FailOnFlag(t, bal.Cmp(value) != 0, fmt.Sprintf("合约地址的余额异常，应该是 %v, 实际上：%s", value, bal.String()))
		fmt.Println("合约地址当前余额", bal)

		bal, err = sim.BalanceAt(ctx, a0.ToAddress(), nil)
		testtool.FailOnErr(t, err, "无法获取地址的余额")
		fmt.Println("地址0当前余额", bal)
	}

	outAddr := a3.Address
	transferValue := big.NewInt(testtool.E18)
	{ // 交易测试
		var (
			sigV                                 []uint8    //签名
			sigR, sigS                           [][32]byte //签名
			privkHex                             string
			multisigContractAddress, fromAddress string //多签合约地址，发起地址
			destination, executor                string //toAddress
			value, gasLimit                      *big.Int
			data                                 []byte
			nonce                                uint64
		)
		// 012由0发起，0和2签名, 把钱赚到1的地址上,executor 为0
		privkHex = a0.PrivkHex
		multisigContractAddress = contractAddressHex
		fromAddress = a0.Address
		executor = a0.Address
		destination = outAddr
		value = transferValue
		gasLimit = big.NewInt(239963)
		data = []byte("")
		nonce, err := sim.PendingNonceAt(context.Background(), contractAddress)
		testtool.FailOnErr(t, err, "Failed to get contract nonce")

		for _, add := range []*testtool.AddrInfo{a0, a2} {
			v, r, s, err := SimpleMultiSigExecuteSign(chainID, add.PrivkHex, multisigContractAddress, destination, executor, nonce, value, gasLimit, data)
			testtool.FailOnErr(t, err, "create sig failed")
			sigV = append(sigV, v)
			sigR = append(sigR, r)
			sigS = append(sigS, s)
		}

		txid, err := ExecuteTX(&TxParams{
			backend:                 sim,
			sigV:                    sigV,
			sigR:                    sigR,
			sigS:                    sigS,
			privkHex:                privkHex,
			multisigContractAddress: multisigContractAddress,
			fromAddress:             fromAddress,
			destination:             destination,
			executor:                executor,
			value:                   value,
			gasLimit:                gasLimit,
			data:                    data,
		})
		testtool.FailOnErr(t, err, "Execute Failed")
		fmt.Println("execute txid", txid)

		sim.Commit()
	}
	sim.Commit()

	{ // 完了检查确实转账成功
		bal, err := sim.BalanceAt(context.Background(), common.HexToAddress(outAddr), nil)
		testtool.FailOnErr(t, err, "FonGetBal")
		fmt.Println("balance of new tx", bal)
		//TODO 检查余额这里暂时有点问题，在ganache上测试没问题，在模拟器上不能正确的查到，
		// testtool.FailOnFlag(t, bal.Cmp(transferValue) != 0, "余额不对", transferValue)
	}

	{ // 调试合约日志
		multisigContract, err := contracts.NewSimpleMultiSig(contractAddress, sim)
		testtool.FailOnErr(t, err, "构建多签合约调用时异常,检查合约地址和rpc server")
		// go func() {
		// 	ito, err := multisigContract.FilterDebugRecover(&bind.FilterOpts{Start: 0})
		// 	testtool.FailOnErr(t, err, "过滤合约日志异常")
		// 	for {
		// 		if !ito.Next() {
		// 			break
		// 		}
		// 		evt := ito.Event
		// 		log.Println("evt recoverd address:", evt.Recovered.Hex())
		// 	}
		// }()

		// go func() {
		// 	ito, err := multisigContract.FilterDbgExecuteParam(&bind.FilterOpts{Start: 0})
		// 	testtool.FailOnErr(t, err, "过滤合约日志异常")
		// 	for {
		// 		if !ito.Next() {
		// 			break
		// 		}
		// 		evt := ito.Event
		// 		log.Println("evt seperator:", hex.EncodeToString(evt.Sperator[:]))
		// 		log.Println("evt TxInputHash:", hex.EncodeToString(evt.TxInputHash[:]))
		// 		log.Println("evt TotalHash:", hex.EncodeToString(evt.TotalHash[:]))
		// 		log.Println("evt txInput:", hex.EncodeToString(evt.TxInput[:]))
		// 	}
		// }()
		executeFilter, err := multisigContract.FilterExecute(&bind.FilterOpts{Start: 0})
		testtool.FailOnErr(t, err, "过滤合约日志异常")
		for {
			if !executeFilter.Next() {
				break
			}
			evt := executeFilter.Event
			log.Println("evt execute, confirmAddrs:", evt.ConfirmAddrs)
			log.Println("evt execute, dest:", evt.Destination.Hex())
		}

		depositLogFilter, err := multisigContract.FilterDeposit(&bind.FilterOpts{Start: 0}, nil)
		testtool.FailOnErr(t, err, "过滤合约日志异常")
		for {
			if !depositLogFilter.Next() {
				break
			}
			evt := depositLogFilter.Event
			log.Println("evt deposit From:", evt.From.Hex())
			log.Println("evt deposit value:", evt.Value)
		}
	}
}

// 本测试要求本地7545端口运行有ganache,并且至少有一个账号余额超过5eth
// 测试erc20多签
func TestSimplemultisigGanacheERC20(t *testing.T) {
	// t.SkipNow()

	const (
		rpcHost = "http://localhost:7545"
	)
	var (
		a0, a1, a2, a3       *testtool.AddrInfo
		addrs                []*testtool.AddrInfo
		rpcClient            *ethclient.Client
		erc20Contract        *contracts.FixedSupplyToken
		erc20ContractAddress common.Address
		err                  error
	)
	{ //init vars
		for i := 0; i < 4; i++ {
			addr := testtool.GenAddr()
			addrs = append(addrs, addr)
		}
		sort.Slice(addrs, func(i, j int) bool {
			return addrs[i].Address < addrs[j].Address
		})
		a0, a1, a2, a3 = addrs[0], addrs[1], addrs[2], addrs[3]
		rpcClient, err = ethclient.Dial(rpcHost)
		testtool.FailOnErr(t, err, "dial failed")
	}

	{ //准备点eth做手续费
		testtool.PrepareFunds4address(t, rpcHost, a0.Address, 1)
	}
	{ // 首先确定addr0里的余额
		bal, err := rpcClient.BalanceAt(context.Background(), a0.ToAddress(), nil)
		testtool.FailOnErr(t, err, "无法获取地址的余额")
		fmt.Println("addr0余额", bal)
	}

	var (
		contractAddressHex string
		contractAddress    common.Address
	)
	{ // 部署多签合约测试 （a0/a1/a2）
		privkHex := a0.PrivkHex
		hexAddress := []string{a0.Address, a1.Address, a2.Address}
		mRequired := uint8(2)

		fmt.Println("owners:", hexAddress)
		chainID := big.NewInt(1)
		got, err := DeploySimpleMultiSigContract(*chainID, rpcClient, privkHex, hexAddress, mRequired)
		testtool.FailOnErr(t, err, "DeployMultiSigWalletContract()")
		fmt.Println("deployMultisigWalletContract got:", got)

		contractAddressHex = got
		contractAddress = common.HexToAddress(contractAddressHex)
		fmt.Println("contractAddressHex", contractAddressHex)
	}

	{ // 部署好后验证合约属性（owners/mRequired）
		owners, mRequired, err := GetContractInfo(rpcClient, contractAddressHex)
		testtool.FailOnErr(t, err, "Failed to get contract info")
		fmt.Println("contract info", owners, mRequired)
		testtool.FailOnFlag(t, len(owners) != 3, "len owners != 3", len(owners))
		testtool.FailOnFlag(t, mRequired != 2, "mRequired != 2", mRequired)
	}

	{ // 部署erc20 合约, owner 为 a0
		auth := bind.NewKeyedTransactor(a0.ToECDSAKey())
		addr, _, contrakt, err := contracts.DeployFixedSupplyToken(auth, rpcClient)
		testtool.FailOnErr(t, err, "Failed to deploy erc20 contract")
		erc20Contract = contrakt
		erc20ContractAddress = addr

		// 首先测试查询数据
		bal, e := erc20Contract.BalanceOf(&bind.CallOpts{}, a0.ToAddress())
		testtool.FailOnErr(t, e, "Failed to get balance of owner")
		fmt.Println("balance of erc20 owner", bal)
		shouldBe, flag := big.NewInt(1).SetString("1000000000000000000000000", 10)
		testtool.FailOnFlag(t, !flag, "set int failed")
		testtool.FailOnFlag(t, bal.Cmp(shouldBe) != 0, "余额异常")
	}

	{ //合约部署后往多签合约转入ERC20资金
		funds := big.NewInt(testtool.E18 * 3)
		_, err := erc20Contract.Transfer(bind.NewKeyedTransactor(a0.ToECDSAKey()), contractAddress, funds)
		testtool.FailOnErr(t, err, "Erc20 转账失败")

		bal, err := erc20Contract.BalanceOf(&bind.CallOpts{}, contractAddress)
		testtool.FailOnErr(t, err, "无法获取合约地址erc20余额")
		fmt.Println("合约地址erc20余额", bal)
		testtool.FailOnFlag(t, bal.Cmp(funds) != 0, "合约地址上的erc20余额不符合预期", bal)
	}

	outAddr := a3.Address
	transferValue := big.NewInt(testtool.E18 * 2)
	{ // 交易测试,a0+a2签名，从合约内转账erc20资金到a3 上
		multisigContract, err := contracts.NewSimpleMultiSig(contractAddress, rpcClient)
		testtool.FailOnErr(t, err, "构建多签合约调用时异常,检查合约地址和rpc server")
		nonce, err := multisigContract.Nonce(&bind.CallOpts{Pending: true})
		testtool.FailOnErr(t, err, "无法获取合约内部nonce")
		var (
			sigV                                 []uint8    //签名
			sigR, sigS                           [][32]byte //签名
			privkHex                             string
			multisigContractAddress, fromAddress string //多签合约地址，发起地址
			destination, executor                string //toAddress
			value, gasLimit                      *big.Int
			data                                 []byte
		)
		// 012由0发起，0和2签名, 把钱赚到1的地址上,executor 为0
		privkHex = a0.PrivkHex
		multisigContractAddress = contractAddressHex
		fromAddress = a0.Address
		executor = a0.Address
		destination = erc20ContractAddress.Hex() //TODO 这里在填错的情况下也可以执行成功，很容易弄错，可以考虑校验或者其他措施（eg:验证地址为erc20地址，验证与其他相关地址fromAddr/toAddr/multisigContractAddr/executorAddr不一样）
		value = big.NewInt(0)
		gasLimit = big.NewInt(239963)
		data, err = contracts.NewERC20InterfaceABIHelper().PackedTransfer(common.HexToAddress(outAddr), transferValue)
		testtool.FailOnErr(t, err, "打包erc20 transfer data失败")

		for _, add := range []*testtool.AddrInfo{a0, a2} {
			v, r, s, err := SimpleMultiSigExecuteSign(chainID, add.PrivkHex, multisigContractAddress, destination, executor, uint64(nonce.Int64()), value, gasLimit, data)
			testtool.FailOnErr(t, err, "create sig failed")
			sigV = append(sigV, v)
			sigR = append(sigR, r)
			sigS = append(sigS, s)
		}

		txid, err := ExecuteTX(&TxParams{
			backend:                 rpcClient,
			sigV:                    sigV,
			sigR:                    sigR,
			sigS:                    sigS,
			privkHex:                privkHex,
			multisigContractAddress: multisigContractAddress,
			fromAddress:             fromAddress,
			destination:             destination,
			executor:                executor,
			value:                   value,
			gasLimit:                gasLimit,
			data:                    data,
		})
		testtool.FailOnErr(t, err, "Execute Failed")
		fmt.Println("execute txid", txid)

	}

	{ // 完了检查确实转账成功
		bal, err := erc20Contract.BalanceOf(&bind.CallOpts{}, common.HexToAddress(outAddr))
		testtool.FailOnErr(t, err, "FonGetBal")
		fmt.Println("balance of new tx (erc20 transfer)", bal)
		testtool.FailOnFlag(t, bal.Cmp(transferValue) != 0, "erc20 多签转账失败，余额不符合预期")
	}

	{ // 查询erc20合约的日志，确认成功
		count := 0
		filter, err := erc20Contract.FilterTransfer(&bind.FilterOpts{}, []common.Address{contractAddress}, []common.Address{common.HexToAddress(outAddr)})
		testtool.FailOnErr(t, err, "无法查询erc20 合约时间日志")
		defer filter.Close()
		for {
			if !filter.Next() {
				break
			}
			evt := filter.Event
			fmt.Println("erv20 evt tokens:", evt.Tokens)
			count++
		}
		testtool.FailOnFlag(t, count != 1, "明确应该包含 1 条erc20转账记录")
	}
	{ // 调试合约日志
		// multisigContract, err := contracts.NewSimpleMultiSig(contractAddress, rpcClient)
		// testtool.FailOnErr(t, err, "构建多签合约调用时异常,检查合约地址和rpc server")
		// go func() {
		// 	ito, err := multisigContract.FilterDebugRecover(&bind.FilterOpts{Start: 0})
		// 	testtool.FailOnErr(t, err, "过滤合约日志异常")
		// 	for {
		// 		if !ito.Next() {
		// 			break
		// 		}
		// 		evt := ito.Event
		// 		log.Println("evt recoverd address:", evt.Recovered.Hex())
		// 	}
		// }()

		// go func() {
		// 	ito, err := multisigContract.FilterDbgExecuteParam(&bind.FilterOpts{Start: 0})
		// 	testtool.FailOnErr(t, err, "过滤合约日志异常")
		// 	for {
		// 		if !ito.Next() {
		// 			break
		// 		}
		// 		evt := ito.Event
		// 		log.Println("evt seperator:", hex.EncodeToString(evt.Sperator[:]))
		// 		log.Println("evt TxInputHash:", hex.EncodeToString(evt.TxInputHash[:]))
		// 		log.Println("evt TotalHash:", hex.EncodeToString(evt.TotalHash[:]))
		// 		log.Println("evt txInput:", hex.EncodeToString(evt.TxInput[:]))
		// 	}
		// }()

		// go func() {
		// 	ito, err := multisigContract.FilterExecute(&bind.FilterOpts{Start: 0})
		// 	testtool.FailOnErr(t, err, "过滤合约日志异常")
		// 	for {
		// 		if !ito.Next() {
		// 			break
		// 		}
		// 		evt := ito.Event
		// 		log.Println("evt confirmAddrs:", evt.ConfirmAddrs)
		// 		log.Println("evt dest:", evt.Destination.Hex())
		// 	}
		// }()

		// go func() {
		// 	ito, err := multisigContract.FilterDeposit(&bind.FilterOpts{Start: 0}, nil)
		// 	testtool.FailOnErr(t, err, "过滤合约日志异常")
		// 	for {
		// 		if !ito.Next() {
		// 			break
		// 		}
		// 		evt := ito.Event
		// 		log.Println("evt From:", evt.From.Hex())
		// 		log.Println("evt value:", evt.Value)
		// 	}
		// }()

	}
}

// 本测试要求本地7545端口运行有ganache,并且至少有一个账号余额超过5eth
func TestSimplemultisigGanache(t *testing.T) {
	// t.SkipNow()

	const (
		rpcHost = "http://localhost:7545"
	)
	var (
		a0, a1, a2, a3 *testtool.AddrInfo
		addrs          []*testtool.AddrInfo
		client         *ethclient.Client
		err            error
	)
	{ // init vars
		// 生成4个地址，并排序
		for i := 0; i < 4; i++ {
			addr := testtool.GenAddr()
			addrs = append(addrs, addr)
		}
		sort.Slice(addrs, func(i, j int) bool {
			return addrs[i].Address < addrs[j].Address
		})
		a0, a1, a2, a3 = addrs[0], addrs[1], addrs[2], addrs[3]

		client, err = ethclient.Dial(rpcHost)
		testtool.FailOnErr(t, err, "dial failed")
	}

	{ //ganache only
		testtool.PrepareFunds4address(t, rpcHost, a0.Address, 3)
	}
	{ // 首先确定addr0里的余额
		bal, err := client.BalanceAt(context.Background(), a0.ToAddress(), nil)
		testtool.FailOnErr(t, err, "无法获取地址的余额")
		fmt.Println("addr0余额", bal)
	}

	var (
		contractAddressHex string
		contractAddress    common.Address
	)
	{ // 部署合约测试
		type args struct {
			privkHex   string
			hexAddress []string
			mRequired  uint8
		}
		arg := args{
			privkHex:   a0.PrivkHex,
			hexAddress: []string{a0.Address, a1.Address, a2.Address},
			mRequired:  2,
		}

		fmt.Println("owners:", arg.hexAddress)
		chainID := big.NewInt(1)
		got, err := DeploySimpleMultiSigContract(*chainID, client, arg.privkHex, arg.hexAddress, arg.mRequired)
		if err != nil {
			t.Errorf("DeployMultiSigWalletContract() error = %v", err)
			t.FailNow()
		}
		fmt.Println("deployMultisigWalletContract got:", got)

		contractAddressHex = got
		contractAddress = common.HexToAddress(contractAddressHex)
		fmt.Println("contractAddressHex", contractAddressHex)
	}

	{ // 部署好后验证合约属性（owners/mRequired）
		owners, mRequired, err := GetContractInfo(client, contractAddressHex)
		testtool.FailOnErr(t, err, "Failed to get contract info")
		fmt.Println("contract info", owners, mRequired)
		testtool.FailOnFlag(t, len(owners) != 3, "len owners != 3", len(owners))
		testtool.FailOnFlag(t, mRequired != 2, "mRequired != 2", mRequired)
	}

	{ //合约部署后往其中转入资金
		value := big.NewInt(1).Mul(big.NewInt(testtool.E18), big.NewInt(2))
		ctx := context.Background()
		nonce, err := client.NonceAt(ctx, contractAddress, nil)
		testtool.FailOnErr(t, err, "Failed to get nonce")
		tx := types.NewTransaction(nonce, contractAddress, value, uint64(6721975), big.NewInt(20000000000), nil)
		signer := types.MakeSigner(params.TestChainConfig, nil)
		signature, err := crypto.Sign(signer.Hash(tx).Bytes(), a0.ToECDSAKey())
		testtool.FailOnErr(t, err, "签名交易失败")
		tx, err = tx.WithSignature(signer, signature)
		testtool.FailOnErr(t, err, "为交易附加签名数据错误")
		err = client.SendTransaction(ctx, tx)
		testtool.FailOnErr(t, err, "充值到合约地址异常")

		bal, err := client.BalanceAt(ctx, contractAddress, nil)
		testtool.FailOnErr(t, err, "无法获取合约地址的余额")
		testtool.FailOnFlag(t, bal.Cmp(value) != 0, fmt.Sprintf("合约地址的余额异常，应该是 %v, 实际上：%s", value, bal.String()))
		fmt.Println("合约地址当前余额", bal)

		bal, err = client.BalanceAt(ctx, a0.ToAddress(), nil)
		testtool.FailOnErr(t, err, "无法获取地址的余额")
		fmt.Println("地址0当前余额", bal)
	}

	outAddr := a3.Address
	transferValue := big.NewInt(testtool.E18)
	{ // 交易测试
		multisigContract, err := contracts.NewSimpleMultiSig(contractAddress, client)
		testtool.FailOnErr(t, err, "构建多签合约调用时异常,检查合约地址和rpc server")
		nonce, err := multisigContract.Nonce(&bind.CallOpts{Pending: true})
		testtool.FailOnErr(t, err, "无法获取合约内部nonce")
		var (
			sigV                                 []uint8    //签名
			sigR, sigS                           [][32]byte //签名
			privkHex                             string
			multisigContractAddress, fromAddress string //多签合约地址，发起地址
			destination, executor                string //toAddress
			value, gasLimit                      *big.Int
			data                                 []byte
		)
		// 012由0发起，0和2签名, 把钱赚到1的地址上,executor 为0
		privkHex = a0.PrivkHex
		multisigContractAddress = contractAddressHex
		fromAddress = a0.Address
		executor = a0.Address
		destination = outAddr
		value = transferValue
		gasLimit = big.NewInt(239963)
		data = []byte("")

		for _, add := range []*testtool.AddrInfo{a0, a2} {
			v, r, s, err := SimpleMultiSigExecuteSign(chainID, add.PrivkHex, multisigContractAddress, destination, executor, uint64(nonce.Int64()), value, gasLimit, data)
			testtool.FailOnErr(t, err, "create sig failed")
			sigV = append(sigV, v)
			sigR = append(sigR, r)
			sigS = append(sigS, s)
		}

		txid, err := ExecuteTX(&TxParams{
			backend:                 client,
			sigV:                    sigV,
			sigR:                    sigR,
			sigS:                    sigS,
			privkHex:                privkHex,
			multisigContractAddress: multisigContractAddress,
			fromAddress:             fromAddress,
			destination:             destination,
			executor:                executor,
			value:                   value,
			gasLimit:                gasLimit,
			data:                    data,
		})
		testtool.FailOnErr(t, err, "Execute Failed")
		fmt.Println("execute txid", txid)

	}

	{ // 完了检查确实转账成功
		bal, err := client.BalanceAt(context.Background(), a3.ToAddress(), nil)
		testtool.FailOnErr(t, err, "FonGetBal")
		fmt.Println("balance of new tx", bal)
		testtool.FailOnFlag(t, bal.Cmp(transferValue) != 0, "余额不对", transferValue)
	}
	{ // 调试合约日志
		multisigContract, err := contracts.NewSimpleMultiSig(contractAddress, client)
		testtool.FailOnErr(t, err, "构建多签合约调用时异常,检查合约地址和rpc server")
		// go func() {
		// 	ito, err := multisigContract.FilterDebugRecover(&bind.FilterOpts{Start: 0})
		// 	testtool.FailOnErr(t, err, "过滤合约日志异常")
		// 	for {
		// 		if !ito.Next() {
		// 			break
		// 		}
		// 		evt := ito.Event
		// 		log.Println("evt recoverd address:", evt.Recovered.Hex())
		// 	}
		// }()

		// go func() {
		// 	ito, err := multisigContract.FilterDbgExecuteParam(&bind.FilterOpts{Start: 0})
		// 	testtool.FailOnErr(t, err, "过滤合约日志异常")
		// 	for {
		// 		if !ito.Next() {
		// 			break
		// 		}
		// 		evt := ito.Event
		// 		log.Println("evt seperator:", hex.EncodeToString(evt.Sperator[:]))
		// 		log.Println("evt TxInputHash:", hex.EncodeToString(evt.TxInputHash[:]))
		// 		log.Println("evt TotalHash:", hex.EncodeToString(evt.TotalHash[:]))
		// 		log.Println("evt txInput:", hex.EncodeToString(evt.TxInput[:]))
		// 	}
		// }()
		executeFilter, err := multisigContract.FilterExecute(&bind.FilterOpts{Start: 0})
		testtool.FailOnErr(t, err, "过滤合约日志异常")
		for {
			if !executeFilter.Next() {
				break
			}
			evt := executeFilter.Event
			log.Println("evt execute, confirmAddrs:", evt.ConfirmAddrs)
			log.Println("evt execute, dest:", evt.Destination.Hex())
		}

		depositLogFilter, err := multisigContract.FilterDeposit(&bind.FilterOpts{Start: 0}, nil)
		testtool.FailOnErr(t, err, "过滤合约日志异常")
		for {
			if !depositLogFilter.Next() {
				break
			}
			evt := depositLogFilter.Event
			log.Println("evt deposit,From:", evt.From.Hex())
			log.Println("evt deposit,value:", evt.Value)
		}
	}
}

// 本测试要求本地7545端口运行有ganache,并且至少有一个账号余额超过5eth
// 仅使用多签合约abiHelper进行开发
// TODO 在execute时同时payable会使交易失败，多签交易时确认执行人设置的value为0
// (多签本质上时合约交易调用，避免调用时又pay)
func TestSimplemultisigAbiHelper(t *testing.T) {
	const (
		rpcHost = "http://localhost:7545"
	)
	var (
		a0, a1, a2, a3 *testtool.AddrInfo
		addrs          []*testtool.AddrInfo
		client         *ethclient.Client
		abiHelper      *contracts.SimpleMultiSigABIHelper
		err            error
	)
	{ // init vars
		// 生成4个地址，并排序
		for i := 0; i < 4; i++ {
			addr := testtool.GenAddr()
			addrs = append(addrs, addr)
		}
		sort.Slice(addrs, func(i, j int) bool {
			return addrs[i].Address < addrs[j].Address
		})
		a0, a1, a2, a3 = addrs[0], addrs[1], addrs[2], addrs[3]

		client, err = ethclient.Dial(rpcHost)
		testtool.FailOnErr(t, err, "dial failed")
		_ = a3

		abiHelper = contracts.NewSimpleMultiSigABIHelper()
	}

	{ //ganache only
		testtool.PrepareFunds4address(t, rpcHost, a0.Address, 5)
	}
	{ // 首先确定addr0里的余额
		bal, err := client.BalanceAt(context.Background(), a0.ToAddress(), nil)
		testtool.FailOnErr(t, err, "无法获取地址的余额")
		fmt.Println("addr0余额", bal)
		testtool.FailOnFlag(t, bal.Cmp(big.NewInt(testtool.E18*5)) != 0, "Wrong balance")
	}

	var (
		contractAddressHex string
		contractAddress    common.Address
	)
	{ // 部署合约测试
		owners := []common.Address{a0.ToAddress(), a1.ToAddress(), a2.ToAddress()}
		mRequired := 2
		chainID := big.NewInt(1)

		fmt.Println("owners:", owners)
		createMultisigData, err := contracts.PackedDeploySimpleMultiSig(big.NewInt(int64(mRequired)), owners, chainID)
		testtool.FailOnErr(t, err, "Failed to pack create simplemultisig contract data")
		a0Nonce, err := client.PendingNonceAt(context.Background(), a0.ToAddress())
		testtool.FailOnErr(t, err, "Failed to get a0 nonce")

		tx := types.NewContractCreation(a0Nonce, big.NewInt(0), uint64(6721975), big.NewInt(20000000000), createMultisigData)
		signer := types.MakeSigner(params.TestChainConfig, nil)
		signature, err := crypto.Sign(signer.Hash(tx).Bytes(), a0.ToECDSAKey())
		testtool.FailOnErr(t, err, "签名交易失败")

		tx, err = tx.WithSignature(signer, signature)
		testtool.FailOnErr(t, err, "tx.WithSignature error")
		err = client.SendTransaction(context.Background(), tx)
		testtool.FailOnErr(t, err, "Failed to create simpleMultisigContract")

		rec, err := client.TransactionReceipt(context.Background(), tx.Hash())
		testtool.FailOnErr(t, err, "无法获取多签地址")
		contractAddress = rec.ContractAddress
		contractAddressHex = contractAddress.Hex()
		testtool.FailOnFlag(t, strings.Index(contractAddressHex, "00000000") != -1, "Create contract failed, address is zero")

	}

	{ // 部署好后验证合约属性（owners/mRequired）
		//get owners length
		packedGetOwnersLengthData, err := abiHelper.PackedGetOwersLength()
		testtool.FailOnErr(t, err, "PackedGetOwnersLen")
		retBytes, err := client.CallContract(context.Background(), ethereum.CallMsg{
			From: a0.ToAddress(),
			To:   &contractAddress,
			Data: packedGetOwnersLengthData,
		}, nil)
		testtool.FailOnErr(t, err, "Failed to call contract")
		ownerLen, err := abiHelper.UnpackGetOwersLength(retBytes)
		testtool.FailOnErr(t, err, "Failed to unpack")
		fmt.Println("get ownerlen: ", ownerLen)

		// get first owner
		getFirstAddrData, err := abiHelper.PackedOwnersArr(big.NewInt(0))
		testtool.FailOnErr(t, err, "FailedToPack")
		firstAddrRet, err := client.CallContract(context.Background(), ethereum.CallMsg{
			From: a0.ToAddress(),
			To:   &contractAddress,
			Data: getFirstAddrData,
		}, nil)
		testtool.FailOnErr(t, err, "Failed to call contract")
		firstAddr, err := abiHelper.UnpackOwnersArr(firstAddrRet)
		testtool.FailOnErr(t, err, "Failed to unpack")
		fmt.Println("get first owner: ", firstAddr.Hex())
	}

	{ //合约部署后往其中转入资金(2 eth)
		value := big.NewInt(testtool.E18 * 2)
		ctx := context.Background()
		nonce, err := client.NonceAt(ctx, a0.ToAddress(), nil)
		testtool.FailOnErr(t, err, "Failed to get nonce")
		tx := types.NewTransaction(nonce, contractAddress, value, uint64(6721975), big.NewInt(20000000000), nil)
		signer := types.MakeSigner(params.TestChainConfig, nil)
		signature, err := crypto.Sign(signer.Hash(tx).Bytes(), a0.ToECDSAKey())
		testtool.FailOnErr(t, err, "签名交易失败")
		tx, err = tx.WithSignature(signer, signature)
		testtool.FailOnErr(t, err, "为交易附加签名数据错误")
		err = client.SendTransaction(ctx, tx)
		testtool.FailOnErr(t, err, "充值到合约地址异常")

		bal, err := client.BalanceAt(ctx, contractAddress, nil)
		testtool.FailOnErr(t, err, "无法获取合约地址的余额")
		testtool.FailOnFlag(t, bal.Cmp(value) != 0, fmt.Sprintf("合约地址的余额异常，应该是 %v, 实际上：%s", value, bal.String()))
		fmt.Println("合约地址当前余额", bal)

		bal, err = client.BalanceAt(ctx, a0.ToAddress(), nil)
		testtool.FailOnErr(t, err, "无法获取地址的余额")
		fmt.Println("地址0当前余额", bal)
	}

	outAddr := a3.Address
	transferValue := big.NewInt(testtool.E18)
	{ // 交易测试
		callNonceData, _ := abiHelper.PackedNonce()
		callNonceBytes, err := client.CallContract(context.Background(), ethereum.CallMsg{
			From: a0.ToAddress(),
			To:   &contractAddress,
			Data: callNonceData,
		}, nil)
		testtool.FailOnErr(t, err, "Failed to call contract")
		nonce, err := abiHelper.UnpackNonce(callNonceBytes)
		testtool.FailOnErr(t, err, "Failed to unpack nonce")

		var (
			sigV                    []uint8    //签名
			sigR, sigS              [][32]byte //签名
			multisigContractAddress string     //多签合约地址，发起地址
			destination, executor   string     //toAddress
			value, gasLimit         *big.Int
			data                    []byte
		)
		// 012由0发起，0和2签名, 把钱赚到1的地址上,executor 为0
		multisigContractAddress = contractAddressHex
		executor = a0.Address
		destination = outAddr
		value = transferValue
		gasLimit = big.NewInt(239963)
		data = []byte("")

		for _, add := range []*testtool.AddrInfo{a0, a2} {
			v, r, s, err := SimpleMultiSigExecuteSign(chainID, add.PrivkHex, multisigContractAddress, destination, executor, uint64(nonce.Int64()), value, gasLimit, data)
			testtool.FailOnErr(t, err, "create sig failed")
			sigV = append(sigV, v)
			sigR = append(sigR, r)
			sigS = append(sigS, s)
		}

		packedTxData, err := abiHelper.PackedExecute(sigV, sigR, sigS, common.HexToAddress(destination), value, data, common.HexToAddress(executor), gasLimit)
		testtool.FailOnErr(t, err, "Pack multisig execute faied")

		ctx := context.Background()
		a0Nonce, err := client.NonceAt(ctx, a0.ToAddress(), nil)
		testtool.FailOnErr(t, err, "Failed to get a0Nonce")
		// tx := types.NewTransaction(a0Nonce, contractAddress, transferValue, uint64(6721975), big.NewInt(20000000000), packedTxData)
		tx := types.NewTransaction(a0Nonce, contractAddress, big.NewInt(0), uint64(6721975), big.NewInt(20000000000), packedTxData)
		// tx := types.NewTransaction(a0Nonce, contractAddress, big.NewInt(testtool.E18), uint64(6721975), big.NewInt(20000000000), packedTxData)
		signer := types.MakeSigner(params.TestChainConfig, nil)
		signature, err := crypto.Sign(signer.Hash(tx).Bytes(), a0.ToECDSAKey())
		testtool.FailOnErr(t, err, "签名交易失败")
		tx, err = tx.WithSignature(signer, signature)
		testtool.FailOnErr(t, err, "为交易附加签名数据错误")
		err = client.SendTransaction(ctx, tx)
		testtool.FailOnErr(t, err, "调用rpc sendTransaction 错误")
	}

	{ // 完了检查确实转账成功
		bal, err := client.BalanceAt(context.Background(), a3.ToAddress(), nil)
		testtool.FailOnErr(t, err, "FonGetBal")
		fmt.Println("balance of new tx", bal)
		testtool.FailOnFlag(t, bal.Cmp(transferValue) != 0, "余额不对", transferValue)
	}
	{ // 调试合约日志
		// multisigContract, err := contracts.NewSimpleMultiSig(contractAddress, client)
		// testtool.FailOnErr(t, err, "构建多签合约调用时异常,检查合约地址和rpc server")
		// go func() {
		// 	ito, err := multisigContract.FilterDebugRecover(&bind.FilterOpts{Start: 0})
		// 	testtool.FailOnErr(t, err, "过滤合约日志异常")
		// 	for {
		// 		if !ito.Next() {
		// 			break
		// 		}
		// 		evt := ito.Event
		// 		log.Println("evt recoverd address:", evt.Recovered.Hex())
		// 	}
		// }()

		// go func() {
		// 	ito, err := multisigContract.FilterDbgExecuteParam(&bind.FilterOpts{Start: 0})
		// 	testtool.FailOnErr(t, err, "过滤合约日志异常")
		// 	for {
		// 		if !ito.Next() {
		// 			break
		// 		}
		// 		evt := ito.Event
		// 		log.Println("evt seperator:", hex.EncodeToString(evt.Sperator[:]))
		// 		log.Println("evt TxInputHash:", hex.EncodeToString(evt.TxInputHash[:]))
		// 		log.Println("evt TotalHash:", hex.EncodeToString(evt.TotalHash[:]))
		// 		log.Println("evt txInput:", hex.EncodeToString(evt.TxInput[:]))
		// 	}
		// }()
		// executeFilter, err := multisigContract.FilterExecute(&bind.FilterOpts{Start: 0})
		// testtool.FailOnErr(t, err, "过滤合约日志异常")
		// for {
		// 	if !executeFilter.Next() {
		// 		break
		// 	}
		// 	evt := executeFilter.Event
		// 	log.Println("evt confirmAddrs:", evt.ConfirmAddrs)
		// 	log.Println("evt dest:", evt.Destination.Hex())
		// }

		// depositLogFilter, err := multisigContract.FilterDeposit(&bind.FilterOpts{Start: 0}, nil)
		// testtool.FailOnErr(t, err, "过滤合约日志异常")
		// for {
		// 	if !depositLogFilter.Next() {
		// 		break
		// 	}
		// 	evt := depositLogFilter.Event
		// 	log.Println("evt From:", evt.From.Hex())
		// 	log.Println("evt value:", evt.Value)
		// }
	}
}
