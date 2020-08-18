package eth

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/dabankio/wallet-core/core/eth/internalized/testtool"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

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

		expireTime := time.Now().Add(time.Hour)
		for _, add := range []*testtool.AddrInfo{a0, a2} {
			v, r, s, err := SimpleMultiSigExecuteSign(expireTime, chainID, add.PrivkHex, multisigContractAddress, destination, executor, nonce, value, gasLimit, data)
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
