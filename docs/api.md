# API 文档说明

遗憾的是由于项目处于早期的阶段，接口设计还没有稳定，暂时没有专门写api文档

不过代码中导出的api都写了注释，这些注释可以在导出的二进制包中找到，熟悉golang的话可以直接参考golang注释

android 导出xxx-sources.jar,可以找到api文档。

iOS 导出的framework 目录下有 Headers目录，可以找到导出的api文档

另外，可以参考每个币种的详细使用介绍文档

golang api 可以在 godoc 上查看： (TBD)


## 导出api

go -> java
    - golang 使用大写字母开头表示为公开的(对应java得public),导出到java后首大写字母变为小写
    - 全局函数会导出为包名类下的静态函数，比如
        - go: btc.CreateMultisigAddress(...) 导出到java为 btc.Btc.createMultisigAdress(...)
    - struct 导出到java为对应名字的class,成员函数导出到类的方法， 比如
go -> objective-c
    - 函数直接加上包名

## React Native 和 Flutter 

参考同目录下的react_native.md 和 flutter.md