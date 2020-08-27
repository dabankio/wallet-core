package wallet

import (
	"fmt"
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
	for _, tt := range []struct {
		skip                       bool
		name, mnemonic, pass, path string
		addrs                      map[string]string
	}{
		{
			name:     "legacy wallet",
			mnemonic: "lecture leg select like delay limit spread retire toward west grape bachelor",
			pass:     bip44.Password,
			path:     bip44.PathFormat,
			addrs: map[string]string{
				"BTC": "13vvVPKZjsStYRZft3RyfgmCVVFsYm8nDT",
				"ETH": "0x947ab281Df5ec46E801F78Ad1363FaaCbe4bfd12",
			},
		},
		{
			name:     "imToken wallet",
			mnemonic: "lecture leg select like delay limit spread retire toward west grape bachelor",
			pass:     "",
			path:     bip44.FullPathFormat,
			addrs: map[string]string{
				"BTC": "1NCvbkHN9bq97JfvTGQAonNn3KpPk73LEZ",
				"ETH": "0x18CACe95E0d5a3E0AC610dD8064490EdC16C176f",
			},
		},
		{
			name:     "legacy wallet2",
			mnemonic: "connect auto goose panda extend ozone absent climb abstract doll west crazy",
			pass:     bip44.Password,
			path:     bip44.PathFormat,
			addrs: map[string]string{
				"BTC": "12X2swpFCeeoVVofn6UHaRpfDAiH9ew2U6",
				"ETH": "0x5f7838c98581f48b9Dc77Cd6410D37AEeAA1e14B",
			},
		},
		{
			name:     "imToken wallet2",
			mnemonic: "connect auto goose panda extend ozone absent climb abstract doll west crazy",
			pass:     "",
			path:     bip44.FullPathFormat,
			addrs: map[string]string{
				"BTC": "12Yj7jHxkQhddZVqQd697Qpq4nhEZiXAzn",
				"ETH": "0xf90b1d47964149Ab7F815F1564E0f41Cac0Dc456",
			},
		},
	} {
		if tt.skip {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			var options WalletOptions
			options.Add(WithPassword(tt.pass)) /*bip44.Password*/
			options.Add(WithPathFormat(tt.path))
			wallet, err := BuildWalletFromMnemonic(
				tt.mnemonic,
				false,
				&options,
			)
			assert.NoError(t, err)
			for symbol, addr := range tt.addrs {
				deriveAddr, err := wallet.DeriveAddress(symbol)
				require.NoError(t, err, fmt.Sprintf("symbol:%s", symbol))
				require.Equal(t, addr, deriveAddr)
			}
		})
	}
}
