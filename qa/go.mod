module github.com/dabankio/wallet-core/qa

go 1.14

require (
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/dabankio/bbrpc v1.1.1-beta.3
	github.com/dabankio/devtools4chains v0.0.0-20200218063219-f22d17f996f4
	github.com/dabankio/wallet-core v0.0.0-00010101000000-000000000000
	github.com/ethereum/go-ethereum v1.9.12
	github.com/stretchr/testify v1.5.1
)

replace github.com/dabankio/wallet-core => ../
