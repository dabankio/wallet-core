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

// NewSymbolCoin symbol 支持 兼容BBC的币种(比如MKF)
func NewSymbolCoin(symbol string, path string, bip44Key string, seed []byte) (core.Coin, error) {
	return internal.NewWallet(symbol, seed, path, bip44Key, nil)
}

// NewSymbolBip44Deriver symbol: BBC | MKF 获取bip44推导
func NewSymbolBip44Deriver(symbol string, bip44Path string, bip44Key string, seed []byte) (bip44.Deriver, error) {
	return internal.NewWallet(symbol, seed, bip44Path, bip44Key, nil)
}

// DecodeSymbolTX 解析原始交易（使用JSON RPC createtransaction 创建的交易）,symbol: BBC | MKF
func DecodeSymbolTX(symbol, rawTX string) (string, error) {
	return internal.DecodeSymbolTx(symbol, rawTX)
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

// CalcTxid 计算txid, symbol: BBC|MKF
func CalcTxid(symbol, rawTx string) (string, error) {
	if txData := internal.TryParseTxDataWithTemplate(rawTx); txData != nil {
		rawTx = txData.TxHex
	}

	se := internal.SymbolSerializer(symbol)
	tx, err := gobbc.DecodeRawTransaction(se, rawTx, true)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse tx data")
	}
	return tx.Txid(se)
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
	price int64,
	fee int32,
	recvAddress string,
	validHeight int32,
	matchAddress string,
	dealAddress string,
	timestamp int64,
) (*TemplateInfo, error) {
	add, raw, err := gobbc.CreateTemplateDataDexOrder(gobbc.DexOrderParam{
		SellerAddress: gobbc.Address(sellerAddress),
		Coinpair:      coinpair,
		Price:         price,
		Fee:           fee,
		RecvAddress:   recvAddress,
		ValidHeight:   validHeight,
		MatchAddress:  gobbc.Address(matchAddress),
		DealAddress:   dealAddress,
		Timestamp:     uint32(timestamp),
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
