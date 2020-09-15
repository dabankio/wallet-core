package internal

import (
	"encoding/hex"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/core"
	"github.com/pkg/errors"
)

const symbol = "BTC"

// BTC key derivation service
type BTC struct {
	core.CoinInfo
	useSegWit bool //是否使用隔离见证地址
}

// New Factory of BTC key derivation service
//
// The order of publicKeys is important.
// using segWit will replace m/44'... ==> m/49'
func New(bip44Path string, isSegWit bool, seed []byte, chainID int) (c *BTC, err error) {
	c = new(BTC)

	c.Symbol = symbol
	c.useSegWit = isSegWit
	if isSegWit {
		bip44Path = strings.Replace(bip44Path, "m/44'", "m/49'", 1)
	}
	c.DerivationPath, err = bip44.GetCoinDerivationPath(bip44Path, c.Symbol)
	if err != nil {
		return nil, errors.Wrap(err, "bip44.GetCoinDerivationPath err")
	}
	c.ChainCfg, err = ChainFlag2ChainParams(chainID)
	if err != nil {
		return nil, err
	}
	c.MasterKey, err = hdkeychain.NewMaster(seed, c.ChainCfg)
	if err != nil {
		err = errors.Wrap(err, "hdkeychain.NewMaster")
		return
	}
	return c, nil
}

// NewFromMetadata .
func NewFromMetadata(metadata core.MetadataProvider) (c *BTC, err error) {
	chainID := ChainMainNet
	if metadata.IsTestNet() {
		chainID = ChainRegtest
	}
	return New(
		metadata.GetPath(),
		metadata.HasFlag(FlagUseSegWitFormat),
		metadata.GetSeed(),
		chainID,
	)
}

// deriveChildKey derives the child key of the derivation path.
func (c *BTC) deriveChildKey() (childKey *hdkeychain.ExtendedKey, err error) {
	childKey = c.MasterKey
	for _, childOpt := range c.DerivationPath {
		childKey, err = childKey.Child(childOpt)
		if err != nil {
			err = errors.Wrapf(err, "childKey.Child for %x", childOpt)
			return
		}
	}
	return
}

// derivePrivateKey derives the private key of the derivation path.
func (c *BTC) derivePrivateKey() (prikey *btcec.PrivateKey, err error) {
	childKey, err := c.deriveChildKey()
	if err != nil {
		err = errors.Wrap(err, "c.deriveChildKey")
		return
	}
	prikey, err = childKey.ECPrivKey()
	if err != nil {
		err = errors.Wrap(err, "childKey.ECPrivKey")
		return
	}

	return
}

// DerivePrivateKey derives the private key of the derivation path, encoded in string with WIF format
func (c *BTC) DerivePrivateKey() (privateKey string, err error) {
	prikey, err := c.derivePrivateKey()
	if err != nil {
		return
	}
	wif, err := btcutil.NewWIF(prikey, c.ChainCfg, true)
	if err != nil {
		return
	}
	privateKey = wif.String()
	return
}

// DerivePublicKey derives the public key of the derivation path.
func (c *BTC) DerivePublicKey() (publicKey string, err error) {
	prikey, err := c.derivePrivateKey()
	if err != nil {
		return
	}
	publicKey = hex.EncodeToString(prikey.PubKey().SerializeCompressed())
	return
}

// DeriveAddress derives the account address of the derivation path.
func (c *BTC) DeriveAddress() (address string, err error) {
	childKey, err := c.deriveChildKey()
	if err != nil {
		err = errors.Wrap(err, "c.deriveChildKey")
		return
	}
	if !c.useSegWit {
		addrP2PKH, er := childKey.Address(c.ChainCfg)
		if er != nil {
			er = errors.Wrap(er, "childKey.Address")
			return
		}
		address = addrP2PKH.String()
	} else {
		pubk, er := childKey.ECPubKey()
		if er != nil {
			return "", errors.Wrap(er, "childKey.ECPubKey err")
		}
		address, er = ConvertPubk2segWitP2WSHAddress(pubk, c.ChainCfg)
		if er != nil {
			return "", errors.Wrap(er, "convert segWit address err")
		}
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
func (c *BTC) DecodeTx(msgTx string) (tx string, err error) {
	var msg = new(CustomHexMsg)
	err = msg.UnmarshalJSON(msgTx)
	if err != nil {
		return
	}
	msg.DecodeTransaction = DecodeRawTransaction
	return msg.MarshalToWalletTxJSON(c.ChainCfg)
}

// Sign signs raw tx with wif privateKey
func (c *BTC) Sign(rawTx, privateKeyWif string) (signedRawTx string, err error) {
	msg := new(CustomHexMsg)
	err = msg.UnmarshalJSON(rawTx)
	if err != nil {
		err = errors.Wrap(err, "btc.sign.unmarshalJson")
		return
	}
	msg.PrivKeys = &[]string{privateKeyWif}
	if msg.Flags == nil {
		var flagALL = "ALL"
		msg.Flags = &flagALL
	}

	if msg.Inputs != nil {
		wif, err := btcutil.DecodeWIF(privateKeyWif)
		if err != nil {
			return "", err
		}
		for i := 0; i < len(*msg.Inputs); i++ {
			if (*msg.Inputs)[i].ScriptPubKey == "" && (*msg.Inputs)[i].RedeemScript == "" { //未提供 ScriptPubKey, RedeemScript 则默认认为是单签,自动填充scriptPubKey
				(*msg.Inputs)[i].ScriptPubKey = GenerateScriptPubKey4PayToPubkeyHash(wif.SerializePubKey())
			}
		}
	}

	signCmd := &SignRawTransactionCmd{
		RawTx:    msg.RawTx,
		Inputs:   msg.Inputs,
		PrivKeys: msg.PrivKeys,
		Flags:    msg.Flags,
	}
	result, err := SignRawTransaction(signCmd, c.ChainCfg)
	if err != nil {
		err = errors.Wrap(err, "btc.sign.signRawTransaction")
		return
	}
	signedRawTx = result.Hex
	return
}

// VerifySignature verifies rawTx's signature is intact
// If signedRawTx is not signed by pubKey, an error will raise.
func (c *BTC) VerifySignature(pubKey, rawTx, signedRawTx string) error {
	// TODO
	return errors.Errorf("%s is not supported signature verify", c.Symbol)
}
