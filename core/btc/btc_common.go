package btcd

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

// BTCAddress
type BTCAddress struct {
	address btcutil.Address
}

// NewBTCAddressFromString converts a string to a address value.
func NewBTCAddressFromString(addr string, testNet bool) (address *BTCAddress, err error) {
	netParams := &chaincfg.MainNetParams
	if testNet {
		netParams = &chaincfg.TestNet3Params
	}

	address = new(BTCAddress)
	decAddr, err := btcutil.DecodeAddress(addr, netParams)
	if err != nil {
		return
	}
	address.address = decAddr
	return
}

// BTCAmount
type BTCAmount struct {
	amount btcutil.Amount
}

func NewBTCAmount(amount float64) (amt *BTCAmount, err error) {
	amt = new(BTCAmount)
	tempAmt, err := btcutil.NewAmount(amount)
	if err != nil {
		return
	}
	amt.amount = tempAmt
	return
}
