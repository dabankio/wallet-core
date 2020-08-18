生成文件说明
----

```erc20_abi_helper.go```及```SimpleMultiSig_abi_helper.go```使用```xabigen```生成，```xabigen```基于go-tehreum 的```abigen```做了简单定制，源码见 https://github.com/dabankio/go-ethereum  (branch feature/abigen)

使用的模板：`mobile_abi_helper.tpl`

在`core/eth`目录下使用的命令生成文件
- `xabigen --sol erc20.sol --pkg geth --out erc20_abi.go --tplgo mobile_abi_helper.tpl `
- `xabigen --sol internalized/contracts/SimpleMultiSig.sol -pkg eth --out SimpleMultiSig_abi.g.go --tplgo mobile_abi_helper.tpl --signal s_gen_bin,`