package trx

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/core"
	"github.com/ethereum/go-ethereum/crypto"
	trxProtoCore "github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

const symbol = "TRX"

// trx key derivation service
type trx struct{ core.CoinInfo }

func NewCoin(bip44Path string, seed []byte) (core.Coin, error) {
	var err error
	c := new(trx)
	c.Symbol = symbol
	c.DerivationPath, err = bip44.GetCoinDerivationPath(bip44Path, symbol)
	if err != nil {
		return nil, errors.Wrap(err, "bip44.GetCoinDerivationPath err:")
	}
	c.MasterKey, err = hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	return c, err
}

func (c *trx) DeriveAddress() (address string, err error) {
	publicKeyECDSA, err := c.derivePublicKey()
	if err != nil {
		return
	}
	const addressPrefix = 0x41
	return base58.CheckEncode(crypto.PubkeyToAddress(*publicKeyECDSA).Bytes(), addressPrefix), nil
}

func (c *trx) derivePublicKey() (publicKey *ecdsa.PublicKey, err error) {
	privateKeyECDSA, err := c.derivePrivateKey()
	if err != nil {
		return
	}
	publicKey, ok := privateKeyECDSA.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to get public key")
	}
	return
}

func (c *trx) DerivePublicKey() (publicKey string, err error) {
	publicKeyECDSA, err := c.derivePublicKey()
	if err != nil {
		return
	}
	return hex.EncodeToString(crypto.FromECDSAPub(publicKeyECDSA)), nil
}

func (c *trx) derivePrivateKey() (privateKey *ecdsa.PrivateKey, err error) {
	childKey := c.MasterKey
	for _, childNum := range c.DerivationPath {
		childKey, err = childKey.Child(childNum)
		if err != nil {
			return nil, err
		}
	}
	priKey, err := childKey.ECPrivKey()
	if err != nil {
		return
	}
	return priKey.ToECDSA(), nil
}

func (c *trx) DerivePrivateKey() (privateKey string, err error) {
	priKey, err := c.derivePrivateKey()
	if err != nil {
		return
	}
	return hex.EncodeToString(crypto.FromECDSA(priKey)), nil
}

func (c *trx) DecodeTx(msg string) (tx string, err error) {
	return "", errors.New("unsupported")
}

// SignWithPrivateKey .
func SignWithPrivateKey(msg, privateKey string) (string, error) {
	txRawBytes, err := hex.DecodeString(msg)
	if err != nil {
		return "", err
	}
	tx := new(trxProtoCore.Transaction)
	err = proto.Unmarshal(txRawBytes, tx)
	if err != nil {
		return "", errors.Wrap(err, "proto.Unmarshal error")
	}

	rawData, err := proto.Marshal(tx.GetRawData())
	if err != nil {
		return "", errors.Wrap(err, "proto.Marsha err")
	}

	privateKey = strings.TrimPrefix(privateKey, "0x")
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}

	txHash := sha256.Sum256(rawData)
	signature, err := crypto.Sign(txHash[:], privateKeyECDSA)
	if err != nil {
		return "", err
	}
	tx.Signature = append(tx.Signature, signature)
	txSigBytes, err := proto.Marshal(tx)
	if err != nil {
		return "", errors.Wrap(err, "proto.Marshal err")
	}
	return hex.EncodeToString(txSigBytes), nil
}

// Sign msg = trx.core.Transaction
func (c *trx) Sign(msg, privateKey string) (sig string, err error) {
	return SignWithPrivateKey(msg, privateKey)
}

func (c *trx) VerifySignature(pubKey, msg, signature string) error {
	return core.ErrThisFeatureIsNotSupported
}
