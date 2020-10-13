package trx

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/dabankio/wallet-core/core/trx"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/golang/protobuf/proto"
	r "github.com/stretchr/testify/require"
)

func BuildTX(t *testing.T, tx *api.TransactionExtention) (txMsg string, txid string) {
	msg, err := proto.Marshal(tx.GetTransaction())
	r.NoError(t, err)
	return hex.EncodeToString(msg), hex.EncodeToString(tx.Txid)
}

func BroadcastTX(t *testing.T, sig string, grpc *client.GrpcClient) {
	trRawBytes, err := hex.DecodeString(sig)
	r.NoError(t, err)
	tr := new(core.Transaction)
	err = proto.Unmarshal(trRawBytes, tr)
	r.NoError(t, err)

	ret, err := grpc.Broadcast(tr)
	r.NoError(t, err)
	fmt.Println("Broadcast RET", ret.Code, string(ret.Message), ret)
}

func SignBroadcastTx(t *testing.T, tx *api.TransactionExtention, privk string, grpc *client.GrpcClient) string {
	msg, txid := BuildTX(t, tx)
	sig, err := trx.SignWithPrivateKey(msg, privk)
	r.NoError(t, err)
	BroadcastTX(t, sig, grpc)
	return txid
}
