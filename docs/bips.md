# bips

BIP 为 Bitcoin Improvement Proposal ,比特币改进提案，虽然说是比特币，不过大部分币种都是通用的。仓库：https://github.com/bitcoin/bips

- bip39 私钥与助记词的相互转换，方便备份私钥. 下面列出了一些可供学习的资料：
  - 廖雪峰的官方网站:助记词， https://www.liaoxuefeng.com/wiki/1207298049439968/1207320517404448
  - Github 上的 bip39 提案原文，https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki
- bip44 多币种私钥推导的方案，旨在使用单一私钥管理多个地址，下面列出了少量的学习资料：
  - Bip44 学习可能同时需要了解 Bip43,Bip32
  - Bip44 提案， https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki

### API 使用

**说明** 此处列出 api 表，仅作为简单了解用，由于文档可能有迟滞性或偏差，实际 api 请以导出的 源文件(android xxx-sources.jar /bip39/\*.java)或头文件(iOS xxx.framework/Headers/Bip39.objc.h)作为参考

## Bip39

| 函数                 | 参数                                 | 返回   | 说明                              | android / iOS                              iOS                           |
|--------------------|------------------------------------|--------|-----------------------------------|---------------------------------------------------------------------|
| 设置词汇表语言       | long 枚举                            | 空     | 默认英文(简繁日韩英 🇫🇷🇪🇸🇮🇹) | bip39.Bip39.setWordListLang          <br> Bip39SetWordListLang          |
| 生成熵               | long 长度                            | byte[] | 32 的倍数 128 到 256 之间         | bip39.Bip39.newEntropy               <br> Bip39NewEntropy               |
| 助记词到熵           | 空格分隔的助记词字符串               | byte[] | 助记词异常时返回错误              | bip39.Bip39.entropyFromMnemonic      <br> Bip39EntropyFromMnemonic      |
| 熵到助记词           | 字节数组表示的熵                     | string | -                                 | bip39.Bip39.newMnemonic              <br> Bip39NewMnemonic              |
| 构造种子(带错误检查) | 空格分隔的助记词字符串               | byte[] | 助记词异常时返回错误              | bip39.Bip39.newSeedWithErrorChecking <br> Bip39NewSeedWithErrorChecking |
| 构造种子             | 1 助记词; 2 密码(不指定则为空字符串) | byte[] | -                                 | bip39.Bip39.newSeed                  <br> Bip39NewSeed                  |
| 验证助记词           | 空格分隔的助记词字符串               | bool   | -                                 | bip39.Bip39.isMnemonicValid          <br> Bip39IsMnemonicValid          |

## Bip44

典型的用法，为特定的币种推导私钥，目前 sdk 支持每个主链币种推导一个地址，具体使用币种包下的：`NewBip44Deriver` 方法，不同币种可能稍有差异

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

## 典型用法

**创建Seed**
- 创建熵 bip39.NewEntropy 
- 设置bip39语言 bip39.SetWordListLang
- 把熵转换为助记词 bip39.NewMnemonic
- 将助记词转换为 bip39.NewSeed


- 使用seed为各币种推导私钥/公钥/地址
  - eg, NewBip44Deriver(seed)

- 验证助记词
