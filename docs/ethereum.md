# 以太坊相关

API使用首先建议阅读 [api文档说明](./api.md)

以太坊的实现都在`core/eth`目录下面

使用之前需要对以太坊有基本的认识，了解基本的交易过程，下面列出了一些建议的学习资料：
- 官方仓库 https://github.com/ethereum/go-ethereum
- 仓库wiki https://github.com/ethereum/go-ethereum/wiki
- 精通以太坊 https://github.com/ethereumbook/ethereumbook
- 以太坊智能合约solidity语言 https://solidity.readthedocs.io/


- 一些开发工具 ganache,
- 区块浏览 https://etherscan.io
- 第三方api https://infura.io/docs
- 了解gas gasPrice gasLimit 的概念和如何获取

## 基本交易过程

- 生成地址
- 转入资金
- 查询交易计数、选择gasPrice gasLimit等
- 构造原始交易
- 使用私钥对交易进行签名
- 广播交易
- 等待打包、确认

## 多重签名一般使用过程

目前，我们选择了一个简单的多重签名合约，参见 ` core/eth/internal/contracts/SimpleMultiSig.sol`,修改自 https://github.com/christianlundkvist/simple-multisig ,该合约一次部署多次使用，一次交易一次上链，不支持变更成员，最多支持128个签名成员，支持对其他合约的调用（比如ERC20）

如果偏好其他多签的实现可以联系我们 support@dabank.io ,末尾也列出了其他几个我们调研过的多签合约

典型的多签使用流程为：
- 收集成员地址
- 部署多签合约
- 往多签合约转入资金
- 从多签合约发起交易，构造交易信息
- 成员分别对交易信息签名
- 收集成员签名，用于构造原始交易
- 广播者对原始交易进行签名
- 广播交易，等待打包

ERC20的多签会更为复杂一点,在集成测试代码里也有示例

## 示例代码

我们建议参考集成测试的代码，了解使用方法，此处的代码仅作展示用，方便了解特性和api设计。可能存在滞后性。

构建和签名交易
```golang
//golang 代码
// :构造交易
toA1Address, err := eth.NewETHAddressFromHex(addrs[1].address)
tx := eth.NewETHTransaction(int64(a0Nonce), toA1Address, eth.NewBigInt(int64(transferAmount)), int64(gasLimit), eth.NewBigInt(int64(gasPrice)), nil)

rawTx, err := tx.EncodeRLP()
// :签名交易
signedHex, err := eth.SignRawTransaction(rawTx, addrs[0].privateKey)
```

ERC20支持
```golang
// demo code in golang
erc20AbiHelper := eth.NewERC20InterfaceABIHelper()

// 打包erc20 合约调用转账的data
data, err := erc20AbiHelper.PackedTransfer(erc20OutEthAddr, erc20TransferValue)
// 打包erc20 合约调用查余额的data
packedErc20GetBalanceData, err : =erc20AbiHelper.PackedBalanceOf(erc20OutEthAddr)

// 通过rpc查询到erc20余额后对数据的解码(此处转换为bigInt)
erc20AbiHelper.UnpackBalanceOf(output []byte)
```

简单多重签名支持：
```golang
// demo code in golang

// :打包部署多签合约的交易数据data
createMultisigData, err := eth.PackedDeploySimpleMultiSig(eth.NewBigInt(int64(mRequired)), ownersAddrWrap, eth.NewBigInt(chainID))
// :构造部署多签合约的交易
ethtx := eth.NewETHTransactionForContractCreation(int64(a0Nonce), gasLimit, eth.NewBigInt(suggestGasPrice.Int64()), createMultisigData)

// :简单多签合约abi工具
abiHelper = eth.NewSimpleMultiSigABIHelper()
// :构造数据，用以读取合约内值（读取合约内的字段值）
callNonceData, _ := abiHelper.PackedNonce()
// :构造数据，用以调用合约的view 函数
packedGetOwnersLengthData, err := abiHelper.PackedGetOwersLength()
// :从rpc读取到合约内的值后进行解码
ownerLen, err := abiHelper.UnpackGetOwersLength(retBytes)

// :用私钥对交易信息进行签名
signRes, err := eth.UtilSimpleMultiSigExecuteSign(chainID, add.PrivkHex, multisigContractAddress, destination, executor, nonce.GetInt64(), value, gasLimit, data)
// :多签签名合约核心方法调用数据打包
packedExecuteData, err := abiHelper.PackedExecute(sigV, sigR, sigS, destAddr, value, data, &eth.ETHAddress{}, gasLimit)

```

[../qa/eth/example_simplemultisig_integration_test.go](../qa/eth/example_simplemultisig_integration_test.go) 包含了完整的golang示例代码，包括
- 部署多签合约
- 多重签名
- ERC20代币支持


## 另外的几个流行的多签合约
（star基于2019-09-03数据）
- 137star,一个简单的 2-3 实现， https://github.com/BitGo/eth-multisig-v2
- 185star,另一个简单的带每日限额的多签实现 https://github.com/ConsenSys/MultiSigWallet
- 668star,对上一个的升级 https://github.com/Gnosis/MultiSigWallet
- dapp的一个实现，https://github.com/ethereum/dapp-bin/blob/master/wallet/wallet.sol

