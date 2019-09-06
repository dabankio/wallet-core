# BTC 使用

API使用首先建议阅读 [api文档说明](./api.md)

比特币的实现都在 `core/btc` 下

使用之前需要对比特币有基本的了解，知道UTXO模型，下面列出了一些建议的学习资料：
- [比特币官方文档](https://bitcoin.org/en/resources)
- [比特币开发者文档](https://bitcoin.org/en/developer-documentation)
- [精通比特币](https://github.com/bitcoinbook/bitcoinbook)

## 基本的流程介绍

`qa/btc`目录下包含了完整的交易示例、多签示例，可以参考其使用方式

基本的使用过程：

- 创建比特币私钥/地址，或者导入现有的私钥
- 转入资金（比如交易所渠道）
- 查询UTXO
- 使用UTXOs构造原始交易
- 使用私钥对交易签名
- 广播签名好的交易
- 等待打包，等待达到确认数量

多重签名过程：
- 收集参与成员的公钥
- 使用公钥创建多签地址，并保存多签地址和解锁脚本
- 往多签地址转入资金
- 查询多签地址的UTXO 并使用UTXO 构造原始交易
- 签名成员顺序对交易进行签名
- 达到指定签名数量后广播交易
- 等待打包，等待达到确认数


## api 举例

示例代码来自[../qa/btc/tx_integration_test.go](../qa/btc/tx_integration_test.go)

构造交易和签名：

```golang
//示例代码为golang，java/objective-c 参考生成的xxx-source.jar 和 xxx.objc.h
//对golang不熟悉可以简单忽略代码中的err

// :构造 unspent output list
unspent := new(btc.BTCUnspent) //java: new btc.BTCUnspent()
unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, utxo.RedeemScript)

// :对数量/金额的封装
amount, err := btc.NewBTCAmount(transferAmount)

// :构造地址
toAddressA1, err := btc.NewBTCAddressFromString(a1.Address, chainID)

// :构造交易输出（准确的说，只是指定了输出的地址和金额）
outputAmount := btc.BTCOutputAmount{} //java: new btc.BTCOutputAmount()
outputAmount.Add(toAddressA1, amount)

feeRate := int64(80)

// :找零地址
changeAddressA0, err := btc.NewBTCAddressFromString(a0.Address, chainID)

// :构造原始交易
tx, err = btc.NewBTCTransaction(unspent, &outputAmount, changeAddressA0, feeRate, chainID)

// :使用私钥对交易进行签名
rs, err := btc.SignTransaction(tx, a0.Privkey, chainID)
signedHex = rs.Hex
```

多重签名请参考 [../qa/btc/multisig_integration_test.go](../qa/btc/multisig_integration_test.go)