package wallet

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dabankio/bbrpc"
	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/core/bbc"
	"github.com/dabankio/wallet-core/core/btc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBTCSegwit(t *testing.T) {
	mnemonic := "connect auto goose panda extend ozone absent climb abstract doll west crazy"
	// normal: 12Yj7jHxkQhddZVqQd697Qpq4nhEZiXAzn
	// segwit: 3BAMVfmYgroYS3XPnfk6x5TYk66p9XH2EL
	want := "3BAMVfmYgroYS3XPnfk6x5TYk66p9XH2EL"
	// public final static String BITCOIN_SEGWIT_MAIN_PATH = "m/49'/0'/0'";
	var options WalletOptions
	options.Add(WithPassword(""))
	options.Add(WithPathFormat(bip44.FullPathFormat))
	options.Add(WithFlag(FlagBTCUseSegWitFormat))
	w, err := BuildWalletFromMnemonic(
		mnemonic,
		false,
		&options,
	)
	assert.NoError(t, err)
	w.ShareAccountWithParentChain = true

	addr, err := w.DeriveAddress("BTC")
	assert.NoError(t, err)
	if addr == want {
		fmt.Println("add:", addr)
		fmt.Println(w.DerivePrivateKey("BTC"))
	}

}

// 该测试验证不同的通用参数推导出不同的地址，以确保path/password确实生效
func TestCoin_DeriveAddressOptions(t *testing.T) {
	const mnemonic = "lecture leg select like delay limit spread retire toward west grape bachelor"
	options := &WalletOptions{}
	options.Add(WithPathFormat(bip44.PathFormat))
	w, err := BuildWalletFromMnemonic(mnemonic, true, options)
	require.NoError(t, err)

	symbols := []string{"BTC", "ETH", "OMNI", "BBC", "MKF"}
	t.Run("不同path生成地址应该不一样", func(t *testing.T) {
		options := &WalletOptions{}
		options.Add(WithPathFormat(bip44.FullPathFormat))
		w2, err := BuildWalletFromMnemonic(mnemonic, true, options)
		require.NoError(t, err)
		for _, s := range symbols {
			t.Run("symbol: "+s, func(t *testing.T) {
				a1, err := w.DeriveAddress(s)
				require.NoError(t, err)

				a2, err := w2.DeriveAddress(s)
				require.NoError(t, err)

				require.NotEqual(t, a1, a2, "path 不同时推导的地址也应该不同")
			})
		}
	})
	t.Run("不同password生成地址应该不一样", func(t *testing.T) {
		options := &WalletOptions{}
		options.Add(WithPathFormat(bip44.FullPathFormat))
		options.Add(WithPassword("some_password"))
		w2, err := BuildWalletFromMnemonic(mnemonic, true, options)
		require.NoError(t, err)
		for _, s := range symbols {
			t.Run("symbol: "+s, func(t *testing.T) {
				a1, err := w.DeriveAddress(s)
				require.NoError(t, err)

				a2, err := w2.DeriveAddress(s)
				require.NoError(t, err)

				require.NotEqual(t, a1, a2, "path 不同时推导的地址也应该不同")
			})
		}
	})

	t.Run("BBC MKF共用地址在不同参数下都有效", func(t *testing.T) {
		passes := []string{bip44.Password, "", "bbc_keys"}
		paths := []string{bip44.PathFormat, bip44.FullPathFormat}

		for _, path := range paths {
			for _, pass := range passes {
				options := &WalletOptions{}
				options.Add(WithPathFormat(path))
				options.Add(WithPassword(pass))
				options.Add(WithFlag(FlagMKFUseBBCBip44ID))
				w, err := BuildWalletFromMnemonic(mnemonic, true, options)
				require.NoError(t, err)

				bbcA, err := w.DeriveAddress("BBC")
				require.NoError(t, err)
				mkfA, err := w.DeriveAddress("MKF")
				require.NoError(t, err)

				require.Equal(t, bbcA, mkfA)
			}
		}
	})

}

// 该测试验证 BTC USDT共享地址时 能始终生成一样的地址
func TestCoin_DeriveAddressPathOMNI_BTC_shareAddress(t *testing.T) {
	const mnemonic = "lecture leg select like delay limit spread retire toward west grape bachelor"

	for _, tt := range []struct {
		name, path string
	}{
		{"短路径", bip44.PathFormat},
		{"长路径", bip44.FullPathFormat},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, _tt := range []struct{ name, pass string }{
				{"其他密码", "passX"},
				{"空密码", ""},
				{"历史默认密码", bip44.Password},
			} { //使用不同的密码也应该可以共享地址
				t.Run(_tt.name, func(t *testing.T) {
					opt := &WalletOptions{}
					opt.Add(WithShareAccountWithParentChain(true))
					opt.Add(WithPassword(_tt.pass))
					opt.Add(WithPathFormat(tt.path))
					w, err := BuildWalletFromMnemonic(mnemonic, false, opt)
					require.NoError(t, err)

					btcAddr, err := w.DeriveAddress("BTC")
					require.NoError(t, err)

					omniAddr, err := w.DeriveAddress("OMNI")
					require.NoError(t, err)

					usdtAddr, err := w.DeriveAddress("USDT(Omni)")
					require.NoError(t, err)

					assert.Equal(t, btcAddr, omniAddr, "OMNI 地址错误")
					assert.Equal(t, btcAddr, usdtAddr, "USDT 地址错误")
				})
			}
		})
	}

}

// 该测试确保历史环境的逻辑兼容性，应该始终保持通过，且测试数据不应该被修改,除非你知道这意味着什么（即兼容性问题）
func TestCoin_DeriveAddress(t *testing.T) {
	const mnemonic = "lecture leg select like delay limit spread retire toward west grape bachelor"
	for _, tt := range []struct {
		name            string
		replaceMnemonic string
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
			apply: func(w *Wallet) {
				w.testNet = true
			},
		},
		{name: "OMNI default",
			symbol:  "OMNI",
			address: "1AzTauTdhZ4VKC88MAb7iu9jU3yNzpx937",
		}, //not: 13vvVPKZjsStYRZft3RyfgmCVVFsYm8nDT
		{name: "BBC default",
			symbol:  "BBC",
			address: "1zebxse3jm1c0jg0a2p22jaqyj7nerh6f1a5ck71g66j7at1w87th34gx",
		},
		{name: "BBCA",
			symbol:  "BBC",
			address: "1w6nq9qdb6dmdfwm7f2bzmxp6qgd8ytcq46t8zxr14rp57706856sern2",
			apply: func(w *Wallet) {
				w.testNet = false
				w.password = ""
				w.path = bip44.FullPathFormat
				w.flags[FlagBBCUseStandardBip44ID] = struct{}{}
			},
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
			mne := mnemonic
			if tt.replaceMnemonic != "" {
				mne = tt.replaceMnemonic
			}
			wt, err := NewHDWalletFromMnemonic(mne, "", false)
			require.NoError(t, err)
			wt.path = bip44.PathFormat
			wt.password = bip44.Password
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
	const testMnemonic = "lecture leg select like delay limit spread retire toward west grape bachelor"
	wallet := new(Wallet)

	wallet, _ = NewHDWalletFromMnemonic(testMnemonic, "", false)
	wallet.path = bip44.PathFormat
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

func TestMoneyOut(t *testing.T) { //BBC找币
	t.Skip("非测试函数") //工具函数，找币用
	mne := "slim industry rival camera chef biology charge omit forget almost craft cycle"
	shouldHasPrefix := []string{"1fkxe"}

	outAmount := 1.2
	toAddr := ""
	rpcConf := bbrpc.ConnConfig{
		Host:       "192.168.50.5:9912",
		User:       "",
		Pass:       "",
		DisableTLS: true,
	}

	paths := []string{bip44.PathFormat, bip44.FullPathFormat}
	passwords := []string{bip44.Password, ""}
	useStdBip44 := []bool{true, false}

	rpc, err := bbrpc.NewClient(&rpcConf)
	require.NoError(t, err)

	for _, path := range paths {
		for _, pass := range passwords {
			for _, stdFlag := range useStdBip44 {

				var options WalletOptions
				options.Add(WithPassword(pass))
				options.Add(WithPathFormat(path))
				if stdFlag {
					options.Add(WithFlag(FlagBBCUseStandardBip44ID))
				}
				w, err := BuildWalletFromMnemonic(
					mne,
					false,
					&options,
				)
				assert.NoError(t, err)
				addr, err := w.DeriveAddress("BBC")
				assert.NoError(t, err)

				fmt.Println("add:", addr)
				for _, pre := range shouldHasPrefix {
					if strings.HasPrefix(addr, pre) {
						fmt.Println("bingo", path, pass, stdFlag)
						fmt.Println("addxx:", addr)
					}
				}

				if 1 > 2 {
					unspents, err := rpc.Listunspent(addr, nil, 999)
					require.NoError(t, err)
					if len(unspents.Addresses[0].Unspents) > 0 {
						privk, err := w.DerivePrivateKey("BBC")
						require.NoError(t, err)
						fmt.Println("-------------", addr, privk, "-----------")

						forks, err := rpc.Listfork(false)
						require.NoError(t, err)

						tb := bbc.NewTxBuilder()
						tb = tb.
							SetAddress(toAddr).
							SetAmount(outAmount).
							SetAnchor(forks[0].Fork).
							SetFee(0.01).
							SetVersion(0).
							SetTimestamp(int(time.Now().Unix()))
						for _, utxo := range unspents.Addresses[0].Unspents {
							tb = tb.AddInput(utxo.Txid, int8(utxo.Out))
						}
						tx, err := tb.Build()
						require.NoError(t, err)

						require.NoError(t, err)
						sig, err := w.Sign("BBC", tx)
						require.NoError(t, err)
						txid, err := rpc.Sendtransaction(sig)
						require.NoError(t, err)
						fmt.Println("txid:", *txid)
						return
					}
				}
			}
		}
	}

}
func TestFindETH(t *testing.T) { //BBC找币
	shouldEqual := ""

	mnes := []string{}
	paths := []string{bip44.PathFormat, bip44.FullPathFormat}
	passwords := []string{bip44.Password, ""}

	for _, mne := range mnes {
		for _, path := range paths {
			for _, pass := range passwords {
				var options WalletOptions
				options.Add(WithPassword(pass))
				options.Add(WithPathFormat(path))
				w, err := BuildWalletFromMnemonic(
					mne,
					false,
					&options,
				)
				assert.NoError(t, err)
				addr, err := w.DeriveAddress("ETH")
				assert.NoError(t, err)

				fmt.Println("add:", addr)
				if strings.ToLower(addr) == strings.ToLower(shouldEqual) {
					fmt.Println("===> bingo:", mne, path, pass)
				}
			}
		}
	}

}

func TestNewBTCTX(t *testing.T) {
	rq := require.New(t)
	unspent := new(btc.BTCUnspent) //java: new btc.BTCUnspent()
	// unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, "")
	//9968812: 0.09968812
	unspent.Add("b4ca4f59a9ce4e2a79ed4cfc3846c33a42c3563c6b5b79903bae366e2a7b186e", int64(0), 9968812/1e8, "", "")
	amount, err := btc.NewBTCAmount(0.0001)
	rq.Nil(err)

	toAddress, err := btc.NewBTCAddressFromString("1P9xfsxQnjyHaynpk5QXMHmHvTPa8VWjFJ", btc.ChainMainNet)
	rq.Nil(err)

	outputAmount := btc.BTCOutputAmount{} //java: new btc.BTCOutputAmount()
	outputAmount.Add(toAddress, amount)

	feeRate := int64(134)

	changeAddress, err := btc.NewBTCAddressFromString("1B2YFWxCJ9v2E69MEBMNT4jJF9MY4vGnM5", btc.ChainMainNet) //找零地址
	rq.Nil(err)

	tx, err := btc.NewBTCTransaction(unspent, &outputAmount, changeAddress, feeRate, btc.ChainMainNet)
	rq.Nil(err)
	enTX, err := tx.Encode()
	rq.Nil(err)
	fmt.Println(enTX)
	fmt.Println(tx.GetFee())

}
