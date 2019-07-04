`SimpleMultisig.sol` originally copied from https://github.com/christianlundkvist/simple-multisig/blob/master/contracts/SimpleMultiSig.sol

## 注意
`SimpleMultiSig.go`由`abigen`从`SimpleMultisig.sol`生成，请勿手工修改

## 做了如下改动：
- pragam 升级为`^0.5.0`
- 语法上的调整（主要包括：memroy external require_message 等）
- 加入了必要的event和view函数 （eg:入账事件，转账事件，查询owners/mRequired）
- 调整常量避免冲突（未实施，参见eip-712）


##TODO
name/version/chainid/salt/等值的更改

## 参考
eip-712
- spec https://eips.ethereum.org/EIPS/eip-712
- 中文翻译 https://www.jianshu.com/p/391cffeb97b3
- 另一个eip-712的实际应用举例 https://www.jianshu.com/p/8903412db62e

go绑定的生成参见cmd.sh,以及 
- Native-DApps:-Go-bindings-to-Ethereum-contracts#bind-solidity-directly https://github.com/ethereum/go-ethereum/wiki/