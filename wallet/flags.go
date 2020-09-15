package wallet

import "github.com/dabankio/wallet-core/core/btc"

const (
	//FlagBBCUseStandardBip44ID BBC使用标准bip44 id (默认不是标准bip44 id)
	FlagBBCUseStandardBip44ID = "bbc_use_std_bip44_id"
	//FlagMKFUseBBCBip44ID MKF使用BBC的bip44 id (即MKF和BBC共用地址)
	FlagMKFUseBBCBip44ID = "mkf_use_bbc_bip44_id"
	// FlagBTCUseSegWitFormat BTC使用隔离见证地址
	FlagBTCUseSegWitFormat = btc.FlagUseSegWitFormat
)
