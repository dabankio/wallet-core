package trx

import (
	"fmt"
	"log"
	"testing"

	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/core"
	"github.com/stretchr/testify/assert"
)

func init() {
	var err error
	seed, err := core.NewSeedFromMnemonic(testMnemonic, "")
	coin, err = NewCoin(bip44.FullPathFormat, seed)
	if err != nil {
		log.Fatal(err)
	}
}

var (
	testMnemonic = "lecture leg select like delay limit spread retire toward west grape bachelor"
	coin         core.Coin
)

func TestNew(t *testing.T) {
	for i := 0; i < 3; i++ {
		seed, _ := core.NewSeed()
		coin, _ := NewCoin(bip44.FullPathFormat, seed)
		fmt.Println(coin.DerivePrivateKey())
		fmt.Println(coin.DerivePublicKey())
		fmt.Println(coin.DeriveAddress())
	}
}

func TestTrx_DerivePrivateKey(t *testing.T) {
	pkey, err := coin.DerivePrivateKey()
	assert.NoError(t, err)
	t.Log(pkey)
}

func TestTrx_DerivePublicKey(t *testing.T) {
	pubKey, err := coin.DerivePublicKey()
	assert.NoError(t, err)
	t.Log(pubKey)
}

func TestTrx_DeriveAddress(t *testing.T) {
	address, err := coin.DeriveAddress()
	assert.NoError(t, err)
	t.Log(address)
}

func TestTrx_DecodeTx(t *testing.T) {
	msg := "0a85010a02a0b322083a36fcb6a57a787540c8c6f7ce9a2d5a67080112630a2d747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412320a1541748ba7f6ba19952a5d91981fa9f71d1035418c841215415faa814c0f4eb571f521530aa2e1cf80244784d718a08d06709ce5e1a59a2d12410fd4a3be763ce07420651153f977066c432509e48a2e5f6a7754e2c3784cc8084166ccdfbd21fd1b6b22e7a27d1f2ed50ffbda16eaa967fa2e0135cd6a0c2ebc001241ece66734ee6014565c4537332f0d0ba01dd5279823cc4e09fe05c8569dc40e116a49e6dfce521a72fa47ceaa24524808384f50168920ba83b55fc4c9c9f3f7b401"
	coin.DecodeTx(msg)
}

func TestTrx_Sign(t *testing.T) {
	bb, err := coin.Sign("0a85010a02ace62208b3f94ba5e72ef48540aefa95ce992d5a67080112630a2d747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412320a1541b917d6261d3bbe00b9183f9b375da811677f6ea31215415faa814c0f4eb571f521530aa2e1cf80244784d718a08d0670ae9dbacc992d124186effa70360c01cbeac998c88a54d9c001863b4643db3fb1faa2d5f37249c550459481e42d6132c10ef15127704f777c9fb1b73f0d332a60eacd3ef31de4799100", "27be6bc7e7a02f310528bbc2a78f724bac9240fd21e9fade9fbcf82405681b8c")
	assert.NoError(t, err)
	t.Log(bb)
	// 0acb010a85010a02ace62208b3f94ba5e72ef48540aefa95ce992d5a67080112630a2d747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412320a1541b917d6261d3bbe00b9183f9b375da811677f6ea31215415faa814c0f4eb571f521530aa2e1cf80244784d718a08d0670ae9dbacc992d124186effa70360c01cbeac998c88a54d9c001863b4643db3fb1faa2d5f37249c550459481e42d6132c10ef15127704f777c9fb1b73f0d332a60eacd3ef31de4799100
	// 0a85010a02ace62208b3f94ba5e72ef48540aefa95ce992d5a67080112630a2d747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412320a1541b917d6261d3bbe00b9183f9b375da811677f6ea31215415faa814c0f4eb571f521530aa2e1cf80244784d718a08d0670ae9dbacc992d124186effa70360c01cbeac998c88a54d9c001863b4643db3fb1faa2d5f37249c550459481e42d6132c10ef15127704f777c9fb1b73f0d332a60eacd3ef31de4799100124157954b88c61fd11d24563800463f99c6af3d447e6055f99019dee5ed365dda6b42a402a4ecd752950aedb1b6ea3d76d62176a0335009171b8a78ee04d8d3a3b500
}
