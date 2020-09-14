package btc

import (
	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/core/btc/internal"
)

// NewBip44Deriver btc bip44 实现
func NewBip44Deriver(bip44Path string, seed []byte, chainID int) (bip44.Deriver, error) {
	coin, err := internal.New(bip44Path, seed, chainID)
	return coin, err
}
