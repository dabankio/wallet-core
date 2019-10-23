package nxt

import (
	"testing"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/dabankio/wallet-core/core"
	"github.com/stretchr/testify/assert"
)

var (
	testMnemonic = "lecture leg select like delay limit spread retire toward west grape bachelor"
	key          *hdkeychain.ExtendedKey
	nxtCoin      *nxt
)

func TestNXT_NEW(t *testing.T) {
	seed, err := core.NewSeedFromMnemonic(testMnemonic)
	assert.NoError(t, err)
	nxtCoin, err = New(seed)
	assert.NoError(t, err)
}

func TestNXT_DerivePrivateKey(t *testing.T) {
	TestNXT_NEW(t)
	privateKey, err := nxtCoin.DerivePrivateKey()
	assert.NoError(t, err)
	t.Log(privateKey)
}

func TestNXT_DerivePublicKey(t *testing.T) {
	TestNXT_NEW(t)
	pubKey, err := nxtCoin.DerivePublicKey()
	assert.NoError(t, err)
	t.Log(pubKey)
}

func TestNXT_DeriveAddress(t *testing.T) {
	TestNXT_NEW(t)
	address, err := nxtCoin.DeriveAddress()
	assert.NoError(t, err)
	t.Log(address)
}

func TestNXT_DecodeTx(t *testing.T) {
	TestNXT_NEW(t)
	tx, err := nxtCoin.DecodeTx("01194f1e4a09a005d2c0f0710f057f21abb52efd2c2cfea1d03e6d073bd7b5e74d3c6c2bebbb893b4add89a5076a2218000000000000000000e1f5050000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e8d81b005a18959dc5525287010159a9c05c8baee9023330cbe7bdadafcdbb94146a305996acdfe04b1403f6e01600000000")
	assert.NoError(t, err)
	t.Log(tx)
}

func TestNXT_GetAccountById(t *testing.T) {
	TestNXT_NEW(t)
	account, err := nxtCoin.GetAccountById("16254860008892909907")
	assert.NoError(t, err)
	assert.Equal(t, account, "NXT-9LCM-3THW-5XVA-GX875")
}

func TestNXT_GetAccountIdByPk(t *testing.T) {
	TestNXT_NEW(t)
	accountid, err := nxtCoin.GetAccountIdByPk("d2c0f0710f057f21abb52efd2c2cfea1d03e6d073bd7b5e74d3c6c2bebbb893b")
	assert.NoError(t, err)
	assert.Equal(t, accountid, "16254860008892909907")
}

func TestNXT_Sign(t *testing.T) {
	TestNXT_NEW(t)
	msg := "01193f204a09a005d2c0f0710f057f21abb52efd2c2cfea1d03e6d073bd7b5e74d3c6c2bebbb893b4add89a5076a2218000000000000000000e1f5050000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e8d81b005a18959dc5525287010159a9c05c8baee9023330cbe7bdadafcdbb94146a305996acdfe04b1403f6e01600000000"
	privateKey := "ahead respond bare seek half special hurry drum someone burden fault shook"
	sign, err := nxtCoin.Sign(msg, privateKey)
	assert.NoError(t, err)
	assert.Equal(t, sign, "01193f204a09a005d2c0f0710f057f21abb52efd2c2cfea1d03e6d073bd7b5e74d3c6c2bebbb893b4add89a5076a2218000000000000000000e1f50500000000000000000000000000000000000000000000000000000000000000000000000025d5c754c22444fd89a38b6d54d9784a452988dc4393a0af89010fb3de20ca01cdf1a4696ea658904f2d3a3a1eaf09d70cf2e105b0775900fbd1aed33bfe0a4a00000000e8d81b005a18959dc5525287010159a9c05c8baee9023330cbe7bdadafcdbb94146a305996acdfe04b1403f6e01600000000")
}

func TestNXT_VerifySignature(t *testing.T) {
	TestNXT_NEW(t)
	pubKey := "9160dfd95d3c240130681a8b9d20a3df9f6ba12e0fc041ee3e91fbe45d5b630c"
	msg := "0010630e4a09a0059160dfd95d3c240130681a8b9d20a3df9f6ba12e0fc041ee3e91fbe45d5b630c53c9c35f0edd94e100e1f5050000000000c2eb0b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000e8d81b005a18959dc552528701206f1f0000020000000000000000000000000000000353c9c35f0edd94e190e69ee6934cd551669584b9d7aa6d5f000000000000000000000000"
	signature := "6cbf4d2064eb779ee1744e32829d886ba5dc2b9cdb2b7b3680c46f34b18d7405adc4606737bb8c278c783942390b3ff84b7bee515050373dc7909d01d2e4ade8"
	err := nxtCoin.VerifySignature(pubKey, msg, signature)
	assert.EqualError(t, err, core.ErrThisFeatureIsNotSupported.Error())
}
