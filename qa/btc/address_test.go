package btc

import (
	"fmt"
	"testing"

	devTools4 "github.com/dabankio/devtools4chains"
	"github.com/stretchr/testify/require"
)

func TestSegwit(t *testing.T) {
	image := "ruimarinho/bitcoin-core:latest"
	killFunc, bitcoinInfo, err := devTools4.DockerRunBitcoin(devTools4.DockerRunOptions{
		AutoRemove: true, Image: &image,
	})
	require.NoError(t, err)
	t.Cleanup(killFunc)

	rpcInfo := devTools4.RPCInfo{
		Host:     fmt.Sprintf("http://127.0.0.1:%d", bitcoinInfo.RPCPort),
		User:     bitcoinInfo.RPCUser,
		Password: bitcoinInfo.RPCPwd,
	}

	var ret string
	_, err = devTools4.RPCCallJSON(rpcInfo, "getnewaddress", []string{"", "p2sh-segwit"}, &ret)
	require.NoError(t, err)
	fmt.Println("segwit addr:", ret)

	_, err = devTools4.RPCCallJSON(rpcInfo, "dumpprivkey", []string{ret}, &ret)
	require.NoError(t, err)
	fmt.Println("private key:", ret)

}
