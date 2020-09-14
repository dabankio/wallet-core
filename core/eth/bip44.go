package eth

import (
	"github.com/dabankio/wallet-core/bip44"
	internal "github.com/dabankio/wallet-core/core/eth/internalized"
)

// NewBip44Deriver eth bip44 实现
func NewBip44Deriver(bip44Path string, seed []byte) (bip44.Deriver, error) {
	coin, err := internal.New(bip44Path, seed)
	return coin, err
}
