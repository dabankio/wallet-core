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

// GenerateScriptPubKey4P2SHP2WPKH https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki#p2wpkh-nested-in-bip16-p2sh
// 相当于listunspent时获取的隔离见证地址的scriptPubk, redeemScript
func GenerateScriptPubKey4P2SHP2WPKH(pubKey *btcec.PublicKey) (redeemScript string, scriptPubk string) {
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
	script := []byte{txscript.OP_0, txscript.OP_DATA_20} //OP_0: 0 -> version
	script = append(script, pubKeyHash...)

	scriptHash := btcutil.Hash160(script)
	scriptPubKeyBytes := append([]byte{txscript.OP_HASH160, txscript.OP_DATA_20}, scriptHash...)
	scriptPubKeyBytes = append(scriptPubKeyBytes, txscript.OP_EQUAL)
	return hex.EncodeToString(script), hex.EncodeToString(scriptPubKeyBytes)
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
