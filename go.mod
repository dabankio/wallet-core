module github.com/dabankio/wallet-core

go 1.12

require (
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412
	github.com/aristanetworks/goarista v0.0.0-20190319235110-489128639c40 // indirect
	github.com/btcsuite/btcd v0.0.0-20190315201642-aa6e0f35703c
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f
	github.com/btcsuite/btcutil v0.0.0-20190316010144-3ac1210f4b38
	github.com/cespare/cp v1.1.1 // indirect
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/edsrzf/mmap-go v1.0.0 // indirect
	github.com/elastic/gosigar v0.10.5 // indirect
	github.com/ethereum/go-ethereum v1.9.12
	github.com/fjl/memsize v0.0.0-20180929194037-2a09253e352a // indirect
	github.com/golang/protobuf v1.3.2-0.20190517061210-b285ee9cfc6c
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/huin/goupnp v1.0.0 // indirect
	github.com/lomocoin/gobbc v1.0.2
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/rs/cors v1.6.0 // indirect
	github.com/status-im/keycard-go v0.0.0-20190424133014-d95853db0f48 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/tyler-smith/go-bip39 v1.0.2 // indirect
	golang.org/x/crypto v0.0.0-20200311171314-f7b00557c8c4
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20190709231704-1e4459ed25ff // indirect
)

// build ios 需要处理下面这个问题
//解决eth相关的下面这个报错
// fatal error: 'libproc.h' file not found
// #include <libproc.h>
//          ^~~~~~~~~~~
// clone github.com/celo-org/gosigar 切换到修复了这个问题的分支
// 这个问题，后续等bug fix PR 合并了，可以删除下面的replace
// replace github.com/elastic/gosigar v0.10.5 => /Users/dev/Documents/workspace/github.com/celo-org/gosigar
