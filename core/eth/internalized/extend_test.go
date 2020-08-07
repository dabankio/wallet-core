package internalized

import (
	"encoding/hex"

	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestOperatingMessage_MarshalJSON(t *testing.T) {
	msg, err := NewOperatingMessageFromHex("0xeb854554484552948888888721e49496726d4bf1c32876c0b41cd01f88018a59e97211800080845baf2e9401")
	assert.NoError(t, err)
	data, err := msg.MarshalJSON()
	assert.NoError(t, err)
	t.Logf("%s", data)
}

func TestNewOperatingMessageFromHex(t *testing.T) {
	msg, err := NewOperatingMessageFromHex("0xeb854554484552948888888721e49496726d4bf1c32876c0b41cd01f88018a59e97211800080845baf2e9401")
	assert.NoError(t, err)
	assert.Equal(t, "0x8888888721E49496726D4Bf1c32876c0b41CD01f", msg.ToAddress.Hex())
	msg, err = NewOperatingMessageFromHex("0xf83b854552433230948888888721e49496726d4bf1c32876c0b41cd01f844190ab00948fd4ba16082b1c5844161bc8121507b4e4dace8c845baf3eb102")
	assert.NoError(t, err)
	assert.Equal(t, "0x8888888721E49496726D4Bf1c32876c0b41CD01f", msg.ToAddress.Hex())
	assert.Equal(t, ERC20, msg.Prefix)
}

func TestOperatingParameter_EncodeRLP(t *testing.T) {
	msg, _ := NewOperatingMessageFromHex("0xeb854554484552948888888721e49496726d4bf1c32876c0b41cd01f88018a59e97211800080845baf2e9401")
	rlpData, err := msg.EncodeRLP()
	assert.NoError(t, err)
	t.Log(hex.EncodeToString(rlpData))
}

func TestOperatingMessage_Hash(t *testing.T) {
	msg, _ := NewOperatingMessageFromHex("0xeb854554484552948888888721e49496726d4bf1c32876c0b41cd01f88018a59e97211800080845badf73c01")
	assert.Equal(t, "0xa7395e23fee31ad3dcd3ee17b48b947846cd060de5e323294f2624235071a3eb", msg.Hash().Hex())
	t.Log(msg.Hash().Hex())
}

func TestOperatingMessage_Sign(t *testing.T) {
	msg, _ := NewOperatingMessageFromHex("0xeb854554484552948888888721e49496726d4bf1c32876c0b41cd01f88018a59e97211800080845badf73c01")
	key, err := crypto.HexToECDSA("816680718cceecedbf5d04b994e3d46c9be6f208629b0209083d3bc246208fa7")
	assert.NoError(t, err)
	sig, err := msg.Sign(key)
	assert.Equal(t,
		"0x145253ea8d9bfff7b8a3daebba608b1a9b0d995b7c7d143428472daca1cc261220ed015966f5feb488ac740443df4a1b1ac4747e9ddd7e11c432c69c024eb0bc01",
		hexutil.Encode(sig),
	)
}
