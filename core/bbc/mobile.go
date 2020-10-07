package bbc

import (
	"crypto/ed25519"

	"github.com/dabankio/gobbc"
	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/core"
	"github.com/dabankio/wallet-core/core/bbc/internal"
	"github.com/pkg/errors"
)

// .
const (
	SymbolBBC = internal.SymbolBBC
	SymbolMKF = internal.SymbolMKF

	templateDexorder = 9
)

// FullnameMap .
var FullnameMap = map[string]string{
	SymbolBBC: "BigBang Core",
	SymbolMKF: "MarketFinance",
}

// NewCoin 该函数已废弃，请使用NewSymbolCoin
func NewCoin(seed []byte, path string) (core.Coin, error) {
	return nil, errors.New("该函数已废弃，请使用NewSymbolCoin")
}

// NewSymbolCoin symbol 支持 兼容BBC的币种(比如MKF)
func NewSymbolCoin(symbol string, seed []byte, path string, bip44Key string) (core.Coin, error) {
	return internal.NewWallet(symbol, seed, path, bip44Key, nil)
}

// NewSimpleBip44Deriver 根据种子获取bip44推导,仅推导1个
// func NewSimpleBip44Deriver(seed []byte) (bip44.Deriver, error) {
// 	return nil, errors.New("该函数已废弃，请使用NewSymbolSimpleBip44Deriver")
// }

// NewSymbolSimpleBip44Deriver 根据种子获取bip44推导,仅推导1个
// func NewSymbolSimpleBip44Deriver(symbol string, seed []byte) (bip44.Deriver, error) {
// 	return internal.NewSimpleWallet(symbol, seed)
// }

// NewBip44Deriver 该函数已废弃，请使用NewSymbolBip44Deriver
func NewBip44Deriver(seed []byte, accountIndex, changeType, index int) (bip44.Deriver, error) {
	return nil, errors.New("该函数已废弃，请使用NewSymbolBip44Deriver")
}

// NewSymbolBip44Deriver 获取bip44推导
// accountIndex 账户索引，以0开始
// changeType 0:外部使用， 1:找零， 通常使用0,BBC通常找零到发送地址
// index 地址索引，以0开始
func NewSymbolBip44Deriver(symbol string, bip44Path string, bip44Key string, seed []byte, accountIndex, changeType, index int) (bip44.Deriver, error) {
	return internal.NewWallet(symbol, seed, bip44Path, bip44Key, &bip44.AdditionalDeriveParam{
		AccountIndex: accountIndex, ChangeType: changeType, Index: index,
	})
}

// DeriveKeySimple 该函数已废弃，请使用NewSymbolCoin
func DeriveKeySimple(seed []byte) (*KeyInfo, error) {
	return nil, errors.New("该函数已废弃，请使用NewSymbolCoin")
}

// DeriveSymbolKeySimple 该函数已废弃，请使用NewSymbolCoin
func DeriveSymbolKeySimple(symbol string, seed []byte) (*KeyInfo, error) {
	return nil, errors.New("该函数已废弃，请使用NewSymbolCoin")
}

// DeriveKey 该该函数已废弃，请使用NewSymbolCoin
func DeriveKey(seed []byte, accountIndex, changeType, index int) (*KeyInfo, error) {
	return nil, errors.New("该函数已废弃，请使用NewSymbolCoin")
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
	return "", errors.New("该函数已废弃，请使用SymbolSignWithPrivateKey")
}

// SymbolSignWithPrivateKey 指定币种使用私钥对交易签名
func SymbolSignWithPrivateKey(symbol, rawTX, templateData, privateKey string) (string, error) {
	service := &internal.SymbolService{Symbol: symbol}
	return service.SignTemplate(rawTX, templateData, privateKey)
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

// TemplateInfo 简要模版信息
type TemplateInfo struct {
	//Type 类型
	Type int
	//Address 地址
	Address string
	//RawHex hex编码的原始数据（TemplateData）
	RawHex string
}

// CreateTemplateDataDexOrder 获取dexOrder模版数据
func CreateTemplateDataDexOrder(
	sellerAddress string,
	coinpair string,
	price int32,
	fee int32,
	recvAddress string,
	validHeight int32,
	matchAddress string,
	dealAddress string,
) (*TemplateInfo, error) {
	add, raw, err := gobbc.CreateTemplateDataDexOrder(gobbc.DexOrderParam{
		SellerAddress: gobbc.Address(sellerAddress),
		Coinpair:      coinpair,
		Price:         price,
		Fee:           fee,
		RecvAddress:   recvAddress,
		ValidHeight:   validHeight,
		MatchAddress:  gobbc.Address(matchAddress),
		DealAddress:   gobbc.Address(dealAddress),
	})
	if err != nil {
		return nil, err
	}
	return &TemplateInfo{
		Type:    templateDexorder,
		Address: add,
		RawHex:  raw,
	}, nil
}
