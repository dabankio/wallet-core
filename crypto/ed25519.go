package crypto

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

// Ed25519sign ed25519签名, seed为私钥, msg: 代签数据
//
// (BBC私钥 hex decode 字符串然后反转字节顺序即可,即: HexDecodeThenReverse)
func Ed25519sign(seed, msg []byte) []byte {
	return ed25519.Sign(ed25519.NewKeyFromSeed(seed), msg)
}

// Ed25519verify 验证签名,公钥 明文 密文
//
// (BBC公钥 hex decode 字符串然后反转字节顺序即可,即: HexDecodeThenReverse)
func Ed25519verify(publicKey, message, sig []byte) bool {
	return ed25519.Verify(publicKey, message, sig)
}

// HexDecodeThenReverse hex decode then reverse byte[]
func HexDecodeThenReverse(str string) ([]byte, error) {
	b, err := hex.DecodeString(str)
	if err != nil {
		return nil, fmt.Errorf("hex decode err: %v", err)
	}
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b, nil
}
