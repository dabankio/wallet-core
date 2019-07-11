package btcd

import (
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/lomocoin/wallet-core/core/btc"
	"strings"
)

// ChainMainNet 链：MainNet
const ChainMainNet = 0

// ChainTestNet3 链：TestNet3
const ChainTestNet3 = 1

// ChainRegtest 链：Regression
const ChainRegtest = 2

// NewMultiSigAddress 工具函数，生成BTC多签地址
// Params:
//  chainParam: 0主链,1测试链,2Regression链
//  cmd.NRequired 需要签名的人数
//  cmd.Keys hex编码的公钥
// 限制：len(cmd.Keys) >= cmd.NRequired
// Return:
//  多签地址,redeemScript
func NewMultiSigAddress(mRequired, chainParam int, keys string) (string, error) {
	var chainParams *chaincfg.Params
	switch chainParam {
	case ChainMainNet:
		chainParams = &chaincfg.MainNetParams
	case ChainTestNet3:
		chainParams = &chaincfg.TestNet3Params
	case ChainRegtest:
		chainParams = &chaincfg.RegressionNetParams
	default:
		return "", fmt.Errorf("期望的链选项: %d > 主链, %d > 测试链, %d > Regtest链,收到的参数：%d", ChainMainNet, ChainTestNet3, ChainRegtest, chainParam)
	}
	arr := strings.Split(keys, ",")
	rs, err := btc.CreateMultiSig(btcjson.NewCreateMultisigCmd(mRequired, arr), chainParams)
	if err != nil {
		return "", err
	}
	return rs.Address + "," + rs.RedeemScript, nil
}
