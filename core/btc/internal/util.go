package internal

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

// https://gist.github.com/t4sk/e251e6f298b533039f0f276cf6f5fb28
const (
	pay2pubkHashScriptPrefix = "76a914"
	pay2pubkHashScriptSuffix = "88ac"
)

// GenerateScriptPubKey4PayToPubkeyHash pay to public key hash -> scriptPubKey
func GenerateScriptPubKey4PayToPubkeyHash(pubk []byte) string {
	return pay2pubkHashScriptPrefix + hex.EncodeToString(btcutil.Hash160(pubk)) + pay2pubkHashScriptSuffix
}

// ConvertPubk2segWitP2WSHAddress 把公钥转换为隔离见证地址(p2sh-p2wpkh / p2wsh)
func ConvertPubk2segWitP2WSHAddress(pubKey *btcec.PublicKey, chainParams *chaincfg.Params) (string, error) {
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
	redeemScript, err := txscript.NewScriptBuilder().AddData(pubKeyHash).Script() //OP_0 OP_DATA_20 pubkeyHash
	if err != nil {
		return "", errors.Wrap(err, "build script err")
	}
	redeemScript = append([]byte{0}, redeemScript...) //0: version
	wsh, err := btcutil.NewAddressScriptHash(redeemScript, chainParams)
	if err != nil {
		return "", err
	}
	return wsh.EncodeAddress(), nil
}
