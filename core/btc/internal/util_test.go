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

func TestGenerateScriptPubKey4P2SHP2WPKH(t *testing.T) {
	type args struct {
		privkWIF string
	}
	tests := []struct {
		name             string
		args             args
		wantScriptPubk   string
		wantRedeemScript string
	}{
		{//data got from bitcoin node rpc (listunspent)
			name:             "should success", //Desc:"sh(wpkh([18150821]020ef0326de554f936189e59f9a68fcc6bce2e3c2b22a01e09098ffaf356e7749a))#5enqgz2r"
			args:             args{"cSUZDTAa3eBGKeB1CZEZFGESkNns5bdw34wFByS65sgXhGu2PW1R"},
			wantRedeemScript: "001418150821b446b8ecc8af607d2c237c62e0951ada",
			wantScriptPubk:   "a914a66488c1c94664e2a409b0a28719d95fb1a3a3c387",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wif, err := btcutil.DecodeWIF(tt.args.privkWIF)
			require.NoError(t, err)
			gotRedeemScript, gotScriptPubk := GenerateScriptPubKey4P2SHP2WPKH(wif.PrivKey.PubKey())
			if gotScriptPubk != tt.wantScriptPubk {
				t.Errorf("GenerateScriptPubKey4P2SHP2WPKH() gotScriptPubk = %v, want %v", gotScriptPubk, tt.wantScriptPubk)
			}
			if gotRedeemScript != tt.wantRedeemScript {
				t.Errorf("GenerateScriptPubKey4P2SHP2WPKH() gotRedeemScript = %v, want %v", gotRedeemScript, tt.wantRedeemScript)
			}
		})
	}
}
