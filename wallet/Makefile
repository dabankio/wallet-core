outdir=out
appName=HDWalletCore
iosfilename=${appName}.framework
androidfilename=${appName}.aar
pkg=github.com/lomocoin/wallet-core/mobile github.com/lomocoin/wallet-core/mobile/btc github.com/lomocoin/wallet-core/mobile/eth

most: clean all

clean:
	if test -d $(outdir); then rm -rf $(outdir); fi

all: ios android package

ios: 
	gomobile bind -ldflags "-s -w -X github.com/lomocoin/wallet-core/mobile.buildTime=`date '+%Y-%m-%d_%I:%M:%S'` -X github.com/lomocoin/wallet-core/mobile.gitHash=`git rev-parse HEAD`" -target=ios -o=${outdir}/${iosfilename} ${pkg}

android: 
	gomobile bind -ldflags "-s -w -X github.com/lomocoin/wallet-core/mobile.buildTime=`date '+%Y-%m-%d_%I:%M:%S'` -X github.com/lomocoin/wallet-core/mobile.gitHash=`git rev-parse HEAD`" -target=android -o=${outdir}/${androidfilename} ${pkg}

package:
	zip android_${appName}_`date "+%Y%m%d_%H%M"`.zip ${outdir}/*.{aar,jar}
	zip -ry iOS_${appName}_`date "+%Y%m%d_%H%M"`.zip ${outdir}/${iosfilename}
