package btc

import (
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"

	"github.com/lomocoin/wallet-core/core/btc/internal"
)

/*
  {
    "txid": "afd9de3c2858a58a6eddef47c6c2920246744950a564291ce3a28d4c91af1ee2",
    "vout": 1,
    "address": "msDp5uuQdtDuUyHrHFyG2CRgFU9BkUPhAi",
    "account": "",
    "scriptPubKey": "76a914806395d3180f74ef570fd0bfdf3548efcc5e655588ac",
    "amount": 0.07000000,
    "confirmations": 37694,
    "spendable": true,
    "solvable": true,
    "safe": true
  },
0200000001e21eaf914c8da2e31c2964a5504974460292c2c647efdd6e8aa558283cded9af0100000000ffffffff01808d5b00000000001976a9142f1246aac6115313689a2f965edb1642176db6c588ac00000000
0100000001e21eaf914c8da2e31c2964a5504974460292c2c647efdd6e8aa558283cded9af0100000000ffffffff01808d5b00000000001976a9142f1246aac6115313689a2f965edb1642176db6c588ac00000000
*/
func TestNewBTCTransaction(t *testing.T) {
	input := new(BTCUnspent)
	input.Add("d67579f1d8a2c45d807a00fe045322c0210a4e15fa32c8ba2aa6eb07326a5ad7", 1, 4.9992, "76a914817db500feded0e25568d5f5357c9bcb31db159488ac", "")

	addr, err := NewBTCAddressFromString("msKe45XX3Sf6bnYM6UXpxbzpf4STqSFDkU", true)
	// addr, err := btcutil.DecodeAddress("mjoqwuqYzCYjHojkeGjAw7M1nPVcA4cyAS", &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("addr", addr.address.String())
	amt, err := NewBTCAmount(0.0)
	// amt, err := btcutil.NewAmount(0.06)
	if err != nil {
		t.Fatal(err)
	}
	aa := new(BTCOutputAmount)
	aa.Add(addr, amt)

	change, err := NewBTCAddressFromString("msKe45XX3Sf6bnYM6UXpxbzpf4STqSFDkU", true)
	if err != nil {
		t.Fatal(err)
	}
	tt, err := NewBTCTransaction(input, aa, change, 2, true)
	if err != nil {
		t.Fatal(err)
	}
	ii, _ := tt.Encode()
	t.Log(ii)
	hh, err := tt.EncodeToSignCmd()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("hahah", hh)
	t.Log(tt.GetFee())

	bb, _ := internal.New(nil, true)
	cc, err := bb.Sign(hh, "cVBp35B945nC4AEHgAdLJQaGewuFJH4PXAgETBxRmmjavJZtQCAB")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cc)

}

// 该测试用于验证多签工作流
func TestWorkflow_Multisig(t *testing.T) {
	// 测试数据从bitcoind获取

	//流程概述：
	// 构造多签地址 > 构造多签交易 > 分别签名 > 签名完成
	type Addr struct {
		Address string
		Privkey string
		Pubkey  string
	}
	var (
		asrt = assert.New(t)

		addrs              []Addr
		a0, a1, a2, a3, a4 Addr
		multisigAddress    *btcjson.CreateMultiSigResult
	)
	{
		addrs = []Addr{
			{Address: "2NBE1izy3ATokdFHUp4CQpwToFA4RUiwjgE", Privkey: "cRZ7digdFSvWwW883utvZ4mv36ZE5W6TtzaxzjDb93Z4Bs2XqzaM", Pubkey: "031c016155fd1cec7e7fcfbbda440de9afaa5c12f2df810b051f3f51bc49f022dc"},
			{Address: "2N2y9WryiJBu6T5jg1QX1eqMBZrQ37Es2Bb", Privkey: "cRRSMtCSJySyonNE4LJ3CHmJZgxtPPzqXGF7C9apePpvaP5srDvH", Pubkey: "023d6e2f2d536f491c818f4490901bbba36d4081e2e33198f6af4a5075cda0579d"},
			{Address: "2NC3uSNfRKkabB17YPRDTP6KjfqGGE1NFA1", Privkey: "cVxBNYJBEbnwVkQZ8vMZENxKD7q1JRjEhkKPd4sF534evHvN1hq1", Pubkey: "02dbbb5449e7569aee58722c074cd7b626e0c7e18ad341315bd632bbcc0467bc21"},
			{Address: "2NCSzpqEjjsFXBHV6dde6aZ57WcLHyqdGja", Privkey: "cQVsaKPr2h1hGbTHYEeWPFCYBJRCXnBXSuL5hMFT9DJTi5XE9ALV", Pubkey: "034511e6edbd0863fa115d6ce62ea94e53764f1126959785b1055770d36fe8f511"},
			{Address: "2NE9NyB4naou217452tA3tU3oipUFQKC9Bs", Privkey: "cTSrkzXWxHtt9CJmWDzjFUvyr1TUaYYb4uT5puCSmm7uodxuHhbZ", Pubkey: "03a59f2a4636bdf6077c32cc35d8480888718af838f964fe4af7037ad4d9ccb976"},
		}
		a0, a1, a2, a3, a4 = addrs[0], addrs[1], addrs[2], addrs[3], addrs[4]
		_ = a0
	}

	{ //a1/a2/a3 生成多签地址
		rs, err := internal.CreateMultiSig(&btcjson.CreateMultisigCmd{
			NRequired: 2,
			Keys:      []string{a1.Pubkey, a2.Pubkey, a3.Pubkey},
		}, &chaincfg.RegressionNetParams)
		asrt.Nil(err, "生成多签地址失败")
		expectedResult := &btcjson.CreateMultiSigResult{
			Address:      "2NDhpvd63GP8Focyqk2abDAPjEyZ6m2dDSA",
			RedeemScript: "5221023d6e2f2d536f491c818f4490901bbba36d4081e2e33198f6af4a5075cda0579d2102dbbb5449e7569aee58722c074cd7b626e0c7e18ad341315bd632bbcc0467bc2121034511e6edbd0863fa115d6ce62ea94e53764f1126959785b1055770d36fe8f51153ae",
		}
		asrt.Equal(expectedResult, rs, "生成的多签地址不符合预期")
		multisigAddress = rs
	}

	var nextSignData string
	{ //构造交易,并且第一个人签名
		// 多签地址上的UTXO
		data := map[string]interface{}{
			"txid":         "80f910495a7bad56b0c0302b74be6b574ad83acf58a567a442ba7b2714dfed5c",
			"vout":         0,
			"address":      "2NDhpvd63GP8Focyqk2abDAPjEyZ6m2dDSA",
			"scriptPubKey": "a914e06a7e98ed71cee4a484e92daf84873805e42af387",
			"amount":       17,
		}
		utxo := BTCUnspent{}
		utxo.Add(data["txid"].(string), int64(data["vout"].(int)), float64(data["amount"].(int)), data["scriptPubKey"].(string), multisigAddress.RedeemScript)

		outputAmount := new(BTCOutputAmount)
		{
			outAddr, err := NewBTCAddressFromString(a4.Address, true)
			asrt.Nil(err)
			amt, err := NewBTCAmount(9)
			asrt.Nil(err)
			outputAmount.Add(outAddr, amt)

			// changeAddr, err := NewBTCAddressFromString(multisigAddress.Address, true)
			// asrt.Nil(err)
			// changeAmt, err := NewBTCAmount(8 - 0.001)
			// asrt.Nil(err)
			// outputAmount.Add(changeAddr, changeAmt)
		}

		change, err := NewBTCAddressFromString(multisigAddress.Address, true)
		asrt.Nil(err)
		tx, err := NewBTCTransaction(&utxo, outputAmount, change, 2, true)
		asrt.Nil(err)
		hh, err := tx.EncodeToSignCmd()
		asrt.Nil(err)

		btcCoin, _ := internal.New(nil, true)
		signedRawHex, err := btcCoin.Sign(hh, a1.Privkey)
		asrt.Nil(err)

		// 下一个人的签名消息
		nextSignData, err = tx.EncodeToSignCmdForNextSigner(signedRawHex)
	}
	{ // 第二个人签名
		btcCoin, err := internal.New(nil, true)
		asrt.NotNil(err) // TODO 这里先不管吧
		signedRawHex, err := btcCoin.Sign(nextSignData, a3.Privkey)
		asrt.Nil(err)

		_ = signedRawHex
	}

}
