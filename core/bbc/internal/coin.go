package internal

import (
	"github.com/dabankio/gobbc"
	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/core"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/pkg/errors"
)

var _ core.Coin = (*Wallet)(nil) //type ensure

//SymbolService 仅和symbol有关的逻辑
type SymbolService struct {
	Symbol string
}

// DecodeTx decodes raw tx to human readable format
func (s *SymbolService) DecodeTx(msg string) (string, error) {
	return DecodeSymbolTx(s.Symbol, msg)
}

// SignTemplate signs raw tx with privateKey
func (s *SymbolService) SignTemplate(msg, templateData, privateKey string) (sig string, err error) {
	if s.Symbol == "" {
		return "", errors.New("symbol not specified")
	}
	//尝试解析为原始交易
	tx, err := gobbc.DecodeRawTransaction(SymbolSerializer(s.Symbol), msg, true)
	if err != nil {
		return msg, errors.Wrap(err, "unable to parse tx data")
	}

	err = tx.SignWithPrivateKey(SymbolSerializer(s.Symbol), templateData, privateKey)
	if err != nil {
		return msg, errors.Wrap(err, "sign failed")
	}
	return tx.Encode(SymbolSerializer(s.Symbol), true)
}

//------------------------wallet----------------------

// Wallet BBC core.Coin implementation
type Wallet struct {
	SymbolService
	DerivationPath accounts.DerivationPath
	MasterKey      *ExtendedKey
}

// NewSimpleWallet new bbc coin implementation, with short bip44 path
func NewSimpleWallet(symbol string, seed []byte) (core.Coin, error) {
	return NewWallet(symbol, seed, bip44.PathFormat, "", nil)
}

// NewWallet new bbc coin implementation, 只推导1个地址
// bip44Key 不为空时用来查找bip44 id，否则使用symbol查找
func NewWallet(symbol string, seed []byte, path string, bip44Key string, additionalDeriveParam *bip44.AdditionalDeriveParam) (core.Coin, error) {
	if e := isKnownSymbol(symbol); e != nil {
		return nil, e
	}
	var bip44ID uint32
	var err error
	if bip44Key == "" {
		bip44Key = symbol
	}
	bip44ID, err = bip44.GetCoinType(bip44Key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get bip44 id")
	}
	w := new(Wallet)
	w.Symbol = symbol
	w.DerivationPath, err = bip44.GetDerivePath(path, bip44ID, additionalDeriveParam)
	if err != nil {
		return nil, errors.Wrap(err, "bip44.GetCoinDerivationPath err:")
	}
	w.MasterKey, err = NewMaster(seed)
	if err != nil {
		return w, errors.Wrap(err, "unable to new master key for symbol")
	}
	return w, nil
}

// DeriveAddress derives the account address of the derivation path.
func (c *Wallet) DeriveAddress() (address string, err error) {
	child, err := c.derive()
	if err != nil {
		return "", err
	}
	pubk, err := gobbc.Seed2pubkString(child.key)
	if err != nil {
		return "", err
	}
	return gobbc.GetPubKeyAddress(pubk)
}

// DerivePublicKey derives the public key of the derivation path.
func (c *Wallet) DerivePublicKey() (publicKey string, err error) {
	child, err := c.derive()
	if err != nil {
		return "", err
	}
	return gobbc.Seed2pubkString(child.key)
}

func (c *Wallet) derive() (*ExtendedKey, error) {
	var err error
	childKey := c.MasterKey
	for _, childNum := range c.DerivationPath {
		childKey, err = childKey.Child(childNum)
		if err != nil {
			return nil, err
		}
	}
	return childKey, nil
}

// DerivePrivateKey derives the private key of the derivation path.
func (c *Wallet) DerivePrivateKey() (privateKey string, err error) {
	child, err := c.derive()
	if err != nil {
		return "", err
	}
	return gobbc.Seed2string(child.key), nil
}

// Sign signs raw tx with privateKey
// 首先尝试解析为带模版数据的待签数据，无法解析则尝试一般原始交易
func (c *Wallet) Sign(msg, privateKey string) (string, error) {
	var err error
	// 1尝试解析为多签数据
	if txData := tryParseTxDataWithTemplate(msg); txData != nil {
		txData.TxHex, err = c.SignTemplate(txData.TxHex, txData.TplHex, privateKey)
		if err != nil {
			return msg, errors.Wrap(err, "failed to encode tx")
		}
		return txData.EncodeString()
	}

	// 2无法解析为带模版待签数据则认为是原始交易
	return c.SignTemplate(msg, "", privateKey)
}

// VerifySignature verifies rawTx's signature is intact
func (c *Wallet) VerifySignature(pubKey, msg, signature string) error {
	return errors.New("verify signature not supported for BBC currently")
}

func tryParseTxDataWithTemplate(msg string) *gobbc.TXData {
	var data gobbc.TXData
	if err := data.DecodeString(msg); err != nil {
		return nil
	}
	return &data
}
