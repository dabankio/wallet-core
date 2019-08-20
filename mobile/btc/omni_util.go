package btcd

import (
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/lomocoin/wallet-core/core/omni"
)

// CreateTransactionForOmni 创建基于btc的omni交易
// (usdt: https://omniexplorer.info/asset/31)
// 特别提醒：不要把找零地址和转账地址传错了，找零地址通常是发送方地址
func CreateTransactionForOmni(
	propertyID int,
	propertyDivisible bool,
	btcUnspentList *BTCUnspent,
	sendToAddress *BTCAddress,
	omniAmount float64,
	changeAddress *BTCAddress,
	btcFeeRate int64,
	testNet bool,
) (btctx *BTCTransaction, err error) {
	opreturnScript, err := omni.GetOpreturnDataScript(uint(propertyID), omniAmount, propertyDivisible)
	if err != nil {
		return nil, err
	}
	opreturnTxOut := wire.NewTxOut(0, opreturnScript)

	outAmounts := new(BTCOutputAmount)
	dustAmount, _ := NewBTCAmount(btcutil.Amount(omni.MinNondustOutput).ToBTC())
	outAmounts.Add(sendToAddress, dustAmount)
	return internalNewBTCTransaction(btcUnspentList, outAmounts, changeAddress, btcFeeRate, testNet, []*wire.TxOut{opreturnTxOut})
}
