package omni

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

// UtilCreatePayloadSimpleSend 对于可分的币会 *1e8 (btcutil.SatoshiPerBitcoin)
func UtilCreatePayloadSimpleSend(propertyID uint, amount float64, divisible bool) (string, error) {
	var intPart int64

	if divisible {
		amt, err := btcutil.NewAmount(amount)
		if err != nil {
			return "", err
		}
		intPart = int64(amt)
	} else {
		intPart = int64(amount)
	}
	return fmt.Sprintf("%016x%016x", propertyID, intPart), nil
}

// GetClassCOpreturnDataScript Create a payload with class C (op-return) encoding to pkScript. Used to create tx out.
func GetClassCOpreturnDataScript(propertyID uint, amount float64, divisible bool) ([]byte, error) {
	payload, err := UtilCreatePayloadSimpleSend(propertyID, amount, divisible)
	if err != nil {
		return nil, err
	}

	b, err := hex.DecodeString(omniHex + payload)
	if err != nil {
		return nil, fmt.Errorf("could not decode payload, %v", err)
	}
	opreturnScript, err := txscript.NullDataScript(b)
	if err != nil {
		return nil, fmt.Errorf("failed to create opreturn data, %v", err)
	}
	return opreturnScript, nil
}

// CreaterawtxOpreturn Impl: https://github.com/OmniLayer/omnicore/blob/master/src/omnicore/doc/rpc-api.md#omni_createrawtx_opreturn
func CreaterawtxOpreturn(rawtx, payload string) (string, error) {
	mtx, err := decodeRawtx2mtx(rawtx)
	if err != nil {
		return rawtx, fmt.Errorf("could not decode rawtx, %v", err)
	}

	b, err := hex.DecodeString(omniHex + payload)
	if err != nil {
		return rawtx, fmt.Errorf("could not decode payload, %v", err)
	}
	opreturnScript, err := txscript.NullDataScript(b)
	if err != nil {
		return rawtx, fmt.Errorf("failed to create opreturn data, %v", err)
	}

	mtx.AddTxOut(wire.NewTxOut(0, opreturnScript))
	return hexEncodeBTCTx(mtx)
}

// CreaterawtxReference amount:btc, Impl: https://github.com/OmniLayer/omnicore/blob/master/src/omnicore/doc/rpc-api.md#omni_createrawtx_reference
func CreaterawtxReference(rawtx, destination string, amount *float64) (string, error) {
	mtx, err := decodeRawtx2mtx(rawtx)
	if err != nil {
		return rawtx, fmt.Errorf("could not decode raw tx, %v", err)
	}

	if amount == nil {
		amount = btcjson.Float64(0)
	}
	outputAmount := *amount
	if int64(*amount*btcutil.SatoshiPerBitcoin) < MinNondustOutput { //关于minOutput的计算参考：https://github.com/OmniLayer/omnicore/blob/master/src/omnicore/createtx.cpp#line:137, (commit hash: 7595200) 此处直接写死
		outputAmount = btcutil.Amount(MinNondustOutput).ToBTC()
	}

	err = appendTxOut(mtx, destination, &outputAmount)
	if err != nil {
		return rawtx, fmt.Errorf("could not append reference output , %v", err)
	}

	return hexEncodeBTCTx(mtx)
}

// PreviousDependentTxOutputAmount .
type PreviousDependentTxOutputAmount struct {
	TxID   string
	Vout   uint32
	Amount float64
}

// CreaterawtxChange Impl: https://github.com/OmniLayer/omnicore/blob/master/src/omnicore/doc/rpc-api.md#omni_createrawtx_change
func CreaterawtxChange(rawtx string, prevtxs []PreviousDependentTxOutputAmount, destination string, fee float64) (string, error) {
	mtx, err := decodeRawtx2mtx(rawtx)
	if err != nil {
		return rawtx, fmt.Errorf("could not decode raw tx, %v", err)
	}

	var totalInputAmountSatoshis btcutil.Amount
	{ // 校验并计算总输入
		//提供的依赖输出没有缺失
		findPreviousDependentTxOutput := func(hash string, vout uint32) *PreviousDependentTxOutputAmount {
			for i := 0; i < len(prevtxs); i++ {
				o := &prevtxs[i]
				if o.TxID == hash && o.Vout == vout {
					return o
				}
			}
			return nil
		}

		for _, in := range mtx.TxIn {
			var pout *PreviousDependentTxOutputAmount
			txHash, vout := in.PreviousOutPoint.Hash.String(), in.PreviousOutPoint.Index
			if pout = findPreviousDependentTxOutput(txHash, vout); pout == nil {
				return rawtx, fmt.Errorf("previous out tx not find, tx: %v, vout: %d", txHash, vout)
			}

			amt, err := btcutil.NewAmount(pout.Amount)
			if err != nil {
				return rawtx, err
			}
			totalInputAmountSatoshis += amt
		}
	}

	{ // 计算找零
		var totalOut int64
		for _, out := range mtx.TxOut {
			totalOut += out.Value
		}
		feeAmt, err := btcutil.NewAmount(fee)
		if err != nil {
			return rawtx, fmt.Errorf("invalid fee %v", fee)
		}
		changeSatoshis := int64(totalInputAmountSatoshis) - totalOut - int64(feeAmt)
		if changeSatoshis < 0 {
			return rawtx, fmt.Errorf("make change failed, insufficient input, total input(satoshis): %d, gived fee(satoshis): %v", totalInputAmountSatoshis, feeAmt)
		}
		changeBTC := btcutil.Amount(changeSatoshis).ToBTC()
		appendTxOut(mtx, destination, &changeBTC)
	}

	return hexEncodeBTCTx(mtx)
}

func decodeRawtx2mtx(rawtx string) (*wire.MsgTx, error) {
	var mtx wire.MsgTx
	hexStr := rawtx
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	serializedTx, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	err = mtx.Deserialize(bytes.NewReader(serializedTx)) //TODO 确定隔离验证支持(mtx.DeserializeNoWitness)
	if err != nil {
		return nil, err
	}
	return &mtx, nil
}

func hexEncodeBTCTx(mtx *wire.MsgTx) (string, error) {
	var buf bytes.Buffer
	if err := mtx.BtcEncode(&buf, wire.ProtocolVersion, wire.WitnessEncoding); err != nil {
		return "", fmt.Errorf("Failed to encode msg of type %T", mtx)
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func appendTxOut(mtx *wire.MsgTx, destination string, amount *float64) error {
	params := &chaincfg.RegressionNetParams //FIXME config able

	if amount == nil {
		amount = btcjson.Float64(0)
	}

	// Ensure amount is in the valid range for monetary amounts.
	if *amount <= 0 || *amount > btcutil.MaxSatoshi {
		return fmt.Errorf("invalid amount %v", *amount)
	}

	// Decode the provided address.
	addr, err := btcutil.DecodeAddress(destination, params)
	if err != nil {
		return fmt.Errorf("Invalid address or key: " + err.Error())
	}

	// Ensure the address is one of the supported types and that
	// the network encoded with the address matches the network the
	// server is currently on.
	switch addr.(type) {
	case *btcutil.AddressPubKeyHash:
	case *btcutil.AddressScriptHash:
	default:
		return fmt.Errorf("Invalid address or key")
	}
	if !addr.IsForNet(params) {
		return fmt.Errorf("Invalid address: " + destination + " is for the wrong network")
	}
	// Create a new script which pays to the provided address.
	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return fmt.Errorf("Failed to generate pay-to-address script, %v", err)
	}

	// Convert the amount to satoshi.
	satoshi, err := btcutil.NewAmount(*amount)
	if err != nil {
		return fmt.Errorf("Failed to convert amount, %v", err)
	}

	mtx.AddTxOut(wire.NewTxOut(int64(satoshi), pkScript))
	return nil
}
