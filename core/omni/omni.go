package omni

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/lomocoin/wallet-core/bip44"
	"github.com/lomocoin/wallet-core/core"
	"github.com/lomocoin/wallet-core/core/btc"
	"github.com/pkg/errors"
)

const symbol = "OMNI"

var _ core.Coin = &omni{}
var _ core.HasParentChain = &omni{}

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

func NewWithOptions(seed []byte, testNet bool, options map[string]interface{}) (c *omni, err error) {
	c = new(omni)
	c.Symbol = symbol
	if _, ok := options["shareAccountWithParentChain"]; ok {
		c.DerivationPath, err = bip44.GetCoinDerivationPath(c.GetParentChainName())
	} else {
		c.DerivationPath, err = bip44.GetCoinDerivationPath(symbol)
	}
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

func (*omni) GetParentChainName() string {
	return core.ParentChainCofnig[symbol]
}
