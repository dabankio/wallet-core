package nxt

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/lomocoin/wallet-core/bip44"
	"github.com/lomocoin/wallet-core/core"
	"github.com/lomocoin/wallet-core/core/wcg"
	"github.com/pkg/errors"
)

const (
	symbol = "NXT"
)

type nxt struct {
	wcg.WCG
}

func New(seed []byte) (c *nxt, err error) {
	c = new(nxt)
	c.Symbol = symbol
	c.DerivationPath, err = bip44.GetCoinDerivationPath(symbol)
	if err != nil {
		err = errors.Wrap(err, "bip44.GetCoinDerivationPath err:")
		return
	}
	c.MasterKey, err = hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return
	}
	return
}

func (c *nxt) VerifySignature(pubKey, msg, signature string) error {
	// TODO
	return core.ErrThisFeatureIsNotSupported
}
