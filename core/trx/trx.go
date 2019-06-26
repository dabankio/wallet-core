package trx

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
	"github.com/lomocoin/wallet-core/bip44"
	"github.com/lomocoin/wallet-core/core"
	trxProto "github.com/lomocoin/wallet-core/core/trx/proto/core"
	"github.com/pkg/errors"
)

const symbol = "TRX"

// trx key derivation service
type trx struct {
	core.CoinInfo
}

func New(seed []byte) (c *trx, err error) {
	c = new(trx)
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

func (c *trx) DeriveAddress() (address string, err error) {
	publicKeyECDSA, err := c.derivePublicKey()
	if err != nil {
		return
	}
	const addressPrefix = 0x41
	address = base58.CheckEncode(crypto.PubkeyToAddress(*publicKeyECDSA).Bytes(), addressPrefix)
	return
}

func (c *trx) derivePublicKey() (publicKey *ecdsa.PublicKey, err error) {
	privateKeyECDSA, err := c.derivePrivateKey()
	if err != nil {
		return
	}

	publicKey, ok := privateKeyECDSA.Public().(*ecdsa.PublicKey)
	if !ok {
		err = errors.New("failed to get public key")
		return
	}
	return

}

func (c *trx) DerivePublicKey() (publicKey string, err error) {
	publicKeyECDSA, err := c.derivePublicKey()
	if err != nil {
		return
	}
	publicKey = hex.EncodeToString(crypto.FromECDSAPub(publicKeyECDSA))
	return
}

func (c *trx) derivePrivateKey() (privateKey *ecdsa.PrivateKey, err error) {
	childKey := c.MasterKey
	for _, childNum := range c.DerivationPath {
		childKey, err = childKey.Child(childNum)
		if err != nil {
			return
		}
	}
	priKey, err := childKey.ECPrivKey()
	if err != nil {
		return
	}
	privateKey = priKey.ToECDSA()
	return
}

func (c *trx) DerivePrivateKey() (privateKey string, err error) {
	priKey, err := c.derivePrivateKey()
	if err != nil {
		return
	}

	privateKey = hex.EncodeToString(crypto.FromECDSA(priKey))
	return
}

func (c *trx) decodeTx(msg string) (tx string, err error) {
	return
}

func (c *trx) DecodeTx(msg string) (tx string, err error) {
	return
}

/*
msg = trx.core.Transaction
*/
func (c *trx) Sign(msg, privateKey string) (sig string, err error) {
	trRawBytes, err := hex.DecodeString(msg)
	if err != nil {
		return
	}
	tr := new(trxProto.Transaction)
	err = tr.XXX_Unmarshal(trRawBytes)
	if err != nil {
		err = errors.Wrap(err, "XXX_Unmarshal error")
		return
	}

	rawData, err := proto.Marshal(tr.GetRawData())
	if err != nil {
		return
	}
	trHash := sha256.Sum256(rawData)

	privateKey = strings.TrimPrefix(privateKey, "0x")
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return
	}

	var signature []byte
	signature, err = crypto.Sign(trHash[:], privateKeyECDSA)
	if err != nil {
		return
	}
	tr.Signature = append(tr.Signature, signature)

	protoBytes, err := proto.Marshal(tr)
	if err != nil {
		return
	}

	sig = hex.EncodeToString(protoBytes)
	return
}

func (c *trx) VerifySignature(pubKey, msg, signature string) error {
	return core.ErrThisFeatureIsNotSupported
}
