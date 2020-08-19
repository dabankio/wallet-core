# Change log

## 2020-08-10
ETH 多签新增 bucketNonce 以支持并行签名
ETH 多签新增 expireTime 参数 以提高安全性

## 2020-03-20
新增币种BBC支持，https://github.com/bigbangcore/BigBang

移除了老代码中市场份额少的几个币种

好消息：现在gomobile 对go module的支持问题已经修复了 (https://github.com/golang/go/issues/27234)，可以使用go module进行构建

## v0.2 (2019-09-03)

First release, features:

- bip39 助记词,bip44 多币种私钥推导
- BTC 
    - new transaction
    - sign transaction
    - create multisig address
- Omni (same as BTC)
- ETH
    - new transaction
    - sign transaction
    - simple multisig support (deploy contract, contract call/transaction)
    - ERC20 support