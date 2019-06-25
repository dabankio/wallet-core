package omni

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/lomocoin/HDWallet-Core/bip44"
	"github.com/lomocoin/HDWallet-Core/core/btc"
	"github.com/pkg/errors"
)

const symbol = "OMNI"

type omni struct {
	btc.BTC
}

func New(seed []byte, testNet bool) (c *omni, err error) {
	c = new(omni)
	c.Symbol = symbol
	c.DerivationPath, err = bip44.GetCoinDerivationPath(symbol)
	if err != nil {
		err = errors.Wrap(err, "bip44.GetCoinDerivationPath err:")
		return
	}
	c.ChainCfg = &chaincfg.MainNetParams
	if testNet {
		c.ChainCfg = &chaincfg.TestNet3Params
	}
	c.MasterKey, err = hdkeychain.NewMaster(seed, c.ChainCfg)
	if err != nil {
		return
	}
	return
}
