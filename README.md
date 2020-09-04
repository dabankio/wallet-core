# wallet-core

## 目录
- 概述
- 特性
- 如何使用、文档、教程
- Demo
- FAQ
- 已知问题
- 贡献代码
- 开源许可证

## 概述

这是一个加密货币签名库，核心目的在于提供离线（在线也可以）环境下加密货币的交易签名。主要使用场景偏向于移动端。

跨平台，使用golang实现，通过[gomobile](https://github.com/golang/mobile) 打包为二进制库供android(.aar) iOS(.framework)使用。(React Native / Flutter (参考实现： https://github.com/dabankio/flutter-wallet-core)

目前支持BTC Omni(USDT) ETH ,后续会加入更多币种支持。

针对的使用场景主要是移动端冷热钱包、多签钱包。

## 特性
- bip39 助记词,bip44 多币种私钥推导
- 单签
- BTC
    - 构造原始交易
    - 签名交易(不支持隔离见证地址)
    - 多重签名
- Omni 由于omni基于bitcoin,故btc支持的特性omni也支持
- ETH
    - 构造原始交易
    - 签名交易
    - 简单多重签名（合约见源码）,支持ERC20或其他任意合约调用
    - ERC20 代币支持
- 跨平台
    - android : xxx.aar
    - iOS : xxx.framework
    - 可以在ReactNative 和 Flutter 中使用
- 可单独使用需要的模块，而不会引入不需要的模块（打包体积小）
    - bip39
    - bip44
    - eth
    - btc
    - btc + omni
- 打包体积
    - android aar (打包支持 arm64-v8a   armeabi-v7a x86 x86_64)
        - 全部打包约18M(仅armeabi-v7a架构约4M)
        - 仅使用btc 13M(仅armeabi-v7a架构约3M)
        - 仅使用eth 9.7M(仅armeabi-v7a架构约2.2M)
    - iOS .framework (TBD)

### 没有的特性

- 私钥存储，建议的做法是使用android keyStore 或者 iOS keyChain
- 需要联网的功能（比如查询余额、广播交易），RPC调用需用户自行处理

---

## 如何使用

总的来说，最小情况下，只需要提供节点服务即可完成钱包app开发。

建议的方式，在使用某一币种之前，首先通过官方文档了解基本的概念，使用流程。

也可以通过集成测试(`qa/`目录)了解详细的使用方式，集成测试包含了完整的使用场景示例。

android 导出xxx-sources.jar,可以找到api文档。

iOS 导出的framework 目录下有 Headers目录，可以找到导出的api文档

### 典型的钱包实现

SDK提供：
- 私钥生成、助记词生成、助记词推导私钥、构造原始交易、对交易签名

您需要自行解决：
- 私钥存储
- 与节点或api服务器的通信，例如
    - 查询余额
    - 查询构造交易的必要参数（btc 为utxo, eth 为交易计数）
    - 广播交易
    - 查询交易
- 通知推送（如果有）


### 更多文档

[docs](docs/readme.md) 目录下包含了更多的文档，包括各币种的使用介绍，开发文档，原生平台使用等


## Demo

- 基于Flutter + 钱包sdk 的ETH 多签 + 冷钱包 demo. https://github.com/dabankio/wallet-sdk-demo


---

## FAQ

### 如何只打包我想要的模块?
参见Makefile `build`开头的命令，支持单币种导出、多币种导出或其他任意独立模块

### 如何进行代码测试?

普通单元测试 go test 即可

集成测试：

搜索：`+build integration` 为集成测试代码

运行集成测试需要增加tag:  `go test -tags=integration`

对于一般币种来说，简单的集成测试会在本地起一条链，在本地链上进行自动验证，详细可以看考`qa`目录下的实现。
另外，`Makefile`中以integrationTest 开头的测试为集成测试。

不同币种的集成测试环境要求不同
- 比特币要求配置环境变量 `BITCOIN_BIN_DIR` 指向bitcoin-core的目录
- Omni要求配置环境变量 `OMNI_BIN_PATH` 指向omni-core目录
- ETH要求本地全局安装有 ganache-cli`

测试链测试：

TBD

## 如何贡献代码
- 任何想法都可以通过issue进行讨论
- 需要新功能的话可以提issue,我们视情况添加
- 新币种支持，计划中
- PR:  fork -> feature/branch -> new PR -> flow
TBD

## 商业支持

目前暂时没有商业支持计划，有需要的话可以通过 support@dabank.io 和我们取得联系

## 关于打包体积问题

建议的方式是打包所有架构出aar,打包apk时再精简架构，这样x86虚拟机也可以调试

- android方面可以自行精简不需要的架构二进制打包,flutter 里有这种 `flutter build apk --target-platform android-arm --split-per-abi`,gradle 方面也有相关的配置 https://developer.android.com/studio/build/configure-apk-splits

- 另一种策略是出aar时就只指定一个架构E.g.`-target=android/arm,android/386.`，可以参考`gomobile help bind` (这种情况下出aar也会更快)

## 一些已知问题

总的来说gomobile并不是一个广泛使用的技术，存在诸多限制，建议阅读官方文档，并浏览现有issues: https://github.com/golang/go/issues?q=is%3Aopen+is%3Aissue+label%3Amobile+sort%3Acomments-desc

- gomobile 导出到二进制存在类型限制，导出的包的导出类型不能包含除了这些数据类型外的类型, https://godoc.org/golang.org/x/mobile/cmd/gobind#hdr-Type_restrictions
    - 补充，虽然有时候使用了不支持的类型在某些环境下还是能用，但总的来说建议参照官方说明
    - 如果需要给iOS打包，建议不要使用byte int8 uint8类型，考虑用int64等替换
- gomobile 在go module环境下打包报错的问题已经修复了，现在可以正常使用go module进行打包, ~~目前在go module环境下无法正确打包，参见：https://github.com/golang/go/issues/27234 ，推荐的做法包括~~, 
    - ~~在GOPATH目录下建立软链接，打包时通过软链接进入目录，执行make~~
- 无法同时使用多个 gomobile 导出的sdk,参见：https://github.com/golang/go/issues/15956 ，解决办法是都把源码下载到本地GOPATH 更改打包命令，多个包合并到一个二进制里(没有测试过)
- 导出的类型不要全部使用大写(比如 type BTC struct{})，建议使用驼峰规则（java方面会有点问题）

- 2020-08-07 更新，该问题在go-ethereum升级到1.9.18时修复了，不再需要处理。 ~~2020-03-23 ,go-ethereum iOS打包报错~~
    ```
    fatal error: 'libproc.h' file not found
    #include <libproc.h>
         ^~~~~~~~~~~
    ```
    - ~~参考 https://github.com/ethereum/go-ethereum/issues/20160 ， https://github.com/elastic/gosigar/pull/134~~
    - ~~这个问题已经有PR修复了，但还没有合并，可以clone到本地后用go.mod replace 替换成本地的依赖，具体参考 `go.mod`,后续PR合并后可以移除replace~~

## License

BSD-3-Clause