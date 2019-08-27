package internal

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
)

/*
cV4nrs2iHooPTayGs3zcUYKW4wyG4gLQrFhUYXZDNswN3CDeRaKN
cQUzSaxub2pHRg418zCwqJWDhGfShn21kqbXq6UoPRFRKpZsdfHn
cQRMakUKWRvxsyWW525NANf6h3LNdwpcJknteiAR3MjuQ5yVpY6C
*/
func TestSignRawTransaction(t *testing.T) {
	createRawHex := "0200000001a798b1665c115e0208fc5a73aca05597703ce2f063e81f98002c6e2d521f862a0000000000ffffffff0240b564000000000017a91457af28534e50bf8f8f8aa022efabd4a74c02a3028780841e00000000001976a9140f4971a4af12a54b9951c56b92dbcfce2870270588ac00000000"
	var signCmd *SignRawTransactionCmd
	var flags = "ALL"
	signCmd = &SignRawTransactionCmd{
		RawTx: createRawHex,
		Inputs: &[]RawTxInput{{
			Txid:         "2a861f522d6e2c00981fe863f0e23c709755a0ac735afc08025e115c66b198a7",
			Vout:         0,
			ScriptPubKey: "a91457af28534e50bf8f8f8aa022efabd4a74c02a30287",
			RedeemScript: "522103a46570124c8d97fd5425134239ffef051a28f2845a12e05fecc66fd699776cf32103b9fee9c62286eefd6052e86bbb00659e5f746b1d63f7c730539efcd3e5c831cb2103ce6b440f1f0f19b918a9cb5b4b36e92ccc52091d1c6df2f9f059a6445c3a312b53ae",
		}},
		PrivKeys: &[]string{"cV4nrs2iHooPTayGs3zcUYKW4wyG4gLQrFhUYXZDNswN3CDeRaKN"},
		// PrivKeys: &[]string{"cQRMakUKWRvxsyWW525NANf6h3LNdwpcJknteiAR3MjuQ5yVpY6C"},
		Flags: &flags,
	}
	msg, err := SignRawTransaction(signCmd, &chaincfg.TestNet3Params)
	assert.NoError(t, err)
	t.Logf("%+v", msg)
	assert.Equal(t, "0200000001a798b1665c115e0208fc5a73aca05597703ce2f063e81f98002c6e2d521f862a00000000b40047304402202abd11f3298e66a08061354ca34941a3343832a6b1d053e54b130608973acd86022044dad04d3ae0882ffe43aad3299ce1e2e308c2f1cefa1f306aefe135e3b4914d014c69522103a46570124c8d97fd5425134239ffef051a28f2845a12e05fecc66fd699776cf32103b9fee9c62286eefd6052e86bbb00659e5f746b1d63f7c730539efcd3e5c831cb2103ce6b440f1f0f19b918a9cb5b4b36e92ccc52091d1c6df2f9f059a6445c3a312b53aeffffffff0240b564000000000017a91457af28534e50bf8f8f8aa022efabd4a74c02a3028780841e00000000001976a9140f4971a4af12a54b9951c56b92dbcfce2870270588ac00000000", msg.Hex)

	signCmd.RawTx = msg.Hex
	signCmd.PrivKeys = &[]string{"cQUzSaxub2pHRg418zCwqJWDhGfShn21kqbXq6UoPRFRKpZsdfHn"}
	msg, err = SignRawTransaction(signCmd, &chaincfg.TestNet3Params)
	assert.NoError(t, err)
	t.Logf("%+v", msg)
	assert.Equal(t, "0200000001a798b1665c115e0208fc5a73aca05597703ce2f063e81f98002c6e2d521f862a00000000fc00473044022003037bd2f3a58d3e2fc6a1f6fd390d2a792a06bff607f468cf7d04dfc1fdc78c022015fd15af9710dc48a2ec83fa33b100ea15b4ea56d8adcc94353bc6dd6da687e00147304402202abd11f3298e66a08061354ca34941a3343832a6b1d053e54b130608973acd86022044dad04d3ae0882ffe43aad3299ce1e2e308c2f1cefa1f306aefe135e3b4914d014c69522103a46570124c8d97fd5425134239ffef051a28f2845a12e05fecc66fd699776cf32103b9fee9c62286eefd6052e86bbb00659e5f746b1d63f7c730539efcd3e5c831cb2103ce6b440f1f0f19b918a9cb5b4b36e92ccc52091d1c6df2f9f059a6445c3a312b53aeffffffff0240b564000000000017a91457af28534e50bf8f8f8aa022efabd4a74c02a3028780841e00000000001976a9140f4971a4af12a54b9951c56b92dbcfce2870270588ac00000000", msg.Hex)

	signCmd.RawTx = createRawHex
	signCmd.PrivKeys = &[]string{"cV4nrs2iHooPTayGs3zcUYKW4wyG4gLQrFhUYXZDNswN3CDeRaKN", "cQRMakUKWRvxsyWW525NANf6h3LNdwpcJknteiAR3MjuQ5yVpY6C"}
	msg2, err := SignRawTransaction(signCmd, &chaincfg.TestNet3Params)
	assert.NoError(t, err)
	assert.Equal(t, true, msg2.Complete)
}
