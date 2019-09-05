package btc

import (
	"github.com/lomocoin/wallet-core/bip44"
	"github.com/lomocoin/wallet-core/core/btc/internal"
)

// NewBip44Deriver btc bip44 实现
func NewBip44Deriver(seed []byte, chainID int) (bip44.Deriver, error) {
	coin, err := internal.New(seed, chainID)
	return coin, err
}
