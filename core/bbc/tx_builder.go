package bbc

import (
	"strings"

	"github.com/dabankio/gobbc"
	"github.com/pkg/errors"
)

// NewTxBuilder new 一个transaction builder
func NewTxBuilder() *TxBuilder { return &TxBuilder{gobbc.NewTXBuilder(), false, []string{}} }

// TxBuilder transaction builder
type TxBuilder struct {
	*gobbc.TXBuilder
	excludeAnchor bool     //默认需要anchor字段，MKF使用则设为true
	tplData       []string //模版原始数据
}

// Build 构造交易,返回hex编码的tx
func (b *TxBuilder) Build() (string, error) {
	rtx, err := b.TXBuilder.Build()
	if err != nil {
		return "", errors.Wrap(err, "build transaction failed")
	}
	ser := gobbc.BBCSerializer
	if b.excludeAnchor {
		ser = gobbc.MKFSerializer
	}
	raw, err := rtx.Encode(ser, false)
	if err != nil {
		return "", err
	}
	if len(b.tplData) == 0 { //无模版数据
		return raw, nil
	}
	d := gobbc.TXData{
		TplHex: strings.Join(b.tplData, gobbc.TemplateDataSpliter),
		TxHex:  raw,
	}
	return d.EncodeString()
}

// SetAnchor 锚定分支id
func (b *TxBuilder) SetAnchor(anchor string) *TxBuilder {
	b.TXBuilder.SetAnchor(anchor)
	return b
}

// ExcludeAnchor MKF需要调用该函数(使得序列化时不会处理anchor字段)
func (b *TxBuilder) ExcludeAnchor() *TxBuilder {
	b.excludeAnchor = true
	return b
}

// SetTimestamp 交易时间戳
func (b *TxBuilder) SetTimestamp(timestamp int) *TxBuilder {
	b.TXBuilder.SetTimestamp(timestamp)
	return b
}

// SetLockUntil lock until
func (b *TxBuilder) SetLockUntil(lockUntil int) *TxBuilder {
	b.TXBuilder.SetLockUntil(lockUntil)
	return b
}

// SetVersion 当前版本 1
func (b *TxBuilder) SetVersion(v int) *TxBuilder {
	b.TXBuilder.SetVersion(v)
	return b
}

// SetType typ
func (b *TxBuilder) SetType(v int) *TxBuilder {
	b.TXBuilder.SetType(v)
	return b
}

// AddInput 参考listunspent,确保输入金额满足amount
func (b *TxBuilder) AddInput(txid string, vout int8) *TxBuilder {
	b.TXBuilder.AddInput(txid, uint8(vout))
	return b
}

// AddTemplateData 添加模版原始数据,多个模版时需要自行确保顺序
func (b *TxBuilder) AddTemplateData(tplData string) *TxBuilder {
	b.tplData = append(b.tplData, tplData)
	return b
}

// SetAddress 转账地址,目前只支持公钥地址
func (b *TxBuilder) SetAddress(add string) *TxBuilder {
	b.TXBuilder.SetAddress(add)
	return b
}

// SetAmount 转账金额
func (b *TxBuilder) SetAmount(amount float64) *TxBuilder {
	b.TXBuilder.SetAmount(amount)
	return b
}

// SetFee 手续费，目前0.01，如果带data则0.03, 额外需咨询BBC
func (b *TxBuilder) SetFee(fee float64) *TxBuilder {
	b.TXBuilder.SetFee(fee)
	return b
}

// SetData 原始data设置,参考 UtilDataEncoding
func (b *TxBuilder) SetData(data []byte) *TxBuilder {
	b.TXBuilder.SetData(data)
	return b
}

// SetDataWithUUID 指定uuid,timestamp,data
func (b *TxBuilder) SetDataWithUUID(_uuid string, timestamp int64, data string) *TxBuilder {
	b.TXBuilder.SetDataWithUUID(_uuid, timestamp, data)
	return b
}

// SetStringData 自动编码数据,自动生成uuid和时间戳
func (b *TxBuilder) SetStringData(data string) *TxBuilder {
	b.TXBuilder.SetStringData(data)
	return b
}
