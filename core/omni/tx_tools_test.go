package omni

import (
	"testing"
)

func TestUtilCreatePayloadSimpleSend(t *testing.T) {
	type args struct {
		propertyID uint
		amount     float64
		divisible  bool
	}
	tests := []struct {
		name        string
		args        args
		wantPayload string
	}{
		{
			name:        "不可再分币",
			args:        args{propertyID: 1, amount: 100000000, divisible: false},
			wantPayload: "00000000000000010000000005f5e100",
		},
		{
			name:        "可再分币",
			args:        args{propertyID: 1, amount: 1, divisible: true},
			wantPayload: "00000000000000010000000005f5e100",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPayload, _ := UtilCreatePayloadSimpleSend(tt.args.propertyID, tt.args.amount, tt.args.divisible); gotPayload != tt.wantPayload {
				t.Errorf("UtilCreatePayloadSimpleSend() = %v, want %v", gotPayload, tt.wantPayload)
			}
		})
	}
}

func TestCreaterawtxOpreturn(t *testing.T) {
	type args struct {
		rawtx   string
		payload string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			args: args{
				rawtx:   "01000000013ee02493f5e3fe0eab42118fd365caa2dcefbdc3b14e6787d2f5402c60e5f50e0000000000ffffffff0000000000",
				payload: "000000008000000300000000000000e9",
			},
			want: "01000000013ee02493f5e3fe0eab42118fd365caa2dcefbdc3b14e6787d2f5402c60e5f50e0000000000ffffffff010000000000000000166a146f6d6e69000000008000000300000000000000e900000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreaterawtxOpreturn(tt.args.rawtx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreaterawtxOpreturn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreaterawtxOpreturn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreaterawtxReference(t *testing.T) {
	type args struct {
		rawtx       string
		destination string
		amount      *float64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "regtest case",
			args: args{
				rawtx:       "01000000013ee02493f5e3fe0eab42118fd365caa2dcefbdc3b14e6787d2f5402c60e5f50e0000000000ffffffff010000000000000000166a146f6d6e69000000008000000300000000000000e900000000",
				destination: "mo3eH2xCGMyjpzj9GcZJtxanExVHdEsiUE",
				amount:      nil,
			},
			want: "01000000013ee02493f5e3fe0eab42118fd365caa2dcefbdc3b14e6787d2f5402c60e5f50e0000000000ffffffff020000000000000000166a146f6d6e69000000008000000300000000000000e922020000000000001976a9145296bc27f62ae8c50d30b0a94cc3559c5f4f882d88ac00000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreaterawtxReference(tt.args.rawtx, tt.args.destination, tt.args.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreaterawtxReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreaterawtxReference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreaterawtxChange(t *testing.T) {
	type args struct {
		rawtx       string
		prevtxs     []PreviousDependentTxOutputAmount
		destination string
		fee         float64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "regtest case",
			args: args{
				rawtx:       "01000000010dc7cc05016641d57ee73f2e4ce2b42f497199293fa02113667adf2cb1d5bb9f0100000000ffffffff020000000000000000166a146f6d6e69000000008000000400000000000000e900e1f505000000001976a914e3eb6fcf349a1afdd78f1fe37a41c2aba467461688ac00000000",
				prevtxs:     []PreviousDependentTxOutputAmount{{TxID: "9fbbd5b12cdf7a661321a03f299971492fb4e24c2e3fe77ed541660105ccc70d", Vout: 1, Amount: 12.4999486}},
				destination: "mjzqMzh1YDA936aMLdcjvoo9a9Y3io28Gg",
				fee:         0.0006,
			},
			want: "01000000010dc7cc05016641d57ee73f2e4ce2b42f497199293fa02113667adf2cb1d5bb9f0100000000ffffffff030000000000000000166a146f6d6e69000000008000000400000000000000e900e1f505000000001976a914e3eb6fcf349a1afdd78f1fe37a41c2aba467461688ac0c9d8a44000000001976a91431265d0d6fe3989194dda320c521c0b5e5b1450d88ac00000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreaterawtxChange(tt.args.rawtx, tt.args.prevtxs, tt.args.destination, tt.args.fee)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreaterawtxChange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreaterawtxChange() = %v, want %v", got, tt.want)
			}
		})
	}
}
