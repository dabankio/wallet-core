package wallet

import (
	"strings"

	"github.com/dabankio/wallet-core/core"
	// "github.com/dabankio/wallet-core/core/bch"
	"github.com/dabankio/wallet-core/core/bbc"
	"github.com/dabankio/wallet-core/core/btc"
	"github.com/dabankio/wallet-core/core/eth"
	"github.com/dabankio/wallet-core/core/omni"
	"github.com/dabankio/wallet-core/core/trx"
	"github.com/dabankio/wallet-core/core/xrp"
	"github.com/pkg/errors"
)

func (c Wallet) initCoin(symbol string) (coin core.Coin, err error) {
	if len(c.seed) == 0 {
		err = errors.New("seed is empty")
		return
	}
	md, err := c.Metadata(symbol)
	if err != nil {
		return
	}
	switch symbol {
	case "BBC":
		coin, err = bbc.NewCoin(c.seed, c.path)
	case "BTC":
		if c.testNet {
			coin, err = btc.New(c.seed, btc.ChainTestNet3)
		} else {
			coin, err = btc.New(c.seed, btc.ChainMainNet)
		}
	case "BTCTest":
		coin, err = btc.NewFromMetadata(md)
	case "USDT(Omni)", "OMNI":
		// TODO more elegant way to support custom options, make the wallet instance a argument?
		if c.ShareAccountWithParentChain {
			coin, err = omni.NewWithOptions(c.seed, c.testNet, map[string]interface{}{
				"shareAccountWithParentChain": true,
			})
		} else {
			coin, err = omni.New(c.seed, c.testNet)
		}
	// case "BCH": //TODO BCH 对 BTC 的代码依赖问题暂时没有解决，先注释掉
	// coin, err = bch.New(c.seed, c.testNet)
	case "ETH", "XT", "THM", "ALI", "RED", "USO", "BTK", "EGT", "HOTC(HOTCOIN)":
		coin, err = eth.New(c.seed)
	case "ETHTest":
		coin, err = eth.NewFromMetadata(md)
	case "XRP":
		coin, err = xrp.New(c.seed)
	case "TRX", "BTT":
		coin, err = trx.New(c.seed)
	default:
		err = errors.Errorf("no entry for coin (%s) was found.", symbol)
	}
	return
}

// GetAvailableCoinList 获取支持币种列表 (case sensitive)
// return : "BTC LMC ETH WCG"
func GetAvailableCoinList() string {
	availableCoin := []string{
		// BTC series
		"BTC",
		// "BCH",//TODO BCH 对 BTC 的代码依赖问题暂时没有解决，先注释掉
		// "MGD",
		// "LMC",

		// OMNI series
		"USDT(Omni)",
		"OMNI",

		// ETH series
		"ETH",
		// ERC20 series
		// "XT",
		// "THM",
		// "ALI",
		// "RED",
		// "USO",
		// "BTK",
		// "EGT",
		// "HOTC(HOTCOIN)",

		// ripple
		// "XRP",

		// tron
		// "TRX",
		// TRC10
		// "BTT",
	}
	return strings.Join(availableCoin, " ")
}
