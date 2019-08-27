package internal

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
)

// CreateMultiSig handles an createmultisig request by returning a
// multisig address for the given inputs.
// This method is originally copied and modified from https://github.com/btcsuite/btcwallet/blob/master/rpc/legacyrpc/methods.go
func CreateMultiSig(cmd *btcjson.CreateMultisigCmd, chainParam *chaincfg.Params) (*btcjson.CreateMultiSigResult, error) {
	script, err := makeMultiSigScript(chainParam, cmd.Keys, cmd.NRequired)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse scirpt %v", err)
	}

	address, err := btcutil.NewAddressScriptHash(script, chainParam)
	if err != nil {
		// above is a valid script, shouldn't happen.
		return nil, err
	}

	return &btcjson.CreateMultiSigResult{
		Address:      address.EncodeAddress(),
		RedeemScript: hex.EncodeToString(script),
	}, nil
}

// makeMultiSigScript is a helper function to combine common logic for
// AddMultiSig and CreateMultiSig.
func makeMultiSigScript(chainParam *chaincfg.Params, keys []string, nRequired int) ([]byte, error) {
	keysesPrecious := make([]*btcutil.AddressPubKey, len(keys))

	// The address list will made up either of addreseses (pubkey hash), for
	// which we need to look up the keys in wallet, straight pubkeys, or a
	// mixture of the two.
	for i, ke := range keys {
		// try to parse as pubkey address
		a, err := decodeAddress(ke, chainParam)
		if err != nil {
			return nil, err
		}

		switch addr := a.(type) {
		case *btcutil.AddressPubKey:
			keysesPrecious[i] = addr
		// default:
		// 	var pubKey btcec.PublicKey
		// 	pubKey, err = w.PubKeyForAddress(addr)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	pubKeyAddr, err := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), chainParam)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	keysesPrecious[i] = pubKeyAddr
		default:
			return nil, fmt.Errorf("Pub key should be 130/66 hex chars")
		}
	}

	return txscript.MultiSigScript(keysesPrecious, nRequired)
}

func decodeAddress(s string, params *chaincfg.Params) (btcutil.Address, error) {
	addr, err := btcutil.DecodeAddress(s, params)
	if err != nil {
		msg := fmt.Sprintf("Invalid address %q: decode failed with %#q", s, err)
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidAddressOrKey,
			Message: msg,
		}
	}
	if !addr.IsForNet(params) {
		msg := fmt.Sprintf("Invalid address %q: not intended for use on %s", addr, params.Name)
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidAddressOrKey,
			Message: msg,
		}
	}
	return addr, nil
}
