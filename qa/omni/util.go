package omni

import (
	"testing"

	"github.com/dabankio/devtools4chains"
	"github.com/stretchr/testify/require"
)

type OmniSenddissuancefixedCmd struct {
	Fromaddress                                    string
	Ecosystem                                      int
	Typ                                            int
	Previousid                                     int
	Category, Subcategory, Name, URL, Data, Amount string
}

// CreateToken 发行代币
func CreateToken(t *testing.T, rpcInfo devtools4chains.RPCInfo, cmd OmniSenddissuancefixedCmd) (propertyID int) {
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
