package bbc

import (
	"crypto/ed25519"

	"github.com/dabankio/gobbc"
	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/core"
	"github.com/dabankio/wallet-core/core/bbc/internal"
	"github.com/pkg/errors"
)

const (
	SymbolBBC = "BBC"
	SymbolMKF = "MKF"
)

func NewCoin(seed []byte, path string) (core.Coin, error) {
	return internal.NewCoinWithPath(SymbolBBC, seed, path)
}

// NewSymbolCoin symbol 支持 兼容BBC的币种(比如MKF)
func NewSymbolCoin(symbol string, seed []byte, path string) (core.Coin, error) {
	return internal.NewCoinWithPath(symbol, seed, path)
}

// NewSimpleBip44Deriver 根据种子获取bip44推导,仅推导1个
func NewSimpleBip44Deriver(seed []byte) (bip44.Deriver, error) {
	return internal.NewCoin(SymbolBBC, seed)
}

// NewSymbolSimpleBip44Deriver 根据种子获取bip44推导,仅推导1个
func NewSymbolSimpleBip44Deriver(symbol string, seed []byte) (bip44.Deriver, error) {
	return internal.NewCoin(symbol, seed)
}

// NewBip44Deriver 根据种子获取bip44推导
// accountIndex 账户索引，以0开始
// changeType 0:外部使用， 1:找零， 通常使用0,BBC通常找零到发送地址
// index 地址索引，以0开始
func NewBip44Deriver(seed []byte, accountIndex, changeType, index int) (bip44.Deriver, error) {
	return internal.NewCoinFullPath(SymbolBBC, seed, accountIndex, changeType, index)
}

func NewSymbolBip44Deriver(symbol string, seed []byte, accountIndex, changeType, index int) (bip44.Deriver, error) {
	return internal.NewCoinFullPath(symbol, seed, accountIndex, changeType, index)
}

// DeriveKeySimple 推导路径 m/44'/%d'
func DeriveKeySimple(seed []byte) (*KeyInfo, error) {
	return DeriveSymbolKeySimple(SymbolBBC, seed)
}

// DeriveSymbolKeySimple 推导路径 m/44'/%d'
func DeriveSymbolKeySimple(symbol string, seed []byte) (*KeyInfo, error) {
	var info KeyInfo
	coin, err := internal.NewCoin(symbol, seed)
	if err != nil {
		return &info, errors.Wrap(err, "无法创建bip44实现")
	}
	privateKey, err := coin.DerivePrivateKey()
	if err != nil {
		return &info, errors.Wrap(err, "DerivePrivateKey failed")
	}
	return ParsePrivateKey(privateKey)
}

// DeriveKey 该函数后面3个参数无效，等同于 DeriveKeySimple，仅保留
func DeriveKey(seed []byte, accountIndex, changeType, index int) (*KeyInfo, error) {
	return nil, errors.New("该函数已失效，请使用 DeriveKeySimple 替换 ,accountIndex, changeType, index 3个参数旧版api也是不会生效的的")
}

// DecodeTX 该函数已废弃，请使用 DecodeSymbolTX
func DecodeTX(rawTX string) (string, error) {
	return "", errors.New("该函数已废弃，请使用 DecodeSymbolTX")
}

// DecodeSymbolTX 解析原始交易（使用JSON RPC createtransaction 创建的交易）,symbol: BBC | MKF
func DecodeSymbolTX(symbol, rawTX string) (string, error) {
	return internal.DecodeSymbolTx(symbol, rawTX)
}

// SignWithPrivateKey 使用私钥对原始交易进行签名,
// 关于templateData的使用参考 https://github.com/dabankio/gobbc/blob/d51d596fa310a5778e3d11eb59bc66d1a6a5e3d6/transaction.go#L197 （SignWithPrivateKey部分）
// 参考测试用例 qa/bbc/example_bbc_test.go
func SignWithPrivateKey(rawTX, templateData, privateKey string) (string, error) {
	bbc := &internal.BBC{Symbol: SymbolBBC}
	return bbc.SignTemplate(rawTX, templateData, privateKey)
}

// SymbolSignWithPrivateKey 指定币种使用私钥对交易签名
func SymbolSignWithPrivateKey(symbol, rawTX, templateData, privateKey string) (string, error) {
	bbc := &internal.BBC{Symbol: symbol}
	return bbc.SignTemplate(rawTX, templateData, privateKey)
}

// KeyInfo 私钥，公钥，地址
type KeyInfo struct {
	PrivateKey, PublicKey, Address string
}

// ParsePrivateKey 解析私钥，返回 privateKey,publicKey,address
func ParsePrivateKey(privateKey string) (*KeyInfo, error) {
	info := KeyInfo{PrivateKey: privateKey}
	ed25Privk, err := gobbc.ParsePrivkHex(privateKey)
	if err != nil {
		return &info, errors.New("解析私钥失败")
	}

	info.PublicKey = gobbc.CopyReverseThenEncodeHex(ed25Privk.Public().(ed25519.PublicKey))
	info.Address, err = gobbc.GetPubKeyAddress(info.PublicKey)
	if err != nil {
		return &info, errors.Wrap(err, "公钥转地址失败")
	}
	return &info, nil
}

// Address2pubk 将地址转换为公钥
func Address2pubk(address string) (string, error) {
	return gobbc.ConvertAddress2pubk(address)
}
