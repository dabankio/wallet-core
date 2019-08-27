package btcd

import (
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/lomocoin/wallet-core/core/omni"
)

// CreateTransactionForOmni 创建基于btc的omni交易,该方法构建比特币交易输出，包括：
// 1. omni layer Class c opreturn data (`propertyID`,`omniAmount` here)
// 2. dust (amount 546 satoshis) output to `sendToAddress`
// 3. change output to `changeAddress`
//
// `propertyID`, `propertyDivisible` 资产id,token是否可分
// `btcUnspentList` bitcoin utxo list,
// `sendToAddress` omni token收款方
// 
// ref: (usdt: https://omniexplorer.info/asset/31)
// **Note**：不要把找零地址和转账地址传错了，找零地址通常是发送方地址
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
	//omni class c opreturn data
	opreturnScript, err := omni.GetClassCOpreturnDataScript(uint(propertyID), omniAmount, propertyDivisible)
	if err != nil {
		return nil, err
	}
	opreturnTxOut := wire.NewTxOut(0, opreturnScript)

	//dust output
	outAmounts := new(BTCOutputAmount)
	dustAmount, _ := NewBTCAmount(btcutil.Amount(omni.MinNondustOutput).ToBTC())
	outAmounts.Add(sendToAddress, dustAmount)
	return internalNewBTCTransaction(btcUnspentList, outAmounts, changeAddress, btcFeeRate, testNet, []*wire.TxOut{opreturnTxOut})
}
