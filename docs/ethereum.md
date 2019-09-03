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


## 另外的几个流行的多签合约
（star基于2019-09-03数据）
- 137star,一个简单的 2-3 实现， https://github.com/BitGo/eth-multisig-v2
- 185star,另一个简单的带每日限额的多签实现 https://github.com/ConsenSys/MultiSigWallet
- 668star,对上一个的升级 https://github.com/Gnosis/MultiSigWallet
- 137star,https://github.com/BitGo/eth-multisig-v2
- dapp的一个实现，https://github.com/ethereum/dapp-bin/blob/master/wallet/wallet.sol
