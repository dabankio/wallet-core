package wallet

const bbcCoreImage = "dabankio/bbccore:0.11"
const mkfImage = "dabankio/mkfdev:1"
const mkfImageDexTest = "dabankio/mkf-dex-test:20201009"

var bbcImage = "ruimarinho/bitcoin-core:latest"
var omniImage = "mpugach/omnicored:v0.8.2-alpine"

type ctx struct {
	pubk, address string
}
