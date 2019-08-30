package wallet

import (
	"github.com/lomocoin/wallet-core/core"
)

var _ core.MetadataProvider = &metadataProviderImpl{}

type metadataProviderImpl struct {
	symbol         string
	password       string
	path           string
	testNet        bool
	chainID        int
	seed           []byte
	derivationPath []uint32
}

func (md *metadataProviderImpl) GetPassword() string {
	return md.password
}

func (md *metadataProviderImpl) GetPath() string {
	return md.path
}

func (md *metadataProviderImpl) IsTestNet() bool {
	return md.testNet
}

func (md *metadataProviderImpl) GetChainID() int { return md.chainID }

func (md *metadataProviderImpl) GetSeed() []byte {
	return md.seed
}

func (md *metadataProviderImpl) GetDerivationPath() []uint32 {
	return md.derivationPath
}
