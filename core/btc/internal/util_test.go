package internal

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/require"
)

func TestConvertPubk2segWitP2WSHAddress(t *testing.T) {
	for _, tt := range []struct {
		name               string
		privkWIF, wantAddr string
		chain              *chaincfg.Params
	}{
		{
			name:     "testnet", //got from bitcoind
			privkWIF: "cRn7drXcdbdEcMeoXcga5fyyggVycFiWs5zcS2rZJUAv1thjc13i",
			wantAddr: "2My3VgSGfdDFvm7DBnyjBMSff46NnzJooHr",
			chain:    &chaincfg.RegressionNetParams,
		},
		{
			name:     "mainnet", //got from bitcoind
			privkWIF: "L4U2wJcUVmozJHUMRJF49SvGnYDxXKAwpAkqnSa6jGs16P5iWx3u",
			wantAddr: "3FVyXet1jeoZ5Q7buSvM1ZVGfDgws3tngN",
			chain:    &chaincfg.MainNetParams,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			privk, err := btcutil.DecodeWIF(tt.privkWIF)
			require.NoError(t, err)
			add, err := ConvertPubk2segWitP2WSHAddress(privk.PrivKey.PubKey(), tt.chain)
			require.NoError(t, err)
			require.Equal(t, tt.wantAddr, add)
		})
	}

}
