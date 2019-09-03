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
pkgCore = ${module}/core

pkgAll = $(pkgGmtypes) $(pkgBip39) $(pkgBip44) $(pkgBtc) $(pkgEth)

#如果没有指定平台，则都构建
platform?=android,ios

fmt:  # 格式化go代码
	@go fmt ./...

test:  # go单元测试
	@go test ./...

#---------------------集成测试  start -----------------

integrationTestBtc:
	#BTC 集成测试需要配置环境变量 BITCOIN_BIN_DIR 指向bitcoin-core目录
	@go test -v -tags=integration github.com/lomocoin/wallet-core/qa/btc

integrationTestOmni:
	#Omni 集成测试需要配置环境变量 OMNI_BIN_PATH 指向omni-core目录
	@go test -v -tags=integration github.com/lomocoin/wallet-core/qa/omni

integrationTestEth:
	#ETH 集成测试需要安装 npm i -g ganache-cli
	@go test -v -tags=integration github.com/lomocoin/wallet-core/qa/eth
#---------------------集成测试  end -----------------


#---------------------构建  start -----------------
#构建android iOS
buildBtcAndroid:
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/btc.aar ${pkgBtc}
buildBtcIOS:
	gomobile bind -ldflags "-s -w" -target=ios -o=${outdir}/btc.framework ${pkgBtc}


buildEthAndroid:
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/eth.aar ${pkgEth}
buildEthIOS:
	gomobile bind -ldflags "-s -w" -target=ios -o=${outdir}/eth.framework ${pkgEth}


buildAllAndroid:
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/wallet.aar ${pkgWallet} ${pkgBtc} ${pkgEth}
buildAllIOS:
	gomobile bind -ldflags "-s -w" -target=ios -o=${outdir}/wallet.framework ${pkgWallet} ${pkgBtc} ${pkgEth}


#---------------------构建  end -----------------

#---------------------依赖图  start -----------------
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

depgraph_core:
	#导出依赖图需要安装godepgraph和graphviz
	godepgraph -horizontal -nostdlib -novendor ${pkgCore} |dot -Tpng -o local_dep_core.png

#---------------------依赖图  end -----------------