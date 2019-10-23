# 开发相关

## Quick Start
golang 环境

```bash
git clone xxx
cd xxx
make test
make buildAllAndroid
#...
```

android 打包需要配置android sdk

iOS 打包需要macOS平台，XCode 等等

## 基本技术参考

使用golang实现，通过[gomobile](https://github.com/golang/mobile) 打包为二进制依赖给移动端使用

gomobile 存在一些限制，不能自由的使用数据结构、数据类型，参见根目录下readme.md 的已知问题部分

## 依赖

`Makefile`里有查看依赖图的命令，对应生成local_dep_xx.png依赖图

比特币的签名库主要依赖了[btcd](https://github.com/btcsuite/btcd) 相关实现

以太坊主要依赖了[geth](https://github.com/ethereum/go-ethereum) 相关代码

omni基于比特币，所以只是在比特币的基础上增加了少量api

## 目录组织


- bip39
- bip44
- core 各币种的实现，币种目录是可以通过gomobile导出的，理念上各个币种的实现是比较独立的，如果只关心某一币种的实现，可以完全忽略其他目录
    - btc 比特币
        - internal 内部实现，go 的 internal 目录对外部目录隐藏
    - eth 以太坊
        - internal
    - omni
    - ... 其他币种
- docs 一些markdown文档
- qa 集成测试目录, `go test -tags=integration`
    - btc 比特币集成测试，需要bitcoin-core
    - eth 以太坊集成测试，需要ganache-cli
    - omni 集成测试，需要omnicore
- wallet 多币种统一封装
- out 打包目录，不包含在版本控制目录中

## 常用命令

参见`Makefile`

## 一些踩过的坑

导出指的是从go代码打包到android 或 iOS 的过程。

- 首先，参见项目根目录下的 `readme.md` “已知的问题”部分
- 虽然部分数据类型无法导出，如果需要导出的包存在公开的不被支持的 函数/数据类型/接口，可能也可以最终导出，这些不支持的类型会被跳过(skipped 在导出的xx-source.jar 和 头文件里可以看到)。最终以打包命令是否执行成功为依据。建议的方式还是尽可能避免go package 里存在不支持的类型
    - 建议提前规划好目录结构，需要导出的目录积累了少量的改动后就尝试导出
- 避免全部大写的名字，可能会影响到导出过程，甚至导出失败
