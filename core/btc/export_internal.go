package btc

import (
	"github.com/dabankio/wallet-core/core"
	"github.com/dabankio/wallet-core/core/btc/internal"
)

// Btc (全部大写在导出到java那边有点问题)
type Btc struct {
	internal.BTC
}

// NewCoin btc impl of core.Coin
func NewCoin(bip44Path string, isSegWit bool, seed []byte, chainID int) (core.Coin, error) {
	coin, err := internal.New(bip44Path, isSegWit, seed, chainID)
	return coin, err
}

// NewFromMetadata .
func NewFromMetadata(metadata core.MetadataProvider) (c core.Coin, err error) {
	return internal.NewFromMetadata(metadata)
}

var New = internal.New

// FlagUseSegWitFormat BTC使用隔离见证地址
const FlagUseSegWitFormat = internal.FlagUseSegWitFormat
