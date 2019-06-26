package xrp

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/lomocoin/wallet-core/bip44"
	"github.com/lomocoin/wallet-core/core"
	"github.com/lomocoin/wallet-core/core/xrp/crypto"
	"github.com/lomocoin/wallet-core/core/xrp/data"
	"github.com/pkg/errors"
)

const (
	symbol = "XRP"
)

type xrp struct {
	core.CoinInfo
}

func New(seed []byte) (c *xrp, err error) {
	c = new(xrp)
	c.Symbol = symbol
	c.DerivationPath, err = bip44.GetCoinDerivationPath(symbol)
	if err != nil {
		err = errors.Wrap(err, "bip44.GetCoinDerivationPath err:")
		return
	}
	c.MasterKey, err = hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return
	}
	return
}

func (c *xrp) DeriveAddress() (address string, err error) {
	ECPublicKey, err := c.derivePublicKey()
	if err != nil {
		return
	}
	accountHash, err := crypto.NewAccountId(crypto.Sha256RipeMD160(ECPublicKey.SerializeCompressed()))
	if err != nil {
		return
	}
	address = accountHash.String()
	return
}

func (c *xrp) derivePublicKey() (publicKey *btcec.PublicKey, err error) {
	ECPrivateKey, err := c.derivePrivateKey()
	if err != nil {
		return
	}
	publicKey = ECPrivateKey.PubKey()
	return
}

func (c *xrp) DerivePublicKey() (publicKey string, err error) {
	ECPublicKey, err := c.derivePublicKey()
	if err != nil {
		return
	}

	publicHash, err := crypto.NewAccountPublicKey(ECPublicKey.SerializeCompressed())
	if err != nil {
		return
	}
	publicKey = publicHash.String()
	return
}

func (c *xrp) derivePrivateKey() (privateKey *btcec.PrivateKey, err error) {
	childKey := c.MasterKey
	for _, childNum := range c.DerivationPath {
		childKey, err = childKey.Child(childNum)
		if err != nil {
			return
		}
	}
	privateKey, err = childKey.ECPrivKey()
	if err != nil {
		return
	}
	return
}

func (c *xrp) DerivePrivateKey() (privateKey string, err error) {
	ECPrivateKey, err := c.derivePrivateKey()
	if err != nil {
		return
	}

	key, err := crypto.NewAccountPrivateKey(ECPrivateKey.Serialize())
	if err != nil {
		return
	}
	privateKey = key.String()
	return
}

func (c *xrp) decodeTx(msg string) (tx data.Transaction, err error) {
	txData, err := hex.DecodeString(msg)
	if err != nil {
		err = errors.Wrap(err, "Hex.DecodeString")
		return
	}
	tx, err = data.ReadTransaction(bytes.NewReader(txData))
	if err != nil {
		err = errors.Wrap(err, "ReadTransaction")
		return
	}
	return
}

func (c *xrp) DecodeTx(msg string) (tx string, err error) {
	transaction, err := c.decodeTx(msg)
	if err != nil {
		err = errors.Wrap(err, "DecodeTx")
		return
	}
	jsData, err := json.MarshalIndent(transaction, "", "  ")
	if err != nil {
		err = errors.Wrap(err, "DecodeTx.JsonMarshal")
		return
	}
	tx = string(jsData)
	return
}

// Only works for multi_sign (xrp native method sign_for)
func (c *xrp) Sign(msg, privateKey string) (sig string, err error) {
	tx, err := c.decodeTx(msg)
	if err != nil {
		err = errors.Wrap(err, "Sign.decode")
		return
	}

	key, err := crypto.NewECDSAKeyFromAccountPrivate(privateKey)
	if err != nil {
		err = errors.Wrap(err, "Sign.NewECDSAKeyFromAccountPrivate")
		return
	}

	err = data.SignFor(tx, key, nil)
	if err != nil {
		err = errors.Wrap(err, "Sign.SignFor")
		return
	}

	_, sigData, err := data.Raw(tx)
	if err != nil {
		err = errors.Wrap(err, "Sign.Raw")
		return
	}

	sig = strings.ToUpper(hex.EncodeToString(sigData))
	return
}

func (c *xrp) VerifySignature(pubKey, msg, signature string) error {
	return core.ErrThisFeatureIsNotSupported
}
