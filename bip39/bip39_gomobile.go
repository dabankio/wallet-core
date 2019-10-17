package bip39

import (
	"strings"

	"github.com/dabankio/wallet-core/bip39/wordlists"
	"github.com/pkg/errors"
)

// wrap some func for gomobile export

// GetWordListString 当前词汇表（以,分隔）
func GetWordListString() string {
	return strings.Join(GetWordList(), ",")
}

const (
	// LangChineseSimplified 简体中文
	LangChineseSimplified = iota
	// LangChineseTraditional 繁体中文
	LangChineseTraditional
	// LangEnglish 英语
	LangEnglish
	// LangFrench 法语
	LangFrench
	// LangItalian 意大利语
	LangItalian
	// LangJapanese 日文
	LangJapanese
	// LangKorean 韩文
	LangKorean
	// LangSpanish 西班牙语
	LangSpanish
)

// SetWordListLang 设置词汇表语言(默认英语)
func SetWordListLang(lang int) error {
	switch lang {
	case LangChineseSimplified:
		SetWordList(wordlists.ChineseSimplified)
	case LangChineseTraditional:
		SetWordList(wordlists.ChineseTraditional)
	case LangEnglish:
		SetWordList(wordlists.English)
	case LangFrench:
		SetWordList(wordlists.French)
	case LangItalian:
		SetWordList(wordlists.Italian)
	case LangJapanese:
		SetWordList(wordlists.Japanese)
	case LangKorean:
		SetWordList(wordlists.Korean)
	case LangSpanish:
		SetWordList(wordlists.Spanish)
	default:
		return errors.Errorf("expected lang [%d - %d], got: %d", LangChineseSimplified, LangSpanish, lang)
	}
	return nil
}

// func GetWordIndex(word string) (int, bool) {
// func MnemonicToByteArray(mnemonic string, raw ...bool) ([]byte, error) {
