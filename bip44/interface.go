package bip44

// Deriver derives address/publicKey/privateKey
type Deriver interface {
	// DeriveAddress derives the account address of the derivation path.
	DeriveAddress() (address string, err error)
	// DerivePublicKey derives the public key of the derivation path.
	DerivePublicKey() (publicKey string, err error)
	// DerivePrivateKey derives the private key of the derivation path.
	DerivePrivateKey() (privateKey string, err error)
}
