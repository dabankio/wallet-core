VERSION=`git describe --tags --dirty`
DATE=`date +%FT%T%z`

outdir=out

module=github.com/lomocoin/wallet-core

pkgGmtypes = ${module}/gmtypes
pkgBip39 = ${module}/bip39
pkgBip44 = ${module}/bip44
pkgBtc = ${module}/core/btc
pkgEth = ${module}/core/eth
pkgWallet = ${module}/wallet

pkgAll = $(pkgGmtypes) $(pkgBip39) $(pkgBip44) $(pkgBtc) $(pkgEth)

#如果没有指定平台，则都构建
platform?=android,ios

#构建android iOS
buildBtcAndroid:
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/btc.aar ${pkgBtc}
buildBtcIOS:
	gomobile bind -ldflags "-s -w" -target=ios -o=${outdir}/btc.framework ${pkgBtc}

buildEthAndroid:
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/eth.aar ${pkgEth}

buildAllAndroid:
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/wallet.aar ${pkgWallet} ${pkgBtc} ${pkgEth}

buildAllIOS:
	gomobile bind -ldflags "-s -w" -target=ios -o=${outdir}/wallet.framework ${pkgWallet} ${pkgBtc} ${pkgEth}

depgraph_btc:
	#导出依赖图需要安装godepgraph和graphviz
	godepgraph -horizontal -nostdlib -novendor ${pkgBtc} |dot -Tpng -o local_dep_btc.png

depgraph_eth:
	#导出依赖图需要安装godepgraph和graphviz
	godepgraph -horizontal -nostdlib -novendor ${pkgEth} |dot -Tpng -o local_dep_eth.png

depgraph_bip39:
	#导出依赖图需要安装godepgraph和graphviz
	godepgraph -horizontal -nostdlib -novendor ${pkgBip39} |dot -Tpng -o local_dep_bip39.png

depgraph_bip44:
	#导出依赖图需要安装godepgraph和graphviz
	godepgraph -horizontal -nostdlib -novendor ${pkgBip44} |dot -Tpng -o local_dep_bip44.png

fmt:
	@go fmt ./...

t:
	@echo $(platform)
