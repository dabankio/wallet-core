package testtool

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// AddrInfo 私钥、公钥、地址
type AddrInfo struct {
	PrivkHex, PubkHex, Address string
}

// ToECDSAKey .
func (ad *AddrInfo) ToECDSAKey() *ecdsa.PrivateKey {
	k, err := crypto.HexToECDSA(ad.PrivkHex)
	if err != nil {
		panic(err)
	}
	return k
}

// ToAddress .
func (ad *AddrInfo) ToAddress() common.Address {
	return common.HexToAddress(ad.Address)
}

// GenAddr 生成地址
func GenAddr() *AddrInfo {
	key, _ := crypto.GenerateKey()
	pubKHex := hexutil.Encode(crypto.FromECDSAPub(&key.PublicKey))
	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	return &AddrInfo{
		PrivkHex: hexutil.Encode(crypto.FromECDSA(key))[2:],
		PubkHex:  pubKHex[4:],
		Address:  address,
	}
}
