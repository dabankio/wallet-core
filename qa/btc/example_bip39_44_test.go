package btc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dabankio/wallet-core/bip39"
	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/core/btc"
	"github.com/stretchr/testify/require"
)

// btc bip39 bip44 示例代码
func TestExampleBip39_44(t *testing.T) {
	rq := require.New(t)

	bip39.SetWordListLang(bip39.LangChineseSimplified)
	ent, err := bip39.NewEntropy(128)
	rq.Nil(err)

	mnemonic, err := bip39.NewMnemonic(ent)
	rq.Nil(err)
	fmt.Println("mnemonic:", mnemonic)

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	rq.Nil(err)

	for _, chainID := range []int{btc.ChainMainNet, btc.ChainTestNet3, btc.ChainRegtest} {
		deriver, err := btc.NewBip44Deriver(bip44.PathFormat, seed, chainID)
		rq.Nil(err)

		address, err := deriver.DeriveAddress()
		rq.Nil(err)
		prvk, err := deriver.DerivePrivateKey()
		rq.Nil(err)
		pubk, err := deriver.DerivePublicKey()
		rq.Nil(err)
		fmt.Printf("chain [%d] address/private/public: %s/%s/%s\n", chainID, address, prvk, pubk)
	}
}

// 该测试展示了如何生成一个“靓号”，大致就是一直试，直到找到一个助记词可以生成对应规则的地址
// 对于随机种子，假设指定n位（大小写敏感），那么单次成功的概率为1/58^n，比如指定以"BTC"开头，那么理论上单次成功的概率为1/58*58*58=1/195112，
func TestExampleGenerateBeautifulAddress(t *testing.T) {

	t.Skip("作为示例用，注释或删除这行代码使用测试")

	const beautifulPrefix = "1btc"
	var address string
	var mnemonic string
	bip39.SetWordListLang(bip39.LangChineseSimplified)
	chainID := btc.ChainMainNet
	count := 0
	for !strings.HasPrefix(strings.ToLower(address), beautifulPrefix) { //此处大小写不敏感
		ent, _ := bip39.NewEntropy(128)
		mnemonic, _ = bip39.NewMnemonic(ent)
		seed, _ := bip39.NewSeedWithErrorChecking(mnemonic, "")
		deriver, _ := btc.NewBip44Deriver(bip44.PathFormat, seed, chainID)
		address, _ = deriver.DeriveAddress()
		count++
		if count%1000 == 0 { //每1000次打印日志
			fmt.Println("count", count)
		}
	}
	fmt.Println("address", address)
	fmt.Println("mnemonic", mnemonic)
}
