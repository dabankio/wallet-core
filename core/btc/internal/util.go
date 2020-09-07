package internal

import (
	"encoding/hex"

	"github.com/btcsuite/btcutil"
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
