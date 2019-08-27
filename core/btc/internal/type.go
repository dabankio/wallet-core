package internal

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
)

// RawTxInput models the data needed for raw transaction input that is used in
// the SignRawTransactionCmd struct.
type RawTxInput struct {
	Txid         string  `json:"txid"`
	Vout         uint32  `json:"vout"`
	ScriptPubKey string  `json:"scriptPubKey"`
	RedeemScript string  `json:"redeemScript"`
	Amount       float64 `json:"amount,omitempty"` // only for bch
}

// SignRawTransactionCmd defines the signrawtransaction JSON-RPC command.
type SignRawTransactionCmd struct {
	RawTx    string
	Inputs   *[]RawTxInput
	PrivKeys *[]string
	Flags    *string `default:"\"ALL\""`
}

// NewSignRawTransactionCmd returns a new instance which can be used to issue a
// signrawtransaction JSON-RPC command.
//
// The parameters which are pointers indicate they are optional.  Passing nil
// for optional parameters will use the default value.
func NewSignRawTransactionCmd(hexEncodedTx string, inputs *[]RawTxInput, privKeys *[]string, flags *string) *SignRawTransactionCmd {
	return &SignRawTransactionCmd{
		RawTx:    hexEncodedTx,
		Inputs:   inputs,
		PrivKeys: privKeys,
		Flags:    flags,
	}
}

// WalletTx 钱包可读数据
/*
* version : 2
* locktime : 0
* vin : [{"txid":"07d25a5793dd24cd6d1a810d8bb2958c271ed1971d7e1fb823217a1947170fed","output":0,"sequence":4294967295,"address":"38pfGw2jtkRvwJqXYqLtcFbPS7gHmkWfsN"}]
* vout : [{"address":"38pfGw2jtkRvwJqXYqLtcFbPS7gHmkWfsN","amount":0.084},{"address":"1QLGpxXUfJUVfVNDUJsuQ4dxBppgeuGX5R","amount":0.1}]
 */
type WalletTx struct {
	Version  int32  `json:"version"`
	Locktime uint32 `json:"locktime"`
	Vin      []vin  `json:"vin"`
	Vout     []vout `json:"vout"`
}

type vin struct {
	Txid     string `json:"txid"`
	Vout     uint32 `json:"output"`
	Sequence uint32 `json:"sequence"`
	Address  string `json:"address"`
}

type vout struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}

// 服务端发送的消息
// "{hex:"",rawTxInput:[{scriptPubKey:"xxx",redeemScript:"xxx",amount:0}]}"
type CustomHexMsg struct {
	SignRawTransactionCmd
	DecodeTransaction func(cmd *btcjson.DecodeRawTransactionCmd, params *chaincfg.Params) (
		btcjson.TxRawDecodeResult, error) `json:"-,omitempty"`
	SignTransaction func()   `json:"-,omitempty"`
	walletTx        WalletTx `json:"-"` // covert from txRawDecodeResult
}

func (c *CustomHexMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(*c)
}

func (c *CustomHexMsg) UnmarshalJSON(msg string) (err error) {
	msg = strings.TrimPrefix(msg, "0x")
	data, err := hex.DecodeString(msg)
	if err != nil {
		err = errors.Wrap(err, "CustomHexMsg.hex.DecodeString")
		return
	}
	err = json.Unmarshal(data, &c)
	if err != nil {
		err = errors.Wrap(err, "CustomHexMsg.UnmarshalJson")
		return
	}
	return
}

func (c *CustomHexMsg) getAddressesFromScriptPubKey(txId string, chainCfg *chaincfg.Params) (addresses []string, err error) {
	if c.Inputs == nil {
		err = errors.New("c.Inputs is nil")
		return
	}
	for item := range *c.Inputs {
		if (*c.Inputs)[item].Txid == txId {
			var reply btcjson.DecodeScriptResult
			reply, err = DecodeScript(&btcjson.DecodeScriptCmd{HexScript: (*c.Inputs)[item].ScriptPubKey}, chainCfg)
			if err != nil {
				return
			}
			addresses = reply.Addresses
		}
	}
	return
}

func (c *CustomHexMsg) MarshalToWalletTxJSON(chainCfg *chaincfg.Params) (tx string, err error) {
	if c.DecodeTransaction == nil {
		err = errors.New("decodeTransaction func not set")
		return
	}
	result, err := c.DecodeTransaction(&btcjson.DecodeRawTransactionCmd{HexTx: c.RawTx}, chainCfg)
	if err != nil {
		return
	}
	c.walletTx.Version = result.Version
	c.walletTx.Locktime = result.Locktime
	for item := range result.Vin {
		var in vin
		var addresses []string
		addresses, err = c.getAddressesFromScriptPubKey(result.Vin[item].Txid, chainCfg)

		if len(addresses) != 1 {
			err = errors.New("decode addresses len not equal 1")
			return
		}
		in.Txid = result.Vin[item].Txid
		in.Vout = result.Vin[item].Vout
		in.Sequence = result.Vin[item].Sequence
		in.Address = addresses[0]
		c.walletTx.Vin = append(c.walletTx.Vin, in)
	}
	for item := range result.Vout {
		var out vout
		if addresses := result.Vout[item].ScriptPubKey.Addresses; len(addresses) != 1 {
			// op_return + omni
			// 6a146f6d6e69
			continue
		} else {
			out.Address = addresses[0]
		}
		out.Amount = result.Vout[item].Value
		c.walletTx.Vout = append(c.walletTx.Vout, out)
	}

	data, err := json.Marshal(c.walletTx)
	if err != nil {
		return
	}
	tx = string(data)
	return
}
