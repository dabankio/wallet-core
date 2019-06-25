package crypto

import (
	"fmt"

	"github.com/agl/ed25519"
	"github.com/btcsuite/btcd/btcec"
)

func Sign(privateKey, hash, msg []byte) ([]byte, error) {
	switch len(privateKey) {
	case ed25519.PrivateKeySize:
		return signEd25519(privateKey, msg)
	case btcec.PrivKeyBytesLen:
		return signECDSA(privateKey, hash)
	default:
		return nil, fmt.Errorf("unknown private key format")
	}
}

func Verify(publicKey, hash, msg, signature []byte) (bool, error) {
	switch publicKey[0] {
	case 0xED:
		return verifyEd25519(publicKey, signature, msg)
	case 0x02, 0x03:
		return verifyECDSA(publicKey, signature, hash)
	default:
		return false, fmt.Errorf("unknown public key format")
	}
}

func signEd25519(privateKey, msg []byte) ([]byte, error) {
	var p [ed25519.PrivateKeySize]byte
	copy(p[:], privateKey)
	return ed25519.Sign(&p, msg)[:], nil
}

func verifyEd25519(pubKey, signature, msg []byte) (bool, error) {
	var (
		p [ed25519.PublicKeySize]byte
		s [ed25519.SignatureSize]byte
	)
	switch {
	case len(pubKey) != ed25519.PublicKeySize+1:
		return false, fmt.Errorf("wrong public key length: %d", len(pubKey))
	case pubKey[0] != 0xED:
		return false, fmt.Errorf("wrong public format: %d", pubKey[0])
	case len(signature) != ed25519.SignatureSize:
		return false, fmt.Errorf("wrong Signature length: %d", len(signature))
	default:
		copy(p[:], pubKey[1:])
		copy(s[:], signature)
		return ed25519.Verify(&p, msg, &s), nil
	}
}

// Returns DER encoded signature from input hash
func signECDSA(privateKey, hash []byte) ([]byte, error) {
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKey)
	sig, err := priv.Sign(hash)
	if err != nil {
		return nil, err
	}
	return sig.Serialize(), nil
}

// Verifies a hash using DER encoded signature
func verifyECDSA(pubKey, signature, hash []byte) (bool, error) {
	sig, err := btcec.ParseDERSignature(signature, btcec.S256())
	if err != nil {
		return false, err
	}
	pk, err := btcec.ParsePubKey(pubKey, btcec.S256())
	if err != nil {
		return false, nil
	}
	return sig.Verify(hash, pk), nil
}
