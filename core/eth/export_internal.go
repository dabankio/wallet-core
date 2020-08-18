package eth

import (
	internal "github.com/dabankio/wallet-core/core/eth/internalized"
)

//先这么解决现有的导出函数问题
var (
	New             = internal.New
	NewFromMetadata = internal.NewFromMetadata
)

// SignRawTransaction .
func SignRawTransaction(msg, privateKey string) (sig string, err error) {
	eth, _ := internal.New(nil)
	return eth.Sign(msg, privateKey)
}
