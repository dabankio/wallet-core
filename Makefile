VERSION=`git describe --tags --dirty`
DATE=`date +%FT%T%z`

outdir=out

module=github.com/lomocoin/wallet-core

pkgBip39 = ${module}/bip39
pkgBip44 = ${module}/bip44
pkgBtc = ${module}/core/btc
pkgOmni = ${module}/core/omni
pkgEth = ${module}/core/eth
pkgWallet = ${module}/wallet
pkgCore = ${module}/core

pkgAll = $(pkgBip39) $(pkgBip44) $(pkgBtc) $(pkgEth) $(pkgOmni) $(pkgWallet)

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
#bip39
buildBip39Android:
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/bip39.aar ${pkgBip39}
buildBip39IOS:
	gomobile bind -ldflags "-s -w" -target=ios -o=${outdir}/bip39.framework ${pkgBip39}

#btc
buildBtcAndroid:
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/btc.aar ${pkgBtc} ${pkgBip44} ${pkgBip39}
buildBtcIOS:
	gomobile bind -ldflags "-s -w" -target=ios -o=${outdir}/btc.framework ${pkgBtc} ${pkgBip44} ${pkgBip39}

#TODO btc+omni
buildOmniBtcAndroid:
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/btcOmni.aar ${pkgBtc} ${pkgOmni} ${pkgBip44} ${pkgBip39}
buildOmniBtcIOS:
	gomobile bind -ldflags "-s -w" -target=ios -o=${outdir}/btcOmni.framework ${pkgBtc} ${pkgOmni} ${pkgBip44} ${pkgBip39}

#eth
buildEthAndroid:
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/eth.aar ${pkgEth} ${pkgBip44} ${pkgBip39}
buildEthIOS:
	gomobile bind -ldflags "-s -w" -target=ios -o=${outdir}/eth.framework ${pkgEth} ${pkgBip44} ${pkgBip39}

#all: bip39,bip44,btc,eth,omni
buildAllAndroid:
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/wallet.aar ${pkgAll}
buildAllIOS:
	gomobile bind -ldflags "-s -w" -target=ios -o=${outdir}/wallet.framework ${pkgAll}


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