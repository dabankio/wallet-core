package btc

import (
	"github.com/lomocoin/wallet-core/core/btc/internal"
)

// Btc (全部大写在导出到java那边有点问题)
type Btc struct {
	internal.BTC
}

//暂时先这么解决现有的代码依赖问题
var (
	New             = internal.New
	NewFromMetadata = internal.NewFromMetadata
)
