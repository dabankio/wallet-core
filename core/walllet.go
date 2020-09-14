package core

import (
	"crypto/rand"

	"github.com/dabankio/wallet-core/bip39"
	"github.com/pkg/errors"
)

// NewEntropy will create random entropy bytes
func NewEntropy(bits int) (entropy []byte, err error) {
	return bip39.NewEntropy(bits)
}

// NewMnemonic returns a randomly generated BIP-39 mnemonic using 128-256 bits of entropy.
func NewMnemonic(entropy []byte) (mnemonic string, err error) {
	return bip39.NewMnemonic(entropy)
}

// NewSeed returns a randomly generated BIP-39 seed.
func NewSeed() (b []byte, err error) {
	b = make([]byte, 64)
	_, err = rand.Read(b)
	return
}

// NewSeedFromMnemonic returns a BIP-39 seed based on a BIP-39 mnemonic.
func NewSeedFromMnemonic(mnemonic, password string) ([]byte, error) {
	if mnemonic == "" {
		return nil, errors.New("mnemonic is required")
	}

	return bip39.NewSeedWithErrorChecking(mnemonic, password)
}

// EntropyFromMnemonic
// returns the input entropy used to generate the given mnemonic
func EntropyFromMnemonic(mnemonic string) (entropy []byte, err error) {
	return bip39.EntropyFromMnemonic(mnemonic)
}
