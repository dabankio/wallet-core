package eth

import (
	"crypto/ecdsa"
	"encoding/hex"
	"reflect"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/lomocoin/HDWallet-Core/bip44"
	"github.com/lomocoin/HDWallet-Core/core"
	"github.com/pkg/errors"
)

const (
	symbol = "ETH"
)

type eth struct {
	core.CoinInfo
}

func New(seed []byte) (c *eth, err error) {
	c = new(eth)
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

func NewTransactionFromHex(hex string) (transaction *types.Transaction, err error) {
	data, err := hexutil.Decode(hex)
	if err != nil {
		return
	}
	transaction = new(types.Transaction)
	err = rlp.DecodeBytes(data, transaction)
	if err != nil {
		return
	}
	return
}

func (c *eth) DeriveAddress() (address string, err error) {
	publicKeyECDSAStr, err := c.DerivePublicKey()
	if err != nil {
		return
	}

	publicKeyECDSABytes, err := hex.DecodeString(publicKeyECDSAStr)
	if err != nil {
		return
	}

	publicKeyECDSA, err := crypto.UnmarshalPubkey(publicKeyECDSABytes)
	if err != nil {
		return
	}

	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return
}

func (c *eth) DerivePublicKey() (publicKey string, err error) {
	privateKeyECDSA, err := c.DerivePrivateKey()
	if err != nil {
		return
	}

	privateKey, err := crypto.HexToECDSA(privateKeyECDSA)
	if err != nil {
		return
	}

	public := privateKey.Public()
	publicKeyECDSA, ok := public.(*ecdsa.PublicKey)
	if !ok {
		err = errors.New("failed to get public key")
		return
	}
	publicKey = hex.EncodeToString(crypto.FromECDSAPub(publicKeyECDSA))
	return
}

func (c *eth) DerivePrivateKey() (privateKey string, err error) {
	childKey := c.MasterKey
	for _, childNum := range c.DerivationPath {
		childKey, err = childKey.Child(childNum)
		if err != nil {
			return "", err
		}
	}

	ECPrivateKey, err := childKey.ECPrivKey()
	if err != nil {
		return
	}
	privateKey = hex.EncodeToString(crypto.FromECDSA(ECPrivateKey.ToECDSA()))
	return
}

func (c *eth) DecodeTx(msg string) (tx string, err error) {
	msgType, err := c.checkMsgType(msg)
	if err != nil {
		return
	}

	var data []byte
	if typeMsg := reflect.TypeOf(msgType); typeMsg == reflect.TypeOf(new(types.Transaction)) {
		tr := msgType.(*types.Transaction)
		data, err = tr.MarshalJSON()
	} else if typeMsg == reflect.TypeOf(new(operatingMessage)) {
		op := msgType.(*operatingMessage)
		data, err = op.MarshalJSON()
	} else {
		err = errors.Errorf("not the expected type: %s", typeMsg)
	}
	if err != nil {
		return
	}
	tx = string(data)
	return
}

func (c *eth) checkMsgType(msg string) (msgType interface{}, err error) {
	msgType, err = NewOperatingMessageFromHex(msg)
	if err == nil {
		return
	}
	msgType, err = NewTransactionFromHex(msg)
	if err != nil {
		err = errors.New("decode message err")
		return
	}
	return
}

func (c *eth) nativeSign(transaction *types.Transaction, privateKey *ecdsa.PrivateKey) (*types.Transaction, error) {
	sigBytes, err := crypto.Sign(types.HomesteadSigner{}.Hash(transaction).Bytes(), privateKey)
	if err != nil {
		return nil, err
	}
	return transaction.WithSignature(types.HomesteadSigner{}, sigBytes)
}

// Sign
// hexPrivateKey can't start with "0x"
func (c *eth) Sign(msg, privateKey string) (sig string, err error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return
	}

	msgType, err := c.checkMsgType(msg)
	if err != nil {
		return
	}

	// native sign
	var sigBytes []byte
	if typeMsg := reflect.TypeOf(msgType); typeMsg == reflect.TypeOf(new(types.Transaction)) {
		tr := msgType.(*types.Transaction)
		tr, err = c.nativeSign(tr, privateKeyECDSA)
		if err != nil {
			err = errors.Wrapf(err, "sign in native failed, msg: %s", msg)
		}
		if rlpData, rlpErr := rlp.EncodeToBytes(tr); rlpErr != nil {
			err = rlpErr
		} else {
			sig = hexutil.Encode(rlpData)
		}
		// contract way sign
	} else if typeMsg == reflect.TypeOf(new(operatingMessage)) {
		op := msgType.(*operatingMessage)
		sigBytes, err = op.Sign(privateKeyECDSA)
		if err != nil {
			err = errors.Wrapf(err, "sign in contract way failed, msg: %s", msg)
		} else {
			sig = hexutil.Encode(sigBytes)
		}
	} else {
		err = errors.Errorf("not the expected type: %s", typeMsg)
	}
	return
}

// pubkey is eth address
func (c *eth) VerifySignature(address, msg, signature string) (err error) {
	if !common.IsHexAddress(address) {
		err = errors.Errorf("Incorrect public key (address) : %s", address)
	}
	msgType, err := c.checkMsgType(msg)
	if err != nil {
		return
	}

	var fromAddress common.Address
	if typeMsg := reflect.TypeOf(msgType); typeMsg == reflect.TypeOf(new(types.Transaction)) {
		var signer types.Signer = types.HomesteadSigner{}
		tr := msgType.(*types.Transaction)
		fromAddress, err = types.Sender(signer, tr)
	} else if typeMsg == reflect.TypeOf(new(operatingMessage)) {
		var sign []byte
		var pubKey *ecdsa.PublicKey
		op := msgType.(*operatingMessage)
		sign, err = hexutil.Decode(signature)
		if err != nil {
			return
		}
		pubKey, err = crypto.SigToPub(op.Hash().Bytes(), sign)
		if err != nil {
			return
		}
		fromAddress = crypto.PubkeyToAddress(*pubKey)
	} else {
		err = errors.Errorf("not the expected type: %s", typeMsg)
	}
	if err != nil {
		return
	} else if fromAddress.Hash() != common.HexToAddress(address).Hash() {
		err = errors.Errorf("signature information and publicKey(address) do not match, 1): %s, 2): %s", address, fromAddress.Hex())
		return
	}
	return
}
