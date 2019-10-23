package main

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/dabankio/wallet-core/bip39"
)

func main() {
	entropy, err := bip39.NewEntropy(128)
	seed1, err := bip39.NewMnemonic(entropy)
	fmt.Println(seed1)
	seed, err := bip39.NewSeedWithErrorChecking("lecture leg select like delay limit spread retire toward west grape bachelor", "dabank2")
	if err != nil {
		fmt.Println(err)
		return
	}
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return
	}
	purpose, err := masterKey.Child(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(purpose.Depth())
	coin_type, err := purpose.Child(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return
	}
	fmt.Println(coin_type.Depth())
	account, err := coin_type.Child(hdkeychain.HardenedKeyStart + 1)
	if err != nil {
		return
	}
	fmt.Println(account.Depth(), account)
	return
}
