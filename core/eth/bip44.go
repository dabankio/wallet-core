package eth

import (
	"github.com/lomocoin/wallet-core/bip44"
	"github.com/lomocoin/wallet-core/core/eth/internal"
)

// NewBip44Deriver eth bip44 实现
func NewBip44Deriver(seed []byte) (bip44.Deriver, error) {
	coin, err := internal.New(seed)
	return coin, err
}
