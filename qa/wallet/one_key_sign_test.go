package wallet

import (
	"testing"

	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/wallet"
	"github.com/stretchr/testify/require"
)

func TestOneKeySign(t *testing.T) {
	mnemonic, err := wallet.NewMnemonic()
	require.NoError(t, err)

	options := &wallet.WalletOptions{}
	options.Add(wallet.WithPathFormat(bip44.FullPathFormat))
	options.Add(wallet.WithPassword(""))
	options.Add(wallet.WithShareAccountWithParentChain(true))
	options.Add(wallet.WithFlag(wallet.FlagBBCUseStandardBip44ID))
	options.Add(wallet.WithFlag(wallet.FlagMKFUseBBCBip44ID))

	w, err := wallet.BuildWalletFromMnemonic(mnemonic, true, options)
	require.NoError(t, err)

	for _, tt := range []struct {
		skip   bool
		symbol string
		testFn func(*testing.T, *wallet.Wallet, ctx)
	}{
		{skip: false, symbol: "ETH", testFn: testERC20PubkSign},
		{skip: true, symbol: "ETH", testFn: testETHPubkSign},
		{skip: true, symbol: "OMNI", testFn: testOmniPubkSign},
		{skip: true, symbol: "BTC", testFn: testBTCPubkSign},
		{skip: true, symbol: "BBC", testFn: testBBCPubkSign},
		{skip: true, symbol: "MKF", testFn: testMKFPubkSign},
	} {
		if tt.skip {
			continue
		}
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

		t.Run(tt.symbol, testFn)
	}
}
