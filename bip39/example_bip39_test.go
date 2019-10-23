package bip39_test

import (
	"fmt"
	"github.com/dabankio/wallet-core/bip39"
	"testing"
)

// example output:
// mnemonic: 流 铸 电 纬 近 俗 搭 插 仓 息 仲 惜
// is mnemonic valid: true
// seed len: 64
func TestBip39(t *testing.T) {
	// 设置语言
	bip39.SetWordListLang(bip39.LangChineseSimplified)

	// 创建熵
	ent, err := bip39.NewEntropy(128)
	if err != nil {
		t.Errorf("failed to new entropy: %v", err)
	}

	// 转换助记词
	mnemonic, err := bip39.NewMnemonic(ent)
	if err != nil {
		t.Errorf("failed to new mnemonic: %v", err)
	}
	fmt.Println("mnemonic:", mnemonic)

	// 验证助记词
	fmt.Println("is mnemonic valid:", bip39.IsMnemonicValid(mnemonic))

	// 构造种子
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		t.Errorf("failed to new seed: %v", err)
	}
	fmt.Println("seed len:", len(seed))

}
