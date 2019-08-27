# wallet-core

## 如何使用
TBD

## 如何只打包我想要的模块
TBD

## 如何贡献代码
TBD

## 如何进行测试
TBD

## 关于打包体积问题

- android方面可以自行精简不需要的架构二进制打包,flutter 里有这种 `flutter build apk --target-platform android-arm --split-per-abi`,gradle 方面也有相关的配置 https://developer.android.com/studio/build/configure-apk-splits

## 一些已知问题

总的来说gomobile并不是一个广泛使用的技术，存在诸多限制，建议阅读官方文档，并浏览现有issues :https://github.com/golang/go/issues?q=is%3Aopen+is%3Aissue+label%3Amobile+sort%3Acomments-desc

- gomobile 导出到二进制存在类型限制，导出的包的导出类型不能包含除了这些数据类型外的类型
    - 补充，虽然有时候使用了不支持的类型在某些环境下还是能用，但总的来说建议参照官方说明
- 目前在go module环境下无法正确打包，参见：https://github.com/golang/go/issues/27234 ，推荐的做法包括
    - 在GOPATH目录下建立软链接，打包时通过软链接进入目录，执行make
    - go vendor
- 无法同时使用多个 gomobile 导出的sdk,参见：https://github.com/golang/go/issues/15956 ，解决办法是都把源码下载到本地GOPATH 更改打包命令，多个包合并到一个二进制里(没有测试过)
- 导出的类型不要全部使用大写(比如 type BTC struct{})，建议使用驼峰规则（java方面会有点问题）

## 开源许可
TBD