package internal

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
)

const (
	// ChainMainNet 链：MainNet
	ChainMainNet = 0
	// ChainTestNet3 链：TestNet3
	ChainTestNet3 = 1
	// ChainRegtest 链：Regression
	ChainRegtest = 2

	// FlagBTCUseSegWitFormat BTC使用隔离见证地址
	FlagUseSegWitFormat = "btc_use_segwit_fmt"
)

// ChainFlag2ChainParams get chainParams from const
func ChainFlag2ChainParams(chainID int) (*chaincfg.Params, error) {
	switch chainID {
	case ChainMainNet:
		return &chaincfg.MainNetParams, nil
	case ChainTestNet3:
		return &chaincfg.TestNet3Params, nil
	case ChainRegtest:
		return &chaincfg.RegressionNetParams, nil
	default:
		return nil, fmt.Errorf("期望的链选项: %d > 主链, %d > 测试链, %d > Regtest链,收到的参数：%d", ChainMainNet, ChainTestNet3, ChainRegtest, chainID)
	}
}
