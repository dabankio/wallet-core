package main

import (
	"flag"
	"fmt"

	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/wallet"
)

var pass string
var mne string

func init() {
	flag.StringVar(&pass, "pass", "", "-pass xxx")
	flag.StringVar(&mne, "mne", "", "-mne xxx")
	flag.Parse()
}

func main() {
	const s = "BBC"

	options := &wallet.WalletOptions{}
	options.Add(wallet.WithPathFormat(bip44.FullPathFormat))       //m/44'/%d'/0'/0/0, 确保兼容imToken
	options.Add(wallet.WithPassword(pass))                         //确保兼容imToken
	options.Add(wallet.WithShareAccountWithParentChain(true))      //USDT BTC共用地址
	options.Add(wallet.WithFlag(wallet.FlagBBCUseStandardBip44ID)) //BBC使用标准bip44 ID
	options.Add(wallet.WithFlag(wallet.FlagMKFUseBBCBip44ID))      //MKF 和BBC共用地址

	w, err := wallet.BuildWalletFromMnemonic(mne, true, options)
	pe(err)

	privk, err := w.DerivePrivateKey(s)
	pe(err)
	pubk, err := w.DerivePublicKey(s)
	pe(err)
	address, err := w.DeriveAddress(s)
	pe(err)

	fmt.Println("addr", address)
	fmt.Println("pubk", pubk)
	fmt.Println("prvk", privk)
}

func pe(e error) {
	if e != nil {
		panic(e)
	}
}
