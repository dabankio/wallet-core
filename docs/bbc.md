## BBC 币种

简介：一个树形多链结构的区块链，理论上支持无限TPS,基于UTXO模型，POW + DPOS 混合共识模式，为物联网使用设计,不支持智能合约但支持固定模版。

浏览器： https://www.bbcexplorer.com/

### 私钥/公钥/地址

基于 ed25519/black2b 等算法

### 示例
参考 [example_bbc_test.go](../qa/bbc/example_bbc_test.go) ,搜 `<<=== sdk`

### API

- 反序列化交易, `bbc.DecodeTX`
- 使用私钥对交易进行签名 `bbc.SignWithPrivateKey`
- 私钥转公钥/地址 `bbc.ParsePrivateKey`

### 文档

releases 下载后解压，见
- android: bbc-sources.jar (tar xzvf bbc-sources.jar)
- ios: bbc.framework/Headers/Bbc.objc.h

go 里的源码
- [core/bbc/mobile.go](../core/bbc/mobile.go)