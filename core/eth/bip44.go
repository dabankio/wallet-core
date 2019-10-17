package eth

import (
	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/core/eth/internal"
)

// NewBip44Deriver eth bip44 实现
func NewBip44Deriver(seed []byte) (bip44.Deriver, error) {
	coin, err := internal.New(seed)
	return coin, err
}
