VERSION=`git describe --tags --dirty`
DATE=`date +%FT%T%z`

outdir=out

module=github.com/dabankio/wallet-core

pkgBip39 = ${module}/bip39
pkgBip44 = ${module}/bip44
pkgBtc = ${module}/core/btc
pkgBBC = ${module}/core/bbc
pkgOmni = ${module}/core/omni
pkgEth = ${module}/core/eth
pkgWallet = ${module}/wallet
pkgCore = ${module}/core

pkgAll = $(pkgBip39) $(pkgBip44) $(pkgBtc) $(pkgEth) $(pkgOmni) $(pkgWallet) $(pkgBBC)

fmt:  # 格式化go代码
	@go fmt ./...

test:  # go单元测试
	@go test ./...

modTidy:
	@go mod tidy

#---------------------集成测试  start -----------------
integrationTest:
	make integrationTestBtc
	make integrationTestEth
	make integrationTestOmni
integrationTestBtc:
	#BTC 集成测试需要配置环境变量 BITCOIN_BIN_DIR 指向bitcoin-core目录
	@go test -v -tags=integration github.com/dabankio/wallet-core/qa/btc

integrationTestOmni:
	#Omni 集成测试需要配置环境变量 OMNI_BIN_PATH 指向omni-core目录
	@go test -v -tags=integration github.com/dabankio/wallet-core/qa/omni

integrationTestEth:
	#ETH 集成测试需要安装 npm i -g ganache-cli
	@go test -v -tags=integration github.com/dabankio/wallet-core/qa/eth
#---------------------集成测试  end -----------------


#---------------------构建  start -----------------
#bip39
buildBip39Android:
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/bip39.aar ${pkgBip39}
buildBip39IOS:
	gomobile bind -ldflags "-s -w" -target=ios -o=${outdir}/bip39.framework ${pkgBip39}

#bbc
buildBBCAndroid:
	gomobile bind -ldflags "-s -w" -target=android -o=${outdir}/bbc.aar ${pkgBBC} ${pkgBip44} ${pkgBip39}
buildBBCIOS:
	gomobile bind -ldflags "-s -w" -target=ios -o=${outdir}/bbc.framework ${pkgBBC} ${pkgBip44} ${pkgBip39}
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
	gomobile bind -ldflags "-s -w" -target=ios -o=${outdir}/Bip39.framework ${pkgAll}


#---------------------构建  end -----------------

#---------------------依赖图  start -----------------
depGraph: #生成工程依赖图，需要安装graphviz 和 https://github.com/loov/goda （go get github.com/loov/goda）
	@goda graph github.com/dabankio/wallet-core/...:root | dot -Tsvg -o local_graph.svg
#---------------------依赖图  end -----------------