 # bips 

 BIP 为Bitcoin Improvement Proposal ,比特币改进提案，虽然说是比特币，不过大部分币种都是通用的。仓库：https://github.com/bitcoin/bips

- bip39 私钥与助记词的相互转换，方便备份私钥. 下面列出了一些可供学习的资料：
    - 廖雪峰的官方网站:助记词， https://www.liaoxuefeng.com/wiki/1207298049439968/1207320517404448
    - Github上的bip39提案原文，https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki
- bip44 多币种私钥推导的方案，旨在使用单一私钥管理多个地址，下面列出了少量的学习资料：
    - Bip44 学习可能同时需要了解Bip43,Bip32
    - Bip44提案， https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki

### API 使用

**说明** 此处列出api表，仅作为简单了解用，由于文档可能有迟滞性或偏差，实际api请以导出的 源文件(android xxx-sources.jar /bip39/*.java)或头文件(iOS xxx.framework/Headers/Bip39.objc.h)作为参考

## Bip39


| 函数                 | 参数                               | 返回   | 说明                             | android                              | iOS                           |
|--------------------|----------------------------------|--------|----------------------------------|--------------------------------------|-------------------------------|
| 设置词汇表语言       | long  枚举                         | 空     | 默认英文(简繁日韩英🇫🇷🇪🇸🇮🇹) | bip39.Bip39.setWordListLang          | Bip39SetWordListLang          |
| 生成熵               | long 长度                          | byte[] | 32的倍数 128 到 256之间          | bip39.Bip39.newEntropy               | Bip39NewEntropy               |
| 助记词到熵           | 空格分隔的助记词字符串             | byte[] | 助记词异常时返回错误             | bip39.Bip39.entropyFromMnemonic      | Bip39EntropyFromMnemonic      |
| 熵到助记词           | 字节数组表示的熵                   | string | -                                | bip39.Bip39.newMnemonic              | Bip39NewMnemonic              |
| 构造种子(带错误检查) | 空格分隔的助记词字符串             | byte[] | 助记词异常时返回错误             | bip39.Bip39.newSeedWithErrorChecking | Bip39NewSeedWithErrorChecking |
| 构造种子             | 1助记词; 2密码(不指定则为空字符串) | byte[] | -                                | bip39.Bip39.newSeed                  | Bip39NewSeed                  |
| 验证助记词           | 空格分隔的助记词字符串             | bool   | -                                | bip39.Bip39.isMnemonicValid          | Bip39IsMnemonicValid          |


## Bip44

典型的用法，为特定的币种推导私钥，目前sdk支持每个主链币种推导一个地址，具体使用币种包下的：`NewBip44Deriver` 方法，不同币种可能稍有差异
- BTC
    - android `btc.Btc.NewBip44Deriver(byte[] seed, long chainID)`
    - iOS `BtcNewBip44Deriver(NSData* _Nullable seed, long chainID...)`
- ETH
    - android `btc.Btc.NewBip44Deriver(byte[] seed)`
    - iOS `BtcNewBip44Deriver(NSData* _Nullable seed...)`

bip44 目录下的 Deriver 接口定义了推导函数，目前支持单个私钥推导
```golang
DeriveAddress() (address string, err error)
DerivePublicKey() (publicKey string, err error)
DerivePrivateKey() (privateKey string, err error)
```
