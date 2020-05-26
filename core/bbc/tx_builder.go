package bbc

import (
	"github.com/lomocoin/gobbc"
	"github.com/pkg/errors"
)

// NewTxBuilder new 一个transaction builder
func NewTxBuilder() *TxBuilder { return &TxBuilder{gobbc.NewTXBuilder()} }

// TxBuilder transaction builder
type TxBuilder struct{ *gobbc.TXBuilder }

// Build 构造交易,返回hex编码的tx
func (b *TxBuilder) Build() (string, error) {
	rtx, err := b.TXBuilder.Build()
	if err != nil {
		return "", errors.Wrap(err, "build transaction failed")
	}
	return rtx.Encode(false)
}

// SetAnchor 锚定分支id
func (b *TxBuilder) SetAnchor(anchor string) *TxBuilder {
	return &TxBuilder{b.TXBuilder.SetAnchor(anchor)}
}

// SetTimestamp 交易时间戳
func (b *TxBuilder) SetTimestamp(timestamp int) *TxBuilder {
	return &TxBuilder{b.TXBuilder.SetTimestamp(timestamp)}
}

// SetLockUntil lock until
func (b *TxBuilder) SetLockUntil(lockUntil int) *TxBuilder {
	return &TxBuilder{b.TXBuilder.SetLockUntil(lockUntil)}
}

// SetVersion 当前版本 1
func (b *TxBuilder) SetVersion(v int) *TxBuilder { return &TxBuilder{b.TXBuilder.SetVersion(v)} }

// AddInput 参考listunspent,确保输入金额满足amount
func (b *TxBuilder) AddInput(txid string, vout int8) *TxBuilder {
	return &TxBuilder{b.TXBuilder.AddInput(txid, uint8(vout))}
}

// SetAddress 转账地址,目前只支持公钥地址
func (b *TxBuilder) SetAddress(add string) *TxBuilder { return &TxBuilder{b.TXBuilder.SetAddress(add)} }

// SetAmount 转账金额
func (b *TxBuilder) SetAmount(amount float64) *TxBuilder {
	return &TxBuilder{b.TXBuilder.SetAmount(amount)}
}

// SetFee 手续费，目前0.01，如果带data则0.03, 额外需咨询BBC
func (b *TxBuilder) SetFee(fee float64) *TxBuilder { return &TxBuilder{b.TXBuilder.SetFee(fee)} }

// SetData 原始data设置,参考 UtilDataEncoding
func (b *TxBuilder) SetData(data []byte) *TxBuilder { return &TxBuilder{b.TXBuilder.SetData(data)} }

// SetDataWithUUID 指定uuid,timestamp,data
func (b *TxBuilder) SetDataWithUUID(_uuid string, timestamp int64, data string) *TxBuilder {
	return &TxBuilder{b.TXBuilder.SetDataWithUUID(_uuid, timestamp, data)}
}

// SetStringData 自动编码数据,自动生成uuid和时间戳
func (b *TxBuilder) SetStringData(data string) *TxBuilder {
	return &TxBuilder{b.TXBuilder.SetStringData(data)}
}
