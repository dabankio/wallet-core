package data

import "github.com/lomocoin/HDWallet-Core/core/xrp/crypto"

func Sign(s SignerAgent, key crypto.Key, sequence *uint32) error {
	s.InitialiseForSigning()
	copy(s.GetPublicKey().Bytes(), key.Public(sequence))
	hash, msg, err := SigningHash(s, nil)
	if err != nil {
		return err
	}
	sig, err := crypto.Sign(key.Private(sequence), hash.Bytes(), append(s.SigningPrefix().Bytes(), msg...))
	if err != nil {
		return err
	}
	*s.GetSignature() = VariableLength(sig)
	hash, _, err = Raw(s)
	if err != nil {
		return err
	}
	copy(s.GetHash().Bytes(), hash.Bytes())
	return nil
}

func SignFor(s SignerAgent, key crypto.Key, sequence *uint32) error {
	s.InitialiseForMultiSigning()
	hash, msg, err := SigningHash(s, key.Id(sequence))
	if err != nil {
		return err
	}
	sig, err := crypto.Sign(key.Private(sequence), hash.Bytes(), append(s.SigningPrefix().Bytes(), msg...))
	if err != nil {
		return err
	}
	signer := &Signer{
		Signer: struct {
			Account       *Account        `json:",omitempty"`
			SigningPubKey *PublicKey      `json:",omitempty"`
			TxnSignature  *VariableLength `json:",omitempty"`
		}{
			Account:       new(Account),
			SigningPubKey: new(PublicKey),
			TxnSignature:  new(VariableLength),
		},
	}
	copy(signer.Signer.Account[:], key.Id(sequence))
	copy(signer.Signer.SigningPubKey.Bytes(), key.Public(sequence))
	*signer.Signer.TxnSignature = VariableLength(sig)
	s.AddSignature(signer)
	hash, _, err = Raw(s)
	if err != nil {
		return err
	}
	copy(s.GetHash().Bytes(), hash.Bytes())
	return nil
}

func CheckSignature(s SignerAgent) (bool, error) {
	hash, msg, err := SigningHash(s, nil)
	if err != nil {
		return false, err
	}
	return crypto.Verify(s.GetPublicKey().Bytes(), hash.Bytes(), msg, s.GetSignature().Bytes())
}
