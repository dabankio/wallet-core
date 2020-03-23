package omni



type Addr struct {
	Address string
	Privkey string
	Pubkey  string
}

var (
	presetAddrs = []Addr{
		Addr{Address: "mmr1MoTRVPTygJBgW4Wp3CKdifumZrUviM", Privkey: "cMuCfzEQznJ9k2obJfVjiDTzz9C1tTWALVweoYC8CL7TcpTGHhae", Pubkey: "035aa171ce873872468dde61d07e5921960ce11a6f0c42e35988a19d3201ed2498"},
		Addr{Address: "mu9ESzbfbXV4jeEvSSFcNDKb1yrwS6G93g", Privkey: "cSPhBkqgBxWoWS4qqchBEgEEkkntejs3fKKiFRU5Kkw3XCbrV8JP", Pubkey: "03c26c864529253892469ed705b1623114a387c1989b402992eec24b6a9f1c7dfb"},
		Addr{Address: "mg5fcedHzPXH3XLtbSm3mnsGuihLCCF6pa", Privkey: "cPvhEsZcBCka5X4kc3iEGWbegDFmsaWP5mmh79ZD3djbQZkGSsPH", Pubkey: "033f62b6222fb59fd90526604d1004795f9da9843ab0e3ec48d30092c0b558866c"},
		Addr{Address: "n4EDrbUE71iquaxAayozcpeYwYXeYRbrUU", Privkey: "cSTPvd9d8M3geUY98FmukjFqv5YmGY5J6RQTCdoF75uLr7ZKJEYj", Pubkey: "026655967cfc86e8175e214246f9bd2615076894ae96aaa10c25a445c0f0984469"},
		Addr{Address: "mt2PAKGMyTHWQCeG1hTkWo4TsyZPmL1T5a", Privkey: "cPaR52uVFyzV7PsVtvdUaDA1K8ZvKhtqHgUfaUo8csbJLozR2gCg", Pubkey: "034bcd6d9bfb0a189ac7893b778e1ce9f76f66e78e9cc4f8269f5c6083c702cb03"},
	}
)
