package internal

import (
	"encoding/json"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
)

func TestDecodeRawTransaction(t *testing.T) {
	cmd := &btcjson.DecodeRawTransactionCmd{
		// cteatrawhex
		HexTx: "0100000001eb29f6d3a024756a66ebf68277b5a816307903657e44da18e98d7d5b00411a0c0000000000ffffffff033f3307000000000017a914611ae902f14f4d1c88a0f06bbb9c6b3c1091fdeb870000000000000000166a146f6d6e690000000000000002000000000393870022020000000000001976a9147598fcf86895d79e81bbf86b308b2c010a8f36eb88ac00000000",
	}
	result, err := DecodeRawTransaction(cmd, &chaincfg.TestNet3Params)
	assert.NoError(t, err)
	j, _ := json.MarshalIndent(result, "", "  ")
	t.Logf("%s\n", j)
}

func TestDecodeScript(t *testing.T) {
	cmd := &btcjson.DecodeScriptCmd{
		HexScript: "0100000001eb29f6d3a024756a66ebf68277b5a816307903657e44da18e98d7d5b00411a0c0000000000ffffffff033f3307000000000017a914611ae902f14f4d1c88a0f06bbb9c6b3c1091fdeb870000000000000000166a146f6d6e690000000000000002000000000393870022020000000000001976a9147598fcf86895d79e81bbf86b308b2c010a8f36eb88ac00000000",
	}
	result, err := DecodeScript(cmd, &chaincfg.TestNet3Params)
	assert.NoError(t, err)
	j, _ := json.MarshalIndent(result, "", "  ")
	t.Logf("%s\n", j)
}
