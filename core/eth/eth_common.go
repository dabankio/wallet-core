package eth

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash struct {
	hash common.Hash
}

// newHashFromHex converts a hex string to a hash value.
func newHashFromHex(hex string) (hash *Hash, err error) {
	h := new(Hash)
	if err = h.SetHex(hex); err != nil {
		return nil, err
	}
	return h, nil
}

// SetBytes sets the specified slice of bytes as the hash value.
func (h *Hash) SetBytes(hash []byte) error {
	if length := len(hash); length != common.HashLength {
		return fmt.Errorf("invalid hash length: %v != %v", length, common.HashLength)
	}
	copy(h.hash[:], hash)
	return nil
}

// GetBytes retrieves the byte representation of the hash.
func (h *Hash) GetBytes() []byte {
	return h.hash[:]
}

// SetHex sets the specified hex string as the hash value.
func (h *Hash) SetHex(hash string) error {
	hash = strings.ToLower(hash)
	if len(hash) >= 2 && hash[:2] == "0x" {
		hash = hash[2:]
	}
	if length := len(hash); length != 2*common.HashLength {
		return fmt.Errorf("invalid hash hex length: %v != %v", length, 2*common.HashLength)
	}
	bin, err := hex.DecodeString(hash)
	if err != nil {
		return err
	}
	copy(h.hash[:], bin)
	return nil
}

// GetHex retrieves the hex string representation of the hash.
func (h *Hash) GetHex() string {
	return h.hash.Hex()
}

// ETHAddress represents the 20 byte address of an Ethereum account.
type ETHAddress struct {
	address common.Address
}

// NewETHAddress 0地址
func NewETHAddress() *ETHAddress {
	return &ETHAddress{}
}

// NewETHAddressFromHex converts a hex string to a address value.
func NewETHAddressFromHex(hex string) (address *ETHAddress, err error) {
	a := new(ETHAddress)
	if err = a.SetHex(hex); err != nil {
		return nil, err
	}
	return a, nil
}

// SetBytes sets the specified slice of bytes as the address value.
func (a *ETHAddress) SetBytes(address []byte) error {
	if length := len(address); length != common.AddressLength {
		return fmt.Errorf("invalid address length: %v != %v", length, common.AddressLength)
	}
	copy(a.address[:], address)
	return nil
}

// GetBytes retrieves the byte representation of the address.
func (a *ETHAddress) GetBytes() []byte {
	return a.address[:]
}

// SetHex sets the specified hex string as the address value.
func (a *ETHAddress) SetHex(address string) error {
	address = strings.ToLower(address)
	if len(address) >= 2 && address[:2] == "0x" {
		address = address[2:]
	}
	if length := len(address); length != 2*common.AddressLength {
		return fmt.Errorf("invalid address hex length: %v != %v", length, 2*common.AddressLength)
	}
	bin, err := hex.DecodeString(address)
	if err != nil {
		return err
	}
	copy(a.address[:], bin)
	return nil
}

// GetHex retrieves the hex string representation of the address.
func (a *ETHAddress) GetHex() string {
	return a.address.Hex()
}
