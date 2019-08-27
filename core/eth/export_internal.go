package eth

import (
	"github.com/lomocoin/wallet-core/core/eth/internal"
)

//先这么解决现有的导出函数问题
var (
	New             = internal.New
	NewFromMetadata = internal.NewFromMetadata
)
