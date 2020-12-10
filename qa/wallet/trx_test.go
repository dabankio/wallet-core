package wallet

import (
	"math/big"
	"testing"
	"time"

	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/wallet"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"

	qaTrx "github.com/dabankio/wallet-core/qa/trx"
	r "github.com/stretchr/testify/require"
)

const prec = 1e6 //precision

var trxCIServerGRPC string

func init() {
	trxCIServerGRPC = "192.168.50.5:50051" //should from env var
}

func testTRXPubkSign(t *testing.T, w *wallet.Wallet, c ctx) {
	grpc := client.NewGrpcClientWithTimeout(trxCIServerGRPC, 5*time.Second)
	r.NoError(t, grpc.Start())
	defer grpc.Stop()

	baseAmt := int64(1 * prec)
	trcAmt := big.NewInt(3000)
	var trc20Addr string
	const richAddr = "TMJTFYx6oQVKLkn5pMmsegsvgqKUoYnaEB" //默认有币
	t.Run("推导地址并准备资金", func(t *testing.T) {               //ci 种子地址
		mne := "mirror increase slot auto memory bicycle flip latin correct humble private online"

		options := &wallet.WalletOptions{}
		options.Add(wallet.WithPathFormat(bip44.FullPathFormat))
		options.Add(wallet.WithPassword(""))

		_w, err := wallet.BuildWalletFromMnemonic(mne, true, options)
		r.NoError(t, err)
		addr, err := _w.DeriveAddress("TRX")
		r.NoError(t, err)
		r.Equal(t, richAddr, addr)

		ttx, err := grpc.Transfer(addr, c.address, baseAmt)
		r.NoError(t, err)

		msg, _ := qaTrx.BuildTX(t, ttx)
		sig, err := _w.Sign("TRX", msg)
		r.NoError(t, err)
		qaTrx.BroadcastTX(t, sig, grpc)

		t.Run("部署trc20合约和准备资金", func(t *testing.T) {
			acct, err := grpc.GetAccount(richAddr)
			r.NoError(t, err)
			if acct.AccountResource.EnergyUsage < 100 { //确保能量充足
				ext, err := grpc.FreezeBalance(addr, "", core.ResourceCode_ENERGY, 500*prec)
				r.NoError(t, err)
				msg, _ := qaTrx.BuildTX(t, ext)
				sig, err := _w.Sign("TRX", msg)
				r.NoError(t, err)
				qaTrx.BroadcastTX(t, sig, grpc)
			}

			dptx, err := grpc.DeployContract(
				addr, "trc20xxx", nil, qaTrx.TestTRC20ContractBytecode, 1e9, 20, 9e10,
			)
			r.NoError(t, err)
			msg, txid := qaTrx.BuildTX(t, dptx)
			sig, err := _w.Sign("TRX", msg)
			r.NoError(t, err)
			qaTrx.BroadcastTX(t, sig, grpc)

			time.Sleep(4 * time.Second)
			txinfo, err := grpc.GetTransactionInfoByID(txid)
			r.NoError(t, err)

			trc20Addr = address.Address(txinfo.GetContractAddress()).String()
			t.Log("token addr", trc20Addr)

			//给地址准备trc资金
			trcTx, err := grpc.TRC20Send(richAddr, c.address, trc20Addr, trcAmt, 999999)
			r.NoError(t, err)
			msg, txid = qaTrx.BuildTX(t, trcTx)
			sig, err = _w.Sign("TRX", msg)
			r.NoError(t, err)
			qaTrx.BroadcastTX(t, sig, grpc)
		})
	})

	// 准备资金
	t.Run("TRX转账", func(t *testing.T) {
		acct, err := grpc.GetAccount(c.address)
		r.NoError(t, err)
		r.Equal(t, baseAmt, acct.Balance)

		ttx, err := grpc.Transfer(c.address, richAddr, baseAmt-100)
		r.NoError(t, err)

		msg, _ := qaTrx.BuildTX(t, ttx)
		sig, err := w.Sign("TRX", msg)
		r.NoError(t, err)
		qaTrx.BroadcastTX(t, sig, grpc)
		// time.Sleep(time.Second)

		acct2, err := grpc.GetAccount(c.address)
		r.NoError(t, err)
		r.Equal(t, int64(100), acct2.Balance)
	})

	t.Run("TRC20转账", func(t *testing.T) {
		bal, err := grpc.TRC20ContractBalance(c.address, trc20Addr)
		r.NoError(t, err)
		r.Equal(t, trcAmt, bal)
	})
	t.Run("TRC10转账", func(t *testing.T) {
		t.Skip("")
	})
}
