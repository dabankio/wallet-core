package btcd

import (
	"testing"

	"github.com/lomocoin/wallet-core/core/btc"
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

	bb, _ := btc.New(nil, true)
	cc, err := bb.Sign(hh, "cVBp35B945nC4AEHgAdLJQaGewuFJH4PXAgETBxRmmjavJZtQCAB")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cc)

}
