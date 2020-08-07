package internalized

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	ETHER = "ETHER"
	ERC20 = "ERC20"
)

// * The operationHash the result of keccak256(prefix, toAddress, value, data, expireTime).
// * For ether transactions, `prefix` is "ETHER".
// * For token transaction, `prefix` is "ERC20" and `data` is the tokenContractAddress.
type operatingMessage struct {
	Prefix     string         `json:"prefix" gencodec:"required"`
	ToAddress  common.Address `json:"to_address" gencodec:"required"`
	Value      *big.Int       `json:"value" gencodec:"required"`
	Data       []byte         `json:"data" gencodec:"required"`
	ExpireTime *big.Int       `json:"expire_time" gencodec:"required"`
	SequenceId *big.Int       `json:"sequence_id" gencodec:"required"`
}

func NewOperatingMessageFromHex(hex string) (opMsg *operatingMessage, err error) {
	data, err := hexutil.Decode(hex)
	if err != nil {
		return
	}
	opMsg = new(operatingMessage)
	err = rlp.DecodeBytes(data, opMsg)
	if err != nil {
		return
	}
	return
}

func NewOperatingMessage(prefix string, toAddress common.Address, value *big.Int, data []byte, expireTime *big.Int, sequenceId *big.Int) *operatingMessage {
	return &operatingMessage{
		Prefix:     prefix,
		ToAddress:  toAddress,
		Value:      value,
		Data:       data,
		ExpireTime: expireTime,
		SequenceId: sequenceId,
	}
}

func (o *operatingMessage) keccak256Hash() (h common.Hash) {
	opMsg := bytes.Join(
		[][]byte{
			[]byte(o.Prefix),
			o.ToAddress.Bytes(),
			math.U256Bytes(o.Value),
			o.Data,
			math.U256Bytes(o.ExpireTime),
			math.U256Bytes(o.SequenceId)},
		nil,
	)
	h = crypto.Keccak256Hash(opMsg)
	return
}

// Hash hashes the RLP encoding of operating parameter.
// It uniquely identifies the operating parameter.
func (o *operatingMessage) Hash() common.Hash {
	return o.keccak256Hash()
}

func (o *operatingMessage) EncodeRLP() (rlpData []byte, err error) {
	rlpData, err = rlp.EncodeToBytes(o)
	return
}

// MarshalJSON encodes the web3 RPC transaction format.
func (o *operatingMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(*o)
}

func (o *operatingMessage) Sign(privateKey *ecdsa.PrivateKey) (sig []byte, err error) {
	hash := o.Hash()
	return crypto.Sign(hash.Bytes(), privateKey)
}

func NewMultiSignTxFromHex(hex string) (*MultiSignTx, error) {
	var multiSignTx MultiSignTx
	bits, err := hexutil.Decode(hex)
	if err != nil {
		return nil, err
	}
	err = rlp.DecodeBytes(bits, &multiSignTx)
	if err != nil {
		return nil, err
	}
	return &multiSignTx, nil
}

type MultiSignTx struct {
	ExpireTime   int64 //time.Time.Unix()值，在该时间之后区块打包会失败
	ChainID      uint64
	Symbol       string
	MultisigAddr string
	ToAddr       string
	Creator      string
	Executor     string
	Nonce        *big.Int
	Value        *big.Int
	GasLimit     *big.Int
	GasPrice     *big.Int
	Data         []byte
}

func (tx *MultiSignTx) Sign(privateKey *ecdsa.PrivateKey) (string, error) {
	signerPrivkHex := hexutil.Encode(crypto.FromECDSA(privateKey))[2:]
	v, r, s, err := SimpleMultiSigExecuteSign(time.Unix(tx.ExpireTime, 0), int64(tx.ChainID), signerPrivkHex, tx.MultisigAddr, tx.ToAddr, tx.Executor, tx.Nonce, tx.Value, tx.GasLimit, tx.Data)
	if err != nil {
		return "", err
	}
	bytes := append(r[:], s[:]...)
	bytes = append(bytes, byte(v))
	return hexutil.Encode(bytes), nil
}
