package omni

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/core"
	"github.com/dabankio/wallet-core/core/btc"
	"github.com/pkg/errors"
)

const (
	symbol = "OMNI"
	// MinNondustOutput Any standard (ie P2PKH) output smaller than this value (in satoshis) will most likely be rejected by the network.
	// This is calculated by assuming a standard output will be 34 bytes
	MinNondustOutput = 546        // satoshis
	omniHex          = "6f6d6e69" // Hex-encoded: "omni"

	OptionShareAccountWithParentChain = "shareAccountWithParentChain"
)

var _ core.Coin = &omni{}
var _ core.HasParentChain = &omni{}

type omni struct {
	btc.Btc
}

func New(seed []byte, testNet bool) (c *omni, err error) {
	return nil, errors.New("该函数已废弃，请使用 NewWithOptions 替换")
}

func NewWithOptions(path string, seed []byte, testNet bool, options map[string]interface{}) (c *omni, err error) {
	c = new(omni)
	c.Symbol = symbol

	bip44Key := symbol
	if _, ok := options[OptionShareAccountWithParentChain]; ok {
		bip44Key = c.GetParentChainName()
	}

	bip44ID, err := bip44.GetCoinType(bip44Key)
	if err != nil {
		return nil, err
	}
	c.DerivationPath, err = bip44.GetDerivePath(path, bip44ID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "bip44.GetDerivePath err:")
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
