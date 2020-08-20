package eth

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/dabankio/devtools4chains"
	"github.com/dabankio/wallet-core/core/eth"
	"github.com/dabankio/wallet-core/core/eth/internalized/testtool"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
)

// 示例代码，展示如何使用多签sdk 进行以太坊开发
// 本测试要求安装又npm i -g ganache-cli
// 仅使用多签合约abiHelper进行开发
// TODO 在execute时同时payable会使交易失败，多签交易时确认执行人设置的value为0
// (多签本质上时合约交易调用，避免调用时又pay)
func TestSimplemultisigAbiHelper(t *testing.T) {
	rq := require.New(t)

	killFunc, port, err := devtools4chains.DockerRunGanacheCli(&devtools4chains.DockerRunOptions{
		AutoRemove: true,
	})
	t.Cleanup(killFunc)
	var rpcHost = fmt.Sprintf("http://localhost:%d", port)
	const (
		bucketNum = 10
		bucketIdx = 0
		chainID   = 2
	)
	var (
		ctx      = context.Background()
		gasLimit = int64(6721975) //这里就写死了，便于测试

		a0, a1, a2, a3  *AddrInfo
		addrs           []*AddrInfo
		client          *ethclient.Client
		abiHelper       *eth.SimpleMultiSigABIHelper
		suggestGasPrice *big.Int
		expireTime      = time.Now().Add(3 * 24 * time.Hour)
	)
	{ // init vars
		// 生成4个地址，并排序
		for i := 0; i < 4; i++ {
			addr := GenAddr()
			addrs = append(addrs, addr)
		}
		sort.Slice(addrs, func(i, j int) bool {
			return strings.ToLower(addrs[i].Address) < strings.ToLower(addrs[j].Address)
		})
		a0, a1, a2, a3 = addrs[0], addrs[1], addrs[2], addrs[3]

		client, err = ethclient.Dial(rpcHost)
		rq.Nil(err, "dial failed")

		testtool.WaitSomething(t, time.Minute, func() error { _, e := client.NetworkID(context.Background()); return e })

		abiHelper = eth.NewSimpleMultiSigABIHelper()

		suggestGasPrice, err = client.SuggestGasPrice(ctx)
		rq.Nil(err, "Failed to get gasPrice")

	}

	{ //ganache only
		PrepareFunds4address(t, rpcHost, a0.Address, 5)
	}
	{ // 首先确定addr0里的余额
		bal, err := client.BalanceAt(context.Background(), a0.ToAddress(), nil)
		rq.Nil(err, "无法获取地址的余额")
		fmt.Println("addr0余额", bal)
		rq.False(bal.Cmp(big.NewInt(E18*5)) != 0, "Wrong balance")
	}

	var (
		contractAddressHex string
		contractAddress    common.Address
	)
	{ // 部署合约测试
		mRequired := 2
		// fmt.Println("owners:", owners)
		ownersAddrWrap := eth.NewAddressesWrap()

		{
			for _, add := range []string{a0.Address, a1.Address, a2.Address} {
				eadd, err := eth.NewETHAddressFromHex(add)
				rq.Nil(err, "new addr from hex failed")
				ownersAddrWrap.AddOne(eadd)
			}
		}

		createMultisigData, err := eth.PackedDeploySimpleMultiSig(bucketNum, eth.NewBigInt(int64(mRequired)), ownersAddrWrap, eth.NewBigInt(chainID))
		rq.Nil(err, "Failed to pack create simplemultisig contract data")
		a0Nonce, err := client.PendingNonceAt(context.Background(), a0.ToAddress())
		rq.Nil(err, "Failed to get a0 nonce")

		ethtx := eth.NewETHTransactionForContractCreation(int64(a0Nonce), gasLimit, eth.NewBigInt(suggestGasPrice.Int64()), createMultisigData)

		encodedRlpTx, err := ethtx.EncodeRLP()
		rq.Nil(err, "failed to encode tx rlp")
		sig, err := eth.SignRawTransaction(encodedRlpTx, a0.PrivkHex)
		rq.Nil(err, "Failed to sign tx")

		{ // 一般步骤下，应该用签名好的数据调用广播api,创建合约，这里处理下数据然后调用geth jsonrpc 进行广播
			sigBytes, err := hexutil.Decode(sig)
			rq.Nil(err, "Failed to decode hexed sig ")
			var rawTx types.Transaction
			err = rlp.DecodeBytes(sigBytes, &rawTx)
			rq.Nil(err, "Failed to rlp decode tx")

			err = client.SendTransaction(context.Background(), &rawTx)
			rq.Nil(err, "Failed to send tx")

			// 获取部署好的合约地址
			rec, err := client.TransactionReceipt(context.Background(), rawTx.Hash())
			rq.Nil(err, "无法获取多签地址")

			contractAddress = rec.ContractAddress
			contractAddressHex = contractAddress.Hex()
			rq.False(strings.Index(contractAddressHex, "00000000") != -1, "Create contract failed, address is zero")
			fmt.Println("Deployed simplemultisig contract address", contractAddressHex)
		}
	}

	{ // 部署好后验证合约属性（owners/mRequired）
		//get owners length
		packedGetOwnersLengthData, err := abiHelper.PackedGetOwersLength()
		rq.Nil(err, "PackedGetOwnersLen")
		retBytes, err := client.CallContract(context.Background(), ethereum.CallMsg{
			From: a0.ToAddress(),
			To:   &contractAddress,
			Data: packedGetOwnersLengthData,
		}, nil)
		rq.Nil(err, "Failed to call contract")
		ownerLen, err := abiHelper.UnpackGetOwersLength(retBytes)
		rq.Nil(err, "Failed to unpack")
		fmt.Println("get ownerlen: ", ownerLen)
		rq.False(ownerLen != 3, "wrong owners length")

		// get first owner
		getFirstAddrData, err := abiHelper.PackedOwnersArr(eth.NewBigInt(0))
		rq.Nil(err, "FailedToPack")
		firstAddrRet, err := client.CallContract(context.Background(), ethereum.CallMsg{
			From: a0.ToAddress(),
			To:   &contractAddress,
			Data: getFirstAddrData,
		}, nil)
		rq.Nil(err, "Failed to call contract")
		firstAddr, err := abiHelper.UnpackOwnersArr(firstAddrRet)
		rq.Nil(err, "Failed to unpack")
		fmt.Println("get first owner: ", firstAddr.GetHex())
		rq.False(firstAddr.GetHex() != a0.Address, "获取到的合约持有人不符合预期")
	}

	{ //合约部署后往其中转入资金(2 eth)
		value := big.NewInt(E18 * 2)
		nonce, err := client.NonceAt(ctx, a0.ToAddress(), nil)
		rq.Nil(err, "Failed to get nonce")
		tx := types.NewTransaction(nonce, contractAddress, value, uint64(6721975), big.NewInt(20000000000), nil)
		signer := types.MakeSigner(params.TestChainConfig, nil)
		signature, err := crypto.Sign(signer.Hash(tx).Bytes(), a0.ToECDSAKey())
		rq.Nil(err, "签名交易失败")
		tx, err = tx.WithSignature(signer, signature)
		rq.Nil(err, "为交易附加签名数据错误")
		err = client.SendTransaction(ctx, tx)
		rq.Nil(err, "充值到合约地址异常")

		bal, err := client.BalanceAt(ctx, contractAddress, nil)
		rq.Nil(err, "无法获取合约地址的余额")
		rq.False(bal.Cmp(value) != 0, fmt.Sprintf("合约地址的余额异常，应该是 %v, 实际上：%s", value, bal.String()))
		fmt.Println("合约地址当前余额", bal)

		bal, err = client.BalanceAt(ctx, a0.ToAddress(), nil)
		rq.Nil(err, "无法获取地址的余额")
		fmt.Println("地址0当前余额", bal)
	}

	outAddr := a3.Address
	transferValue := eth.NewBigInt(E18)
	{ // 多签交易测试
		callNonceData, _ := abiHelper.PackedNonceBucket(eth.NewBigInt(bucketIdx))
		callNonceBytes, err := client.CallContract(context.Background(), ethereum.CallMsg{
			From: a0.ToAddress(),
			To:   &contractAddress,
			Data: callNonceData,
		}, nil)
		rq.Nil(err, "Failed to call contract")
		nonce, err := abiHelper.UnpackNonceBucket(callNonceBytes)
		rq.Nil(err, "Failed to unpack nonce")

		var (
			sigV                    *eth.Uint8ArrayWrap  //签名
			sigR, sigS              *eth.Byte32ArrayWrap //签名
			multisigContractAddress string               //多签合约地址，发起地址
			destination, executor   string               //toAddress
			value, gasLimit         *eth.BigInt
			data                    []byte
		)
		// 012由0发起，0和2签名, 把钱转出到3的地址上,executor 为0
		sigV = eth.NewUint8ArrayWrap()
		sigR, sigS = eth.NewByte32ArrayWrap(), eth.NewByte32ArrayWrap()
		multisigContractAddress = contractAddressHex
		// executor = a0.Address
		destination = outAddr
		value = transferValue
		gasLimit = eth.NewBigInt(239963)
		data = []byte("")

		for _, add := range []*AddrInfo{a0, a2} {
			//实际的使用场景中，应该把需要签名的数据分发给需要签名的人，分别签名，然后在合起来
			signRes, err := eth.UtilSimpleMultiSigExecuteSign(expireTime, chainID, add.PrivkHex, multisigContractAddress, destination, executor, nonce.GetInt64(), value, gasLimit, data)
			rq.Nil(err, "Failed to sign execute")

			sigV.AddOne(int8(signRes.V))
			sigR.AddOne(signRes.R.Get())
			sigS.AddOne(signRes.S.Get())
		}

		destAddr, err := eth.NewETHAddressFromHex(destination)
		rq.Nil(err, "Failed to new eth addr from hex")
		// executorAddr, err := NewETHAddressFromHex(executor)
		// rq.Nil(err, "Failed to new eth addr from hex")
		packedExecuteData, err := abiHelper.PackedExecute(bucketIdx, eth.NewBigInt(expireTime.Unix()), sigV, sigR, sigS, destAddr, value, data, &eth.ETHAddress{}, gasLimit)
		rq.Nil(err, "Pack multisig execute faied")

		a0Nonce, err := client.NonceAt(ctx, a0.ToAddress(), nil)
		rq.Nil(err, "Failed to get a0Nonce")
		// tx := types.NewTransaction(a0Nonce, contractAddress, transferValue, uint64(6721975), big.NewInt(20000000000), packedExecuteData)
		// tx := types.NewTransaction(a0Nonce, contractAddress, big.NewInt(0), uint64(6721975), big.NewInt(20000000000), packedExecuteData)
		contractETHAddr, err := eth.NewETHAddressFromHex(contractAddressHex)
		rq.Nil(err, "Failed to new eth addr from hex")

		ethtx := eth.NewETHTransaction(int64(a0Nonce), contractETHAddr, eth.NewBigInt(0), gasLimit.GetInt64(), eth.NewBigInt(suggestGasPrice.Int64()), packedExecuteData)

		encodedRlpTx, err := ethtx.EncodeRLP()
		rq.Nil(err, "failed to encode tx rlp")
		sig, err := eth.SignRawTransaction(encodedRlpTx, a0.PrivkHex)
		rq.Nil(err, "Failed to sign tx")

		{ // 一般步骤下，应该用签名好的数据调用广播api,创建合约，这里处理下数据然后调用 jsonrpc 进行广播
			sigBytes, err := hexutil.Decode(sig)
			rq.Nil(err, "Failed to decode hexed sig ")
			var rawTx types.Transaction
			err = rlp.DecodeBytes(sigBytes, &rawTx)
			rq.Nil(err, "Failed to rlp decode tx")

			err = client.SendTransaction(context.Background(), &rawTx)
			rq.Nil(err, "Failed to send tx")
		}
	}

	{ // 完了检查确实转账成功
		bal, err := client.BalanceAt(context.Background(), a3.ToAddress(), nil)
		rq.Nil(err, "FonGetBal")
		fmt.Println("多签转出账户余额", bal)
		rq.False(bal.Cmp(big.NewInt(transferValue.GetInt64())) != 0, "余额不对", transferValue)
	} //至此，简单多签的使用就没问题了

	// 接下来是多签在ERC20代币中的使用
	fmt.Println("-----------------------下面为简单多签合约在ERC20代币中的用法----------------------------")

	var (
		erc20Contract        *FixedSupplyToken
		erc20ContractAddress common.Address

		erc20AbiHelper = eth.NewERC20InterfaceABIHelper()
	)
	{ // 首先部署一个erc20代币,(在实际的使用场景中，代币是已经部署好的，不会有这个环节，直接使用地址即可)
		auth := bind.NewKeyedTransactor(a0.ToECDSAKey())
		addr, _, contrakt, err := DeployFixedSupplyToken(auth, client)
		rq.Nil(err, "Failed to deploy erc20 contract")
		erc20Contract = contrakt
		erc20ContractAddress = addr

		// 测试erc20查询数据
		bal, e := erc20Contract.BalanceOf(&bind.CallOpts{}, a0.ToAddress())
		rq.Nil(e, "Failed to get balance of owner")
		fmt.Println("balance of erc20 owner", bal)
		shouldBe, flag := big.NewInt(1).SetString("1000000000000000000000000", 10)
		rq.False(!flag, "set int failed")
		rq.False(bal.Cmp(shouldBe) != 0, "余额异常")
	}

	{ //erc20合约部署后往多签合约转入ERC20资金
		funds := big.NewInt(E18 * 3)
		_, err := erc20Contract.Transfer(bind.NewKeyedTransactor(a0.ToECDSAKey()), contractAddress, funds)
		rq.Nil(err, "Erc20 转账失败")

		bal, err := erc20Contract.BalanceOf(&bind.CallOpts{}, contractAddress)
		rq.Nil(err, "无法获取合约地址erc20余额")
		fmt.Println("合约地址erc20余额", bal)
		rq.False(bal.Cmp(funds) != 0, "合约地址上的erc20余额不符合预期", bal)
	}

	erc20OutAddr := a3.Address
	erc20TransferValue := eth.NewBigInt(E18 * 2)
	{ // 交易测试,a0+a2签名，从合约内转账erc20资金到a3 上
		//获取多签合约内部的nonce
		callNonceData, _ := abiHelper.PackedNonceBucket(eth.NewBigInt(bucketIdx))
		callNonceBytes, err := client.CallContract(context.Background(), ethereum.CallMsg{
			From: a0.ToAddress(),
			To:   &contractAddress,
			Data: callNonceData,
		}, nil)
		rq.Nil(err, "Failed to call contract")
		nonce, err := abiHelper.UnpackNonceBucket(callNonceBytes)
		rq.Nil(err, "Failed to unpack nonce")

		var (
			sigV                    *eth.Uint8ArrayWrap  //签名
			sigR, sigS              *eth.Byte32ArrayWrap //签名
			multisigContractAddress string               //多签合约地址，发起地址
			destination, executor   string               //toAddress
			value, gasLimit         *eth.BigInt
			data                    []byte
		)
		// 012由0发起，0和2签名, 把钱赚到3的地址上,executor 为0
		sigV = eth.NewUint8ArrayWrap()
		sigR, sigS = eth.NewByte32ArrayWrap(), eth.NewByte32ArrayWrap()
		multisigContractAddress = contractAddressHex
		executor = a0.Address

		// 区别于主币转账，erc20转账关键在于此处，特别注意data的创建
		value = eth.NewBigInt(0) //本质上为合约调用，所以不需要value
		gasLimit = eth.NewBigInt(239963)
		destination = erc20ContractAddress.Hex()
		erc20OutEthAddr, err := eth.NewETHAddressFromHex(erc20OutAddr)
		rq.Nil(err, "Failed to new eth addr from hex")
		data, err = erc20AbiHelper.PackedTransfer(erc20OutEthAddr, erc20TransferValue)

		for _, add := range []*AddrInfo{a0, a2} {
			//实际的使用场景中，应该把需要签名的数据分发给需要签名的人，分别签名，然后在合起来
			signRes, err := eth.UtilSimpleMultiSigExecuteSign(expireTime, chainID, add.PrivkHex, multisigContractAddress, destination, executor, nonce.GetInt64(), value, gasLimit, data)
			rq.Nil(err, "Failed to sign execute")

			sigV.AddOne(int8(signRes.V))
			sigR.AddOne(signRes.R.Get())
			sigS.AddOne(signRes.S.Get())
		}

		destAddr, err := eth.NewETHAddressFromHex(destination)
		rq.Nil(err, "Failed to new eth addr from hex")
		executorAddr, err := eth.NewETHAddressFromHex(executor)
		rq.Nil(err, "Failed to new eth addr from hex")
		packedExecuteData, err := abiHelper.PackedExecute(bucketIdx, eth.NewBigInt(expireTime.Unix()), sigV, sigR, sigS, destAddr, value, data, executorAddr, gasLimit)
		rq.Nil(err, "Pack multisig execute faied")

		a0Nonce, err := client.NonceAt(ctx, a0.ToAddress(), nil)
		rq.Nil(err, "Failed to get a0Nonce")
		contractETHAddr, err := eth.NewETHAddressFromHex(contractAddressHex)
		rq.Nil(err, "Failed to new eth addr from hex")

		ethtx := eth.NewETHTransaction(int64(a0Nonce), contractETHAddr, eth.NewBigInt(0), gasLimit.GetInt64(), eth.NewBigInt(suggestGasPrice.Int64()), packedExecuteData)

		encodedRlpTx, err := ethtx.EncodeRLP()
		rq.Nil(err, "failed to encode tx rlp")
		sig, err := eth.SignRawTransaction(encodedRlpTx, a0.PrivkHex)
		rq.Nil(err, "Failed to sign tx")

		{ // 一般步骤下，应该用签名好的数据调用广播api,创建合约，这里处理下数据然后调用 jsonrpc 进行广播
			sigBytes, err := hexutil.Decode(sig)
			rq.Nil(err, "Failed to decode hexed sig ")
			var rawTx types.Transaction
			err = rlp.DecodeBytes(sigBytes, &rawTx)
			rq.Nil(err, "Failed to rlp decode tx")

			err = client.SendTransaction(context.Background(), &rawTx)
			rq.Nil(err, "Failed to send tx")
		}
	}

	{ // 完了检查确实转账成功
		bal, err := erc20Contract.BalanceOf(&bind.CallOpts{}, common.HexToAddress(erc20OutAddr))
		rq.Nil(err, "FonGetBal")
		fmt.Println("balance of new tx (erc20 transfer)", bal)
		rq.False(bal.Cmp(big.NewInt(erc20TransferValue.GetInt64())) != 0, "erc20 多签转账失败，余额不符合预期")
	}

	{ // 调试合约日志
		// multisigContract, err := contracts.NewSimpleMultiSig(contractAddress, client)
		// rq.Nil(err, "构建多签合约调用时异常,检查合约地址和rpc server")
		// go func() {
		// 	ito, err := multisigContract.FilterDebugRecover(&bind.FilterOpts{Start: 0})
		// 	rq.Nil(err, "过滤合约日志异常")
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
		// 	rq.Nil(err, "过滤合约日志异常")
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
		// rq.Nil(err, "过滤合约日志异常")
		// for {
		// 	if !executeFilter.Next() {
		// 		break
		// 	}
		// 	evt := executeFilter.Event
		// 	log.Println("evt confirmAddrs:", evt.ConfirmAddrs)
		// 	log.Println("evt dest:", evt.Destination.Hex())
		// }

		// depositLogFilter, err := multisigContract.FilterDeposit(&bind.FilterOpts{Start: 0}, nil)
		// rq.Nil(err, "过滤合约日志异常")
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
