package btc

import (
	"strings"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/lomocoin/wallet-core/core/btc/internal"
)

// NewMultiSigAddress 工具函数，生成BTC多签地址
// Params:
//  chainID: 0主链,1测试链,2Regression链
//  cmd.NRequired 需要签名的人数
//  cmd.Keys hex编码的公钥
// 限制：len(cmd.Keys) >= cmd.NRequired
// Return:
//  多签地址,redeemScript
func NewMultiSigAddress(mRequired, chainID int, keys string) (string, error) {
	chainParams, err := internal.ChainFlag2ChainParams(chainID)
	if err != nil {
		return "", err
	}

	arr := strings.Split(keys, ",")
	rs, err := internal.CreateMultiSig(btcjson.NewCreateMultisigCmd(mRequired, arr), chainParams)
	if err != nil {
		return "", err
	}
	return rs.Address + "," + rs.RedeemScript, nil
}
