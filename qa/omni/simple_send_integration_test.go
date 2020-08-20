package omni

import (
	"fmt"
	"testing"
	"time"

	"github.com/dabankio/devtools4chains"
	"github.com/dabankio/wallet-core/core/btc"
	"github.com/dabankio/wallet-core/core/eth/internalized/testtool"
	"github.com/dabankio/wallet-core/core/omni"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testImage = "mpugach/omnicored:v0.8.2-alpine"

// 简单转账测试
// 本地起一条全新的链，发布一个omni代币，在此基础上进行测试
func TestSimpleSend(t *testing.T) {
	rq := require.New(t)

	killFunc, nodeInfo, err := devtools4chains.DockerRunOmnicored(devtools4chains.DockerRunOptions{
		AutoRemove: true,
		Image:      &testImage,
	})
	rq.NoError(err)
	t.Cleanup(killFunc)

	rpcInfo := devtools4chains.RPCInfo{
		Host:     fmt.Sprintf("http://127.0.0.1:%d", nodeInfo.RPCPort),
		User:     nodeInfo.RPCUser,
		Password: nodeInfo.RPCPwd,
		// Debug: true,
	}

	testtool.WaitSomething(t, time.Minute, func() error {
		_, err := devtools4chains.RPCCallJSON(rpcInfo, "getblockcount", nil, nil)
		return err
	})
	for _, add := range presetAddrs {
		_, err := devtools4chains.RPCCallJSON(rpcInfo, "importprivkey", []string{add.Privkey}, nil)
		rq.Nil(err)
	}

	a0, a1 := presetAddrs[0], presetAddrs[1]
	// time.Sleep(time.Second)

	{ // 生成多个块，获取utxo
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{101, a0.Address}, nil)
		rq.Nil(err)
	}

	propertyID := createToken(t, rpcInfo, OmniSenddissuancefixedCmd{
		Fromaddress: a0.Address,
		Ecosystem:   2, //2 fot test
		Typ:         1, // 1 for indivisible
		Previousid:  0, // 0 for new tokens
		Category:    "test_omniii",
		Subcategory: "unit_test",
		Name:        "FakeUSDT",
		Amount:      "10000",
	})

	var unspents []ListUnspentResult
	_, err = devtools4chains.RPCCallJSON(rpcInfo, "listunspent", []interface{}{0, 999, []string{a0.Address}}, &unspents)
	rq.Nil(err)
	utxo := unspents[0]

	chainID := btc.ChainRegtest
	transferAmount := 233.0
	{ // simple send from a0 to a1,
		unspent := new(btc.BTCUnspent)
		unspent.Add(utxo.TxID, int64(utxo.Vout), utxo.Amount, utxo.ScriptPubKey, utxo.RedeemScript)
		toAddr, err := btc.NewBTCAddressFromString(a1.Address, chainID)
		rq.Nil(err, "failed to create addr")
		changeAddr, err := btc.NewBTCAddressFromString(a0.Address, chainID)
		rq.Nil(err, "failed to create addr")
		feeRate := int64(80)
		btctx, err := omni.CreateSimpleSendTransaction(propertyID, false, unspent, toAddr, transferAmount, changeAddr, feeRate, chainID)
		rq.Nil(err, "Failed to crate btctx")

		toSignMsg, err := btctx.EncodeToSignCmd()
		rq.Nil(err, "failed to encode to sign")

		btcCoin, _ := btc.New(nil, chainID)
		signedRawHex, err := btcCoin.Sign(toSignMsg, a0.Privkey)
		assert.Nil(t, err)

		// 广播交易
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "sendrawtransaction", []interface{}{signedRawHex}, nil)
		rq.Nil(err)
	}

	{ // 生成一个块确认代币转账
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{1, a0.Address}, nil) //生成几个块，确认刚才的交易
		rq.Nil(err)
	}

	{ // 确认代币转账成功
		var balanceMap map[string]string
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "omni_getbalance", []interface{}{a1.Address, propertyID}, &balanceMap) //生成几个块，确认刚才的交易
		rq.Nil(err, "Failed to get omni balance")
		rq.Equal(balanceMap["balance"], "233")
	}
}

type OmniSenddissuancefixedCmd struct {
	Fromaddress                                    string
	Ecosystem                                      int
	Typ                                            int
	Previousid                                     int
	Category, Subcategory, Name, URL, Data, Amount string
}

// 发行代币
func createToken(t *testing.T, rpcInfo devtools4chains.RPCInfo, cmd OmniSenddissuancefixedCmd) (propertyID int) {
	rq := require.New(t)
	{ // create token
		var txHash string
		_, err := devtools4chains.RPCCallJSON(rpcInfo, "omni_sendissuancefixed", []interface{}{
			cmd.Fromaddress, cmd.Ecosystem, cmd.Typ, cmd.Previousid,
			cmd.Category, cmd.Subcategory, cmd.Name, cmd.URL, cmd.Data, cmd.Amount,
		}, &txHash)
		rq.Nil(err, "Failed to create omni coin")

		_, err = devtools4chains.RPCCallJSON(rpcInfo, "generatetoaddress", []interface{}{1, cmd.Fromaddress}, nil) //生成几个块，确认刚才的交易
		rq.Nil(err)

		var tx map[string]interface{}
		_, err = devtools4chains.RPCCallJSON(rpcInfo, "omni_gettransaction", []string{txHash}, &tx)
		rq.NoError(err)
		propertyID = int(tx["propertyid"].(float64))
		rq.False(propertyID == 0, "Got property id error", propertyID)
	}

	{ // 代币创建完成后查询代币持有人的余额，应该等于总的发行量
		// fmt.Println("-------then balance of new created property-----")
		// bal, err := cli.OmniGetbalance(cmd.Fromaddress, propertyID)
		// rq.Nil(err, "Failed to get balance of owner")
		// rq.False(bal.Balance != cmd.Amount, "余额不符合预期")
	}
	return
}

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
