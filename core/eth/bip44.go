package eth

import (
	"github.com/dabankio/wallet-core/bip44"
	internal "github.com/dabankio/wallet-core/core/eth/internalized"
)

// NewBip44Deriver eth bip44 实现
func NewBip44Deriver(seed []byte) (bip44.Deriver, error) {
	coin, err := internal.New(seed)
	return coin, err
}
