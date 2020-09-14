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
func NewCoin(bip44Path string, seed []byte, chainID int) (core.Coin, error) {
	coin, err := internal.New(bip44Path, seed, chainID)
	return coin, err
}

// NewFromMetadata .
func NewFromMetadata(metadata core.MetadataProvider) (c core.Coin, err error) {
	return internal.NewFromMetadata(metadata)
}

var New = internal.New

// 下面这些暴露的是给bch那边用的，建议在找到更好的方案后删除，暂时没有解决bch依赖问题，先保留

// type SignRawTransactionCmd internal.SignRawTransactionCmd
// type CustomHexMsg internal.CustomHexMsg
// type RawTxInput internal.RawTxInput

// var DecodeRawTransaction = internal.DecodeRawTransaction
// var DecodeHexStr = internal.DecodeHexStr
