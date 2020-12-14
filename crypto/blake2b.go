package crypto

import "golang.org/x/crypto/blake2b"

// Blake2b256 blake2b.Sum256
func Blake2b256(msg []byte) []byte {
	b := blake2b.Sum256(msg)
	return b[:]
}

// Blake2b512 blake2b.Sum512
func Blake2b512(msg []byte) []byte {
	b := blake2b.Sum512(msg)
	return b[:]
}

// Blake2b384 blake2b.Sum384
func Blake2b384(msg []byte) []byte {
	b := blake2b.Sum384(msg)
	return b[:]
}
