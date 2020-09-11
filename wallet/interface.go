// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/dabankio/wallet-core/bip39"
	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/core"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/pkg/errors"
)

type Wallet struct {
	mnemonic string
	// ShareAccountWithParentChain break the HD rule, use the metadata of the parent chain to generate keys and addresses.
	ShareAccountWithParentChain bool
	seed                        []byte
	testNet                     bool
	password                    string
	path                        string
	flags                       map[string]struct{} //存在一些特殊情况，使用通用的配置可能产生级连影响，所有用了flag以实现灵活的配置，缺点是逻辑比较分散。 (举个例子:ShareAccountWithParentChain来控制BTC和OMNI用1个地址，但如果这时候BBC和MKF不需要用同一个地址则会有问题)
}

// 检查并设置默认值
func (c *Wallet) init() error {
	if len(c.seed) == 0 {
		return errors.New("seed not provided")
	}
	if c.path == "" {
		c.path = bip44.PathFormat //默认使用短格式bip44 path
	}
	return nil
}

// HasFlag 是否存在flag
func (c *Wallet) HasFlag(flag string) bool { _, ok := c.flags[flag]; return ok }

// MnemonicFromEntropy 根据 entropy， 获取对应助记词
func MnemonicFromEntropy(entropy string) (mnemonic string, err error) {
	entropyBytes, err := hex.DecodeString(entropy)
	if err != nil {
		return
	}
	return core.NewMnemonic(entropyBytes)
}

// EntropyFromMnemonic 根据 助记词, 获取 Entropy
// returns the input entropy used to generate the given mnemonic
func EntropyFromMnemonic(mnemonic string) (entropy string, err error) {
	entropyBytes, err := core.EntropyFromMnemonic(mnemonic)
	if err != nil {
		return
	}
	entropy = hex.EncodeToString(entropyBytes)
	return
}

// NewMnemonic 生成助记词
// 默认使用128位密钥，生成12个单词的助记词
func NewMnemonic() (mnemonic string, err error) {
	const bitSize = 128
	entropy, err := core.NewEntropy(bitSize)
	if err != nil {
		return
	}
	return core.NewMnemonic(entropy)
}

// ValidateMnemonic 验证助记词的正确性
func ValidateMnemonic(mnemonic string) (err error) {
	_, err = core.NewSeedFromMnemonic(mnemonic)
	if err != nil {
		return
	}
	return
}

// NewMnemonic 通过助记词得到一个 HD 对象
func NewHDWalletFromMnemonic(mnemonic string, testNet bool) (w *Wallet, err error) {
	seed, err := core.NewSeedFromMnemonic(mnemonic)
	if err != nil {
		return
	}
	w = new(Wallet)
	w.mnemonic = mnemonic
	w.seed = seed
	w.testNet = testNet
	w.flags = make(map[string]struct{})
	if e := w.init(); e != nil {
		return nil, e
	}
	return
}

// DeriveAddress 获取对应币种代号的 地址
func (c Wallet) DeriveAddress(symbol string) (address string, err error) {
	coin, err := c.initCoin(symbol)
	if err != nil {
		return "", errors.Wrap(err, "init coin impl err")
	}
	return coin.DeriveAddress()
}

// DerivePublicKey 获取对应币种代号的 公钥
func (c Wallet) DerivePublicKey(symbol string) (publicKey string, err error) {
	coin, err := c.initCoin(symbol)
	if err != nil {
		return
	}
	return coin.DerivePublicKey()
}

// DerivePrivateKey 获取对应币种代号的 私钥
func (c Wallet) DerivePrivateKey(symbol string) (privateKey string, err error) {
	coin, err := c.initCoin(symbol)
	if err != nil {
		return
	}
	return coin.DerivePrivateKey()
}

// DecodeTx 解析交易数据
// return: json 数据
func (c Wallet) DecodeTx(symbol, msg string) (tx string, err error) {
	coin, err := c.initCoin(symbol)
	if err != nil {
		return
	}
	return coin.DecodeTx(msg)
}

// Sign 签名交易
func (c Wallet) Sign(symbol, msg string) (sig string, err error) {
	coin, err := c.initCoin(symbol)
	if err != nil {
		return
	}

	privateKey, err := coin.DerivePrivateKey()
	if err != nil {
		return
	}

	return coin.Sign(msg, privateKey)
}

func (c Wallet) Metadata(symbol string) (core.MetadataProvider, error) {
	seed, err := bip39.NewSeedWithErrorChecking(c.mnemonic, c.password)
	if err != nil {
		return nil, err
	}
	path := c.path
	if strings.Contains(path, "%d") {
		symbolBip44ID, err := bip44.GetCoinType(symbol)
		if err != nil {
			return nil, errors.Wrap(err, "get coin bip44 id failed")
		}
		path = fmt.Sprintf(path, symbolBip44ID)
	}
	derivationPath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return nil, err
	}
	md := metadataProviderImpl{
		symbol:         symbol,
		path:           c.path,
		testNet:        c.testNet,
		seed:           seed,
		derivationPath: derivationPath,
	}
	return &md, nil
}
