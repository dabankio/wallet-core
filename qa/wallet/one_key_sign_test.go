package wallet

import (
	"testing"

	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/wallet"
	"github.com/stretchr/testify/require"
)

func TestOneKeySign(t *testing.T) {
	mnemonic, err := wallet.NewMnemonic() //生成新的助记词
	require.NoError(t, err)
	t.Log("mnemonic:", mnemonic)

	for _, tt := range []struct {
		skip   bool
		name   string
		symbol string
		testFn func(*testing.T, *wallet.Wallet, ctx)
	}{
		{skip: false, symbol: "ETH", testFn: testETHPubkSign},
		{skip: false, name: "ERC20", symbol: "ETH", testFn: testERC20PubkSign},
		{skip: false, symbol: "OMNI", testFn: testOmniPubkSign},
		{skip: false, name: "BTC_p2pkh", symbol: "BTC", testFn: testBTCPubkSign},
		{skip: false, name: "BTC隔离见证", symbol: "BTC", testFn: testBTCPubkSignSegwit},
		{skip: false, symbol: "BBC", testFn: testBBCPubkSign},
		{skip: false, symbol: "MKF", testFn: testMKFPubkSign},
	} {
		if tt.skip {
			continue
		}
		if tt.name == "" {
			tt.name = tt.symbol
		}
		options := &wallet.WalletOptions{}
		options.Add(wallet.WithPathFormat(bip44.FullPathFormat))       //m/44'/%d'/0'/0/0, 确保兼容imToken
		options.Add(wallet.WithPassword(""))                           //确保兼容imToken
		options.Add(wallet.WithShareAccountWithParentChain(true))      //USDT BTC共用地址
		options.Add(wallet.WithFlag(wallet.FlagBBCUseStandardBip44ID)) //BBC使用标准bip44 ID
		options.Add(wallet.WithFlag(wallet.FlagMKFUseBBCBip44ID))      //MKF 和BBC共用地址

		if 2 < 1 {
			// options.Add(wallet.WithFlag(wallet.FlagBBCUseStandardBip44ID)) //兼容pockmine不要这个
			options.Add(wallet.WithPathFormat("m/44'/%d'"))        //pockmine的path
			options.Add(wallet.WithPassword("$pockmine的password")) //确保兼容pockmine
		}

		w, err := wallet.BuildWalletFromMnemonic(mnemonic, true, options)
		require.NoError(t, err)

		pubk, err := w.DerivePublicKey(tt.symbol)
		require.NoError(t, err)
		address, err := w.DeriveAddress(tt.symbol)
		require.NoError(t, err)

		testFn := func(c ctx, fn func(*testing.T, *wallet.Wallet, ctx)) func(*testing.T) {
			return func(t *testing.T) {
				t.Parallel()
				fn(t, w, c)
			}
		}(ctx{pubk: pubk, address: address}, tt.testFn)

		t.Run(tt.name, testFn)
	}
}
