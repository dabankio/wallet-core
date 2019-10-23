#用于生成合约-go绑定代码
abigen -sol SimpleMultiSig.sol -pkg contracts --out SimpleMultiSig.go
# abigen -sol erc20.sol -pkg contracts --out erc20.go
# solc --abi erc20.sol -o erc20.abi
abigen --abi erc20.abi/ERC20Interface.abi -pkg contracts --out erc20.go

# clean 
# rm -rf erc20.abi/
abigen -sol fixed_erc20.sol -pkg contracts --out FixedERC20.go