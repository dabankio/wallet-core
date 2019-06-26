package lmc

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/lomocoin/wallet-core/bip44"
	"github.com/lomocoin/wallet-core/core/btc"
	"github.com/pkg/errors"
)

const symbol = "LMC"

type lmc struct {
	btc.BTC
}

func New(seed []byte) (c *lmc, err error) {
	c = new(lmc)
	c.Symbol = symbol
	c.DerivationPath, err = bip44.GetCoinDerivationPath(symbol)
	if err != nil {
		err = errors.Wrap(err, "bip44.GetCoinDerivationPath err:")
		return
	}
	c.ChainCfg = &chaincfg.MainNetParams
	c.ChainCfg.PubKeyHashAddrID = 0x30
	c.ChainCfg.ScriptHashAddrID = 0x7d
	c.ChainCfg.PrivateKeyID = 0xb0
	c.MasterKey, err = hdkeychain.NewMaster(seed, c.ChainCfg)
	if err != nil {
		return
	}
	return
}

// DecodeTx decodes raw tx to human readable format
/*
return:
* version : 2
* locktime : 0
* vin : [{"txid":"07d25a5793dd24cd6d1a810d8bb2958c271ed1971d7e1fb823217a1947170fed","output":0,"sequence":4294967295,"address":"38pfGw2jtkRvwJqXYqLtcFbPS7gHmkWfsN"}]
* vout : [{"address":"38pfGw2jtkRvwJqXYqLtcFbPS7gHmkWfsN","amount":0.084},{"address":"1QLGpxXUfJUVfVNDUJsuQ4dxBppgeuGX5R","amount":0.1}]
*/
func (c *lmc) DecodeTx(msgTx string) (tx string, err error) {
	var msg = new(btc.CustomHexMsg)
	err = msg.UnmarshalJSON(msgTx)
	if err != nil {
		return
	}
	msg.DecodeTransaction = DecodeRawTransaction
	return msg.MarshalToWalletTxJSON(c.ChainCfg)
}

// Sign signs raw tx with wif privateKey
func (c *lmc) Sign(rawTx, privateKeyWif string) (signedRawTx string, err error) {
	msg := new(btc.CustomHexMsg)
	err = msg.UnmarshalJSON(rawTx)
	if err != nil {
		err = errors.Wrap(err, "lmc.sign.unmarshalJson")
		return
	}
	msg.PrivKeys = &[]string{privateKeyWif}
	if msg.Flags == nil {
		var flagALL = "ALL"
		msg.Flags = &flagALL
	}
	signCmd := &btc.SignRawTransactionCmd{
		RawTx:    msg.RawTx,
		Inputs:   msg.Inputs,
		PrivKeys: msg.PrivKeys,
		Flags:    msg.Flags,
	}
	result, err := SignRawTransaction(signCmd, c.ChainCfg)
	if err != nil {
		err = errors.Wrap(err, "lmc.sign.signRawTransaction")
		return
	}
	signedRawTx = result.Hex
	return
}
