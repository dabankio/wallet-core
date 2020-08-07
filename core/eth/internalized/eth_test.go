package internalized

import (
	"testing"

	"github.com/dabankio/wallet-core/core"
	"github.com/stretchr/testify/assert"
)

var (
	testMnemonic = "lecture leg select like delay limit spread retire toward west grape bachelor"
	ethCoin      = &eth{}
)

func TestEth_VerifySignature(t *testing.T) {
	TestNewETH(t)
	pkey := "0x9207f1e00b9e6b6fe2f8cffec52f2fb36029cc9df5a73096a59769f24f4d49e6"
	msgeth := "0xf8e90a843b9aca008307d00094e2112d55f0d5b94143dfd2e5bd18dc3a65862bf380b884c6427474000000000000000000000000d783ae7026cf8d8081ae3d5b4ea8b4b82bda3d3900000000000000000000000000000000000000000000000000470de4df820000000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000001ca024dad1adab099b389862aa200b6d4e2b544700b0d8de8e10b131a0c5c2fcc6f2a052683a572a2bbb95bd48404cd024e1b025b399c81e3ba72fc0368d8077a633f8"
	msgcontract := "0xea85455448455294d783ae7026cf8d8081ae3d5b4ea8b4b82bda3d3987470de4df82000080845baf694104"
	sig, err := ethCoin.Sign(msgeth, pkey)
	assert.NoError(t, err)
	t.Log(sig)
	err = ethCoin.VerifySignature("0x947ab281Df5ec46E801F78Ad1363FaaCbe4bfD12", sig, "")
	assert.NoError(t, err)
	sig, err = ethCoin.Sign(msgcontract, pkey)
	assert.NoError(t, err)
	t.Log(sig)
	err = ethCoin.VerifySignature("0x947ab281Df5ec46E801F78Ad1363FaaCbe4bfd12", msgcontract, sig)
	assert.NoError(t, err)
}

func TestEth_DecodeTx(t *testing.T) {
	TestNewETH(t)
	msgeth := "0xf8e90a843b9aca008307d00094e2112d55f0d5b94143dfd2e5bd18dc3a65862bf380b884c6427474000000000000000000000000d783ae7026cf8d8081ae3d5b4ea8b4b82bda3d3900000000000000000000000000000000000000000000000000470de4df820000000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000001ca024dad1adab099b389862aa200b6d4e2b544700b0d8de8e10b131a0c5c2fcc6f2a052683a572a2bbb95bd48404cd024e1b025b399c81e3ba72fc0368d8077a633f8"
	msgcontract := "0xea85455448455294d783ae7026cf8d8081ae3d5b4ea8b4b82bda3d3987470de4df82000080845baf694104"
	data, err := ethCoin.DecodeTx(msgeth)
	assert.NoError(t, err)
	t.Log(data)
	data, err = ethCoin.DecodeTx(msgcontract)
	assert.NoError(t, err)
	t.Log(data)
}

func TestEth_Sign(t *testing.T) {
	TestNewETH(t)
	pkey := "816680718cceecedbf5d04b994e3d46c9be6f208629b0209083d3bc246208fa7"
	msgeth := "0xf8e90a843b9aca008307d00094e2112d55f0d5b94143dfd2e5bd18dc3a65862bf380b884c6427474000000000000000000000000d783ae7026cf8d8081ae3d5b4ea8b4b82bda3d3900000000000000000000000000000000000000000000000000470de4df820000000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000001ca024dad1adab099b389862aa200b6d4e2b544700b0d8de8e10b131a0c5c2fcc6f2a052683a572a2bbb95bd48404cd024e1b025b399c81e3ba72fc0368d8077a633f8"
	msgcontract := "0xea85455448455294d783ae7026cf8d8081ae3d5b4ea8b4b82bda3d3987470de4df82000080845baf694104"
	sig, err := ethCoin.Sign(msgeth, pkey)
	assert.NoError(t, err)
	t.Log(sig)
	sig, err = ethCoin.Sign(msgcontract, pkey)
	assert.NoError(t, err)
	t.Log(sig)
}

func TestEth_Sign_SimpleMultiSign(t *testing.T) {
	pkey := "f0b24f7d0e472855a57d5ff4c8239c297dbcb33e2f94860e8a7cea763c0e4509"
	msg := "0xf89b0183455448aa307835303835426338393832383763333931436232353942363830384542643064623439324633303739aa307845306446313937413539304331383645413342306234663061346334313864613431393039323837aa3078323862376332643342376238614632334134426137393041304165463439656236633345613136328080880de0b6b3a76400008257d28504a817c80080"
	sig, err := ethCoin.Sign(msg, pkey)
	assert.NoError(t, err)
	t.Log(sig)
}

func TestEth_DeriveAddress(t *testing.T) {
	TestNewETH(t)
	address, err := ethCoin.DeriveAddress()
	assert.NoError(t, err)
	t.Logf("%s", address)
}

func TestEth_DerivePublicKey(t *testing.T) {
	TestNewETH(t)
	pkey, err := ethCoin.DerivePublicKey()
	assert.NoError(t, err)
	t.Logf("%s", pkey)
}

func TestEth_DerivePrivateKey(t *testing.T) {
	TestNewETH(t)
	pkey, err := ethCoin.DerivePrivateKey()
	assert.NoError(t, err)
	t.Logf("%s", pkey)
}

func TestNewETH(t *testing.T) {
	seed, err := core.NewSeedFromMnemonic(testMnemonic)
	assert.NoError(t, err)
	ethCoin, err = New(seed)
	assert.NoError(t, err)
}
