package mgd

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/lomocoin/HDWallet-Core/bip44"
	"github.com/lomocoin/HDWallet-Core/core/btc"
	"github.com/pkg/errors"
)

const symbol = "MGD"

type mgd struct {
	btc.BTC
}

func New(seed []byte) (c *mgd, err error) {
	c = new(mgd)
	c.Symbol = symbol
	c.DerivationPath, err = bip44.GetCoinDerivationPath(symbol)
	if err != nil {
		err = errors.Wrap(err, "bip44.GetCoinDerivationPath err:")
		return
	}
	c.ChainCfg = &chaincfg.MainNetParams
	c.ChainCfg.PubKeyHashAddrID = 0x32 // starts with M
	c.ChainCfg.ScriptHashAddrID = 0x26 // starts with G
	c.ChainCfg.PrivateKeyID = 0x19     // starts with 4 (compressed)
	c.MasterKey, err = hdkeychain.NewMaster(seed, c.ChainCfg)
	if err != nil {
		return
	}
	return
}
