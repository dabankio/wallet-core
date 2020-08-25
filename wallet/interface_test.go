package wallet

import (
	"strings"
	"testing"

	"github.com/dabankio/wallet-core/bip44"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	wallet       = new(Wallet)
	testMnemonic = "lecture leg select like delay limit spread retire toward west grape bachelor"
)

func init() {
	wallet, _ = NewHDWalletFromMnemonic(testMnemonic, false)
	wallet.path = bip44.PathFormat
}

func TestCoin_DeriveAddress(t *testing.T) {
	for _, tt := range []struct {
		Symbol, Address string
		Apply           func(*Wallet)
	}{
		{"ETH", "0x947ab281Df5ec46E801F78Ad1363FaaCbe4bfd12", nil},
		{"BTC", "13vvVPKZjsStYRZft3RyfgmCVVFsYm8nDT", nil},
		{"OMNI", "1AzTauTdhZ4VKC88MAb7iu9jU3yNzpx937", nil},                                                            //not: 13vvVPKZjsStYRZft3RyfgmCVVFsYm8nDT
		{"USDT(Omni)", "1AzTauTdhZ4VKC88MAb7iu9jU3yNzpx937", nil},                                                      //not: 13vvVPKZjsStYRZft3RyfgmCVVFsYm8nDT
		{"USDT(Omni)", "13vvVPKZjsStYRZft3RyfgmCVVFsYm8nDT", func(w *Wallet) { w.ShareAccountWithParentChain = true }}, //not: 13vvVPKZjsStYRZft3RyfgmCVVFsYm8nDT
		{"BBC", "1zebxse3jm1c0jg0a2p22jaqyj7nerh6f1a5ck71g66j7at1w87th34gx", nil},
	} {
		wallet, _ = NewHDWalletFromMnemonic(testMnemonic, false)
		wallet.path = bip44.PathFormat
		if tt.Apply != nil {
			tt.Apply(wallet)
		}

		addr, err := wallet.DeriveAddress(tt.Symbol)
		require.NoError(t, err)
		assert.Equal(t, tt.Address, addr)
	}
}

func TestWallet_GetAvailableCoinList(t *testing.T) {
	bb := GetAvailableCoinList()
	t.Log(bb)
	cc := strings.Split(bb, " ")
	for i := range cc {
		addr, err := wallet.DeriveAddress(cc[i])
		assert.NoError(t, err)
		t.Log(cc[i], addr)
	}
}

func TestNewMnemonic(t *testing.T) {
	mn, err := NewMnemonic()
	assert.NoError(t, err)
	en, err := EntropyFromMnemonic(mn)
	assert.NoError(t, err)
	mn1, err := MnemonicFromEntropy(en)
	assert.NoError(t, err)
	assert.EqualValues(t, mn, mn1)
}

func TestGetVersion(t *testing.T) {
	t.Log(GetVersion())
	t.Log(GetBuildTime())
	t.Log(GetGitHash())
}

func TestIMTokenCompatibility(t *testing.T) {
	testMnemonic := "lecture leg select like delay limit spread retire toward west grape bachelor"
	var options WalletOptions
	options.Add(
		WithPassword( /*bip44.Password*/ ""),
	)
	options.Add(
		WithPathFormat("m/44'/0'/0'/0/0"),
	)
	wallet, err := BuildWalletFromMnemonic(
		testMnemonic,
		false,
		&options,
	)
	assert.NoError(t, err)
	//btc
	{
		coin, err := wallet.initCoin("BTCTest")
		assert.NoError(t, err)
		addr, err := coin.DeriveAddress()
		assert.NoError(t, err)
		imTokenBTCAddr := "1NCvbkHN9bq97JfvTGQAonNn3KpPk73LEZ"
		assert.Equal(t, imTokenBTCAddr, addr)
		t.Log(addr)
	}
	//eth
	{
		var options WalletOptions
		options.Add(
			WithPathFormat("m/44'/60'/0'/0/0"),
		)
		wallet, err := wallet.Clone(&options)
		assert.NoError(t, err)
		coin, err := wallet.initCoin("ETHTest")
		assert.NoError(t, err)
		addr, err := coin.DeriveAddress()
		assert.NoError(t, err)
		imTokenETHAddr := "0x18CACe95E0d5a3E0AC610dD8064490EdC16C176f"
		assert.Equal(t, imTokenETHAddr, addr)
		t.Log(addr)
	}
}
