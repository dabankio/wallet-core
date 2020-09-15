package core

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
)

type CoinInfo struct {
	Symbol         string
	DerivationPath accounts.DerivationPath
	MasterKey      *hdkeychain.ExtendedKey
	ChainCfg       *chaincfg.Params
}

// Coin 定义了某种货币常用的可以提供的方法
type Coin interface {
	// DeriveAddress derives the account address of the derivation path.
	DeriveAddress() (address string, err error)
	// DerivePublicKey derives the public key of the derivation path.
	DerivePublicKey() (publicKey string, err error)
	// DerivePrivateKey derives the private key of the derivation path.
	DerivePrivateKey() (privateKey string, err error)
	// DecodeTx decodes raw tx to human readable format
	DecodeTx(msg string) (tx string, err error)
	// Sign signs raw tx with privateKey
	Sign(msg, privateKey string) (sig string, err error)
	// VerifySignature verifies rawTx's signature is intact
	VerifySignature(pubKey, msg, signature string) error
}

type HasParentChain interface {
	// GetParentChainName get the symbol name of the parent chain
	GetParentChainName() string
}

// MetadataProvider we need a configuration data container per-symbol.
type MetadataProvider interface {
	GetChainID() int //获取链id
	GetPath() string //含%d的路径
	IsTestNet() bool
	GetSeed() []byte             //bip39 seed,受助记词和密码影响
	GetDerivationPath() []uint32 //推导出的bip44路径，受path, symbol影响
	HasFlag(string) bool         //是否存在flag
}
