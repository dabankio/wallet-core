genMultisigGo:
	abigen -sol SimpleMultiSig.sol -pkg contracts --out SimpleMultiSig.g.go
genMultisigAbiHelper:
	xabigen --sol SimpleMultiSig.sol -pkg contracts --out SimpleMultiSig_abi.g.go --tplgo ../../mobile_abi_helper.tpl --signal s_gen_bin,

genERC20InterfaceGo:
	abigen --abi erc20.abi/ERC20Interface.abi -pkg contracts --out erc20.go
genFixedERC20:
	abigen -sol fixed_erc20.sol -pkg contracts --out FixedERC20.go