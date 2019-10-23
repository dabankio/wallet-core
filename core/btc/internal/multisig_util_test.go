package internal

import (
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
)

// TODO 增加异常测试用例，增加edge case测试用例
func TestCreateMultiSig(t *testing.T) {
	t.Parallel()

	type args struct {
		cmd        *btcjson.CreateMultisigCmd
		chainParam *chaincfg.Params
	}

	//验证数据来源于  https://coinb.in/multisig/ ，以及bitcoind生成
	tests := []struct {
		name    string
		args    args
		want    *btcjson.CreateMultiSigResult
		wantErr bool
	}{
		{
			name: "Mainnet正常生成多签地址(未压缩pubkey)",
			args: args{
				cmd: &btcjson.CreateMultisigCmd{
					NRequired: 2,
					Keys: []string{
						"04cf70b115e925df94ddc977560f0a32c963a664234e16eec5ee824e3725477e1739dc77e49273796baa3e532285a5a5699719d84ad3057dc2b840f2f1225ee75d",
						"04d536bc47783172ce6a0db8a46978115e14fa7ebb237a6358e28522b0b614000c0c06369b49be3837e5793359e292736d8d7a12afeffffeef675658570ecfbb21",
						"048381b5c19aab35784d0b8073533c8c4c94bba6949f3efbb864c1838e9b713ad2d9435d1db0de1ec963fb662e646e6fdff0d44796caf095cd81a30f8b98c9737a",
					},
				}, chainParam: &chaincfg.MainNetParams,
			},
			want: &btcjson.CreateMultiSigResult{
				Address:      "3ErQJHeeRQ86c4C8FERFLtf6ynJDpnEeuj",
				RedeemScript: "524104cf70b115e925df94ddc977560f0a32c963a664234e16eec5ee824e3725477e1739dc77e49273796baa3e532285a5a5699719d84ad3057dc2b840f2f1225ee75d4104d536bc47783172ce6a0db8a46978115e14fa7ebb237a6358e28522b0b614000c0c06369b49be3837e5793359e292736d8d7a12afeffffeef675658570ecfbb2141048381b5c19aab35784d0b8073533c8c4c94bba6949f3efbb864c1838e9b713ad2d9435d1db0de1ec963fb662e646e6fdff0d44796caf095cd81a30f8b98c9737a53ae",
			},
		},
		{
			name: "Mainnet正常生成多签地址(压缩pubkey)",
			args: args{
				cmd: &btcjson.CreateMultisigCmd{
					NRequired: 2,
					Keys: []string{
						"032de7b8c0757f3bc7b76e853372a0f81b4e8030db7ea31f5ae93ccdb8d5578c16",
						"032153a7cfa8148fe2f582b55bbeb7cf0cb97e3619e9149f3cde5cc7e4d4d7b08e",
						"033ba9e6a9508b00e3202b3497ed0c84efb3dd20f2f71720ebca28aa1ee3014ae9",
					},
				}, chainParam: &chaincfg.MainNetParams,
			},
			want: &btcjson.CreateMultiSigResult{
				Address:      "3J8A1K8jL1m4MUx2mmyoUDGUPgejbRWLLJ",
				RedeemScript: "5221032de7b8c0757f3bc7b76e853372a0f81b4e8030db7ea31f5ae93ccdb8d5578c1621032153a7cfa8148fe2f582b55bbeb7cf0cb97e3619e9149f3cde5cc7e4d4d7b08e21033ba9e6a9508b00e3202b3497ed0c84efb3dd20f2f71720ebca28aa1ee3014ae953ae",
			},
		},
		{
			name: "Testnet正常生成多签地址(压缩的pubkey)",
			args: args{
				cmd: &btcjson.CreateMultisigCmd{
					NRequired: 2,
					Keys: []string{
						"0344ed331d03e8f75f836b36ee0c6eb985f59e1b1a8338b3bec930b10b227ff961",
						"03948271d95d0aa4117a89f1871634fef88bc8cd242306c1f2bcfcceacb0e7b890",
						"0332197ef06a09e510d4112d00f8887a395e80ee32d2d59f9a87e9637c8a3c60d1",
					},
				}, chainParam: &chaincfg.TestNet3Params,
			},
			want: &btcjson.CreateMultiSigResult{
				Address:      "2N784vMEL4gRTRaim2XEZHp2ay3sXRAjm2j",
				RedeemScript: "52210344ed331d03e8f75f836b36ee0c6eb985f59e1b1a8338b3bec930b10b227ff9612103948271d95d0aa4117a89f1871634fef88bc8cd242306c1f2bcfcceacb0e7b890210332197ef06a09e510d4112d00f8887a395e80ee32d2d59f9a87e9637c8a3c60d153ae",
			},
		},
		{
			name: "Regtest正常生成多签地址(压缩的pubkey)",
			args: args{
				cmd: &btcjson.CreateMultisigCmd{
					NRequired: 2,
					Keys: []string{
						"02eb34360a90f6138cd949eab24d7dd247427831a818d4981fcd1d1b89743098bc",
						"0379c801b428dd6ac5da394986209374084fe7de84fc016fd9a8bc85db53e2d682",
						"02688d7bc6153de9f01200a1f42628455b938ba7fad01c6f09be50a30097e8b68c",
					},
				}, chainParam: &chaincfg.RegressionNetParams,
			},
			want: &btcjson.CreateMultiSigResult{
				Address:      "2NDjLpneNQcdjPCCfTAzR4uadyDQB1pNME2",
				RedeemScript: "522102eb34360a90f6138cd949eab24d7dd247427831a818d4981fcd1d1b89743098bc210379c801b428dd6ac5da394986209374084fe7de84fc016fd9a8bc85db53e2d6822102688d7bc6153de9f01200a1f42628455b938ba7fad01c6f09be50a30097e8b68c53ae",
			},
		},
		{
			name: "Regtest正常生成多签地址(压缩的pubkey)",
			args: args{
				cmd: &btcjson.CreateMultisigCmd{
					NRequired: 2,
					Keys: []string{
						"021813b82c08ed54ad25997e808fef13b313243e45243b9a559ef26a8efdc5d75d",
						"02779080efbc1a2cd3b7bde36a1e448a899654814ea78f36a8de5c6e33c02da0d3",
						"02440e5e5351aa0214eb4c91a875456de24cc4704635c3cd175c06a584256bf5bd",
					},
				}, chainParam: &chaincfg.RegressionNetParams,
			},
			want: &btcjson.CreateMultiSigResult{
				Address:      "2MyMzun7qPzvJmNqGHfmwGnJZNo4U3PbD4L",
				RedeemScript: "5221021813b82c08ed54ad25997e808fef13b313243e45243b9a559ef26a8efdc5d75d2102779080efbc1a2cd3b7bde36a1e448a899654814ea78f36a8de5c6e33c02da0d32102440e5e5351aa0214eb4c91a875456de24cc4704635c3cd175c06a584256bf5bd53ae",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateMultiSig(tt.args.cmd, tt.args.chainParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMultiSig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateMultiSig() = %v, want %v", got, tt.want)
			}
		})
	}
}
