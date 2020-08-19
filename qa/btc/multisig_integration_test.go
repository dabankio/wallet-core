package btc

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dabankio/devtools4chains"
	"github.com/dabankio/wallet-core/core/btc"
	"github.com/dabankio/wallet-core/core/eth/internalized/testtool"
	"github.com/stretchr/testify/require"
)

// ListUnspentResult models a successful response from the listunspent request.
type ListUnspentResult struct {
	TxID          string  `json:"txid"`
	Vout          uint32  `json:"vout"`
	Address       string  `json:"address"`
	Account       string  `json:"account"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	RedeemScript  string  `json:"redeemScript,omitempty"`
	Amount        float64 `json:"amount"`
	Confirmations int64   `json:"confirmations"`
	Spendable     bool    `json:"spendable"`

	Label         string `json:"label"`         //        (string) The associated label, or "" for the default label
	WitnessScript string `json:"witnessScript"` // (string) witnessScript if the scriptPubKey is P2WSH or P2SH-P2WSH
	Solvable      bool   `json:"solvable"`      //         (bool) Whether we know how to spend this output, ignoring the lack of keys
	Desc          string `json:"desc"`          //             (string, only when solvable) A descriptor for spending this output
	Safe          bool   `json:"safe"`          //             (bool) Whether this output is considered safe to spend. Unconfirmed transactions from outside keys and unconfirmed replacement transactions are considered unsafe and are not eligible for spending by fundrawtransaction and sendtoaddress.
}

// 多重签名测试
func TestMultisig(t *testing.T) {
	rq := require.New(t)

	image := "ruimarinho/bitcoin-core:latest"
	killFunc, bitcoinInfo, err := devtools4chains.DockerRunBitcoin(devtools4chains.DockerRunOptions{
		AutoRemove: true, Image: &image,
	})
	require.NoError(t, err)
	t.Cleanup(killFunc)

	rpcInfo := devtools4chains.RPCInfo{
		Host:     fmt.Sprintf("http://127.0.0.1:%d", bitcoinInfo.RPCPort),
		User:     bitcoinInfo.RPCUser,
		Password: bitcoinInfo.RPCPwd,
	}

	// 导入a0 private Key, a1 a2 a3 address
	// 首先为a0生成 101 个块
	// 用a1 a2 a3 生成多签地址 (2-3)
	// 往多签地址转入btc
	// a1 a2 签名，转出到 a3
	// 查询a3 utxo, 应该不为0

	testtool.WaitSomething(t, time.Minute, func() error {
		b, err := devtools4chains.RPCCallJSON(rpcInfo, "getblockcount", nil, nil)
		// _, err := devtools4chains.RPCCallJSON(rpcInfo, "getblockcount", []interface{}{})
		if b != nil && strings.Contains(string(b), "Loading wallet") {
			return fmt.Errorf("Loading wallet")
		}
		return err
	})

	{ // import addresses
		_, err := devtools4chains.RPCCallJSON(rpcInfo, "importprivkey", []string{a0.Privkey}, nil)
		rq.Nil(err)

		for _, add := range []string{a1.Address, a2.Address, a3.Address} {
			_, err := devtools4chains.RPCCallJSON(rpcInfo, "importaddress", []string{add}, nil)
			rq.Nil(err)
		}
	}

	_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{101, a0.Address}, nil)
	rq.Nil(err)

	var multisigAddress, redeemScript string
	{ // create multisig address,and import to bitcoind
		rs, err := btc.NewMultiSigAddress(2, btc.ChainRegtest, strings.Join([]string{a1.Pubkey, a2.Pubkey, a3.Pubkey}, ","))
		rq.Nil(err)
		arr := strings.Split(rs, ",")
		rq.Len(arr, 2, "")
		multisigAddress, redeemScript = arr[0], arr[1]

		_, err = devtools4chains.RPCCallJSON(rpcInfo, "importaddress", []string{multisigAddress}, nil)
		rq.Nil(err)
	}

	{ //send to multisig address for next step
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendtoaddress", []interface{}{multisigAddress, 23.3}, nil)
		rq.Nil(err)
	}

	var unspents []ListUnspentResult
	resp, err := devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{multisigAddress}}, nil)
	rq.Nil(err)
	rq.NoError(json.Unmarshal(resp, &unspents), string(resp))
	rq.Len(unspents, 1, string(resp))

	utxo := unspents[0]
	fmt.Printf("%#v\n", utxo)

	var signedHex string
	transferAmount := 2.3
	chainID := btc.ChainRegtest
	{ //SDK 的使用可以参考这里
		var tx *btc.BTCTransaction
		unspent := new(btc.BTCUnspent) //java: new btc.BTCUnspent()
		{                              // build tx
			unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, redeemScript)

			amount, err := btc.NewBTCAmount(transferAmount)
			rq.Nil(err)

			toAddressA1, err := btc.NewBTCAddressFromString(a1.Address, chainID)
			rq.Nil(err)

			outputAmount := btc.BTCOutputAmount{} //java: new btc.BTCOutputAmount()
			outputAmount.Add(toAddressA1, amount)

			feeRate := int64(80)

			changeAddressMultisig, err := btc.NewBTCAddressFromString(multisigAddress, chainID)
			rq.Nil(err)

			tx, err = btc.NewBTCTransaction(unspent, &outputAmount, changeAddressMultisig, feeRate, chainID)
			rq.Nil(err)
		}

		{ //sdk sign
			// member a1 sign
			signRs, err := btc.SignTransaction(tx, a1.Privkey, chainID)
			rq.Nil(err)
			rq.True(signRs.Changed, "")
			rq.False(signRs.Complete, "")
			signedHex = signRs.Hex

			// member a3 sign
			signRs, err = btc.SignRawTransactionWithKey(signedHex, a3.Privkey, unspent, chainID)
			rq.Nil(err)
			rq.True(signRs.Changed)
			rq.True(signRs.Complete)
			signedHex = signRs.Hex
		}
	}

	{ // relay tx
		_, err := devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{signedHex}, nil)
		rq.Nil(err)
	}

	{ //generate 1 block
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{1, a0.Address}, nil)
		rq.Nil(err)
	}

	{ // validate utxo for receiver
		resp, err := devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{a1.Address}}, nil)
		rq.Nil(err)

		fmt.Println("utxo for a1", string(resp))
	}

}
