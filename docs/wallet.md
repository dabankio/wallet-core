# wallet 导出 API

wallet 目录作为多币种的统一抽象，抽离出几个简单的API使用， 具体来说参考：[interface.go](../wallet/interface.go) 中 `struct wallet`的几个函数

一般性用法：生成或导入助记词 -> 创建钱包 -> 推导地址/公钥/私钥 -> 使用私钥进行签名

## 交易创建和格式
不同币种交易格式不同，单签和多签也不同，具体请参考各币种的 Sign 函数

下面简单说明几个币种的单签格式， 单签可以参考集成测试代码`qa/wallet/xxx.go`(包含完整的交易过程、数据结构等)

BTC/OMNI，你需要将交易组装为下面的示例数据，然后json encode -> hex eoncode， 交易通常可以用RPC创建
```json
{
    //交易相关的输入
    "Inputs": [
        {
            "scriptPubKey": "76a91481e8446570f9f0db967bfee69cda1069fa42653588ac",
            "txid": "8225f6d0c3b4d44d81b1986c04458725250d30708728455860dc8dbea6656f78",
            "vout": 1
        }
    ],
    //原始交易
    "RawTx": "0200000001786f65a6be8ddc605845288770300d25258745046c98b1814dd4b4c3d0f625820100000000ffffffff03fe0e651e000000001976a91481e8446570f9f0db967bfee69cda1069fa42653588ac0000000000000000166a146f6d6e690000000080000003000000000000006522020000000000001976a9142f49fb7ce9b9de814a65f393fcb8ee69b878443488ac00000000"
}
```

ETH/ERC20, 单签简单采用 eth transaction RLP 编码，再 hex 编码。 ETH节点通常不提供创建交易的RPC,所以你可能需要自行构建交易数据，建议的做法是找到开发语言的相关库， 或者直接使用go封装为服务

BBC/MKF, 交易使用RPC创建，可以直接用来签名， SDK会尝试解析为多签交易数据或原始交易数据

## 配置和参数

推荐使用 `wallet/builder.go` 中的建造模式创建钱包，可提供的选项通常类似 `WithXXX`, 下面列出已有的选项（文档可能不能保持同步，细节请参考生成的文档)
- WithShareAccountWithParentChain 使OMNI 与 BTC 生成一样的地址
- WithPathFormat，指定BIP44路径，如果不提供的话默认使用`m/44'/%d'`, 我们建议你使用 `m/44'/%d'/0'/0/0` 这与imtoken兼容
- WithFlag, 一些特殊的开关，具体参考 `wallet/flags.go` 包括
    - `FlagBBCUseStandardBip44ID (bbc_use_std_bip44_id)`, BBC使用标准bip44 id (默认不是标准bip44 id)
    - `FlagMKFUseBBCBip44ID (mkf_use_bbc_bip44_id)`, MKF使用BBC的bip44 id (即MKF和BBC共用地址)
- WithPassword, 提供salt, 注意：salt 是的你定制了推导参数，变得"不在标准"，这可能会影响与其他钱包的兼容性, 建议使用 "" (空字符串)
