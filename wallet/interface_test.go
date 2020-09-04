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

// 该测试验证不同的path推导出不同的地址，以确保path确实生效
func TestCoin_DeriveAddressPath(t *testing.T) {
	const mnemonic = "lecture leg select like delay limit spread retire toward west grape bachelor"
	w, err := BuildWalletFromMnemonic(mnemonic, true, nil)
	require.NoError(t, err)

	options := &WalletOptions{}
	options.Add(WithPathFormat(bip44.FullPathFormat))
	w2, err := BuildWalletFromMnemonic(mnemonic, true, options)
	require.NoError(t, err)

	for _, s := range []string{
		"BTC", "ETH", "OMNI", "BBC", "MKF",
	} {
		t.Run("symbol: "+s, func(t *testing.T) {
			a1, err := w.DeriveAddress(s)
			require.NoError(t, err)

			a2, err := w2.DeriveAddress(s)
			require.NoError(t, err)

			require.NotEqual(t, a1, a2, "path 不同时推导的地址也应该不同")
		})
	}
}

// 该测试确保历史环境的逻辑兼容性，应该始终保持通过，且测试数据不应该被修改,除非你知道这意味着什么（即兼容性问题）
func TestCoin_DeriveAddress(t *testing.T) {
	const mnemonic = "lecture leg select like delay limit spread retire toward west grape bachelor"
	for _, tt := range []struct {
		name            string
		symbol, address string
		apply           func(*Wallet)
	}{
		{name: "ETH default",
			symbol:  "ETH",
			address: "0x947ab281Df5ec46E801F78Ad1363FaaCbe4bfd12",
		},
		{name: "BTC default",
			symbol:  "BTC",
			address: "13vvVPKZjsStYRZft3RyfgmCVVFsYm8nDT",
		},
		{name: "BTC testnet",
			symbol:  "BTC",
			address: "miSsnSQYYtt9KY3HbcQMVbyXMUraV9u9Qa",
			apply:   func(w *Wallet) { w.testNet = true },
		},
		{name: "OMNI default",
			symbol:  "OMNI",
			address: "1AzTauTdhZ4VKC88MAb7iu9jU3yNzpx937",
		}, //not: 13vvVPKZjsStYRZft3RyfgmCVVFsYm8nDT
		{name: "BBC default",
			symbol:  "BBC",
			address: "1zebxse3jm1c0jg0a2p22jaqyj7nerh6f1a5ck71g66j7at1w87th34gx",
		},
		{name: "BBC using std bip44 id",
			symbol:  "BBC",
			address: "126xdeftrb77mg6vy78zdn9rcny3zgvm9rp1wek3npqc2w8s142pfjdtz",
			apply:   func(w *Wallet) { w.flags[FlagBBCUseStandardBip44ID] = struct{}{} }},
		{name: "MKF default",
			symbol:  "MKF",
			address: "1vx6bd4d0jvhte4qndwgcf0hdc4cstmz3zqg8eh2bfsrarewv65xezpdz",
		},
		{name: "MKF share address with BBC",
			symbol:  "MKF",
			address: "1zebxse3jm1c0jg0a2p22jaqyj7nerh6f1a5ck71g66j7at1w87th34gx",
			apply:   func(w *Wallet) { w.flags[FlagMKFUseBBCBip44ID] = struct{}{} },
		},
		{name: "BBC use std bip44 id and MKF share address with BBC",
			symbol:  "MKF",
			address: "126xdeftrb77mg6vy78zdn9rcny3zgvm9rp1wek3npqc2w8s142pfjdtz",
			apply: func(w *Wallet) {
				w.flags[FlagMKFUseBBCBip44ID] = struct{}{}
				w.flags[FlagBBCUseStandardBip44ID] = struct{}{}
			},
		},
		{name: "USDT(Omni) default",
			symbol:  "USDT(Omni)",
			address: "1AzTauTdhZ4VKC88MAb7iu9jU3yNzpx937",
		}, //not: 13vvVPKZjsStYRZft3RyfgmCVVFsYm8nDT
		{name: "omni share address with btc",
			symbol:  "USDT(Omni)",
			address: "13vvVPKZjsStYRZft3RyfgmCVVFsYm8nDT",
			apply:   func(w *Wallet) { w.ShareAccountWithParentChain = true }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			wt, err := NewHDWalletFromMnemonic(mnemonic, false)
			require.NoError(t, err)
			wt.path = bip44.PathFormat
			if tt.apply != nil {
				tt.apply(wt)
			}

			addr, err := wt.DeriveAddress(tt.symbol)
			require.NoError(t, err)
			assert.Equal(t, tt.address, addr)
		})
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
