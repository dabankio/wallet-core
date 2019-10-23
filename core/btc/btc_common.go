package btc

import (
	"github.com/btcsuite/btcutil"
	"github.com/dabankio/wallet-core/core/btc/internal"
)

// BTCAddress .
type BTCAddress struct {
	address btcutil.Address
}

// NewBTCAddressFromString converts a string to a address value.
func NewBTCAddressFromString(addr string, chainID int) (address *BTCAddress, err error) {
	netParams, err := internal.ChainFlag2ChainParams(chainID)
	if err != nil {
		return nil, err
	}

	address = new(BTCAddress)
	decAddr, err := btcutil.DecodeAddress(addr, netParams)
	if err != nil {
		return
	}
	address.address = decAddr
	return
}

// BTCAmount 对float64 的封装
type BTCAmount struct {
	amount btcutil.Amount
}

// NewBTCAmount 数量in BTC (not in satoshi)
func NewBTCAmount(amount float64) (amt *BTCAmount, err error) {
	amt = new(BTCAmount)
	tempAmt, err := btcutil.NewAmount(amount)
	if err != nil {
		return
	}
	amt.amount = tempAmt
	return
}
