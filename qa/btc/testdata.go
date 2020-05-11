package btc

type Addr struct {
	Address string
	Privkey string
	Pubkey  string
}

var (
	addrs = [5]Addr{
		Addr{Address: "mwuHuUeUF8oUjjFKFaURYtTUSZwM1zGgwW", Privkey: "cRMM3dxdoQZjdo8CYJuoZAKZJj5kWxY7QHwws7u5y96nEM8NjoZD", Pubkey: "02540da55aacbb3ea8a072e137bbd177180c1649482f62d9abb1e66c1db59a4d1c"},
		Addr{Address: "mqQDjZPwFhRWucTPbuULPYA8ouwEnqbWEQ", Privkey: "cQfn7uKptR51ga5RSpFB7kDAY6qT44reYfwaWEi1drzhP3srfBvA", Pubkey: "02e06a99f8b6ad7ed986a6d13429cced62c71127dff704ba21028640f816e4afd8"},
		Addr{Address: "mxkC1eUFX284kpTMPFx5KUqdumRAXdPxoB", Privkey: "cQUgmhQwh7AEmZE5kCncmUV9kmuPN75TGbaJR71uWcyXrAMUAyGe", Pubkey: "024078d3e4c3d14659c02ace5eb03c5c1ec041cc2a0ab4fb23e72f9955aef4b024"},
		Addr{Address: "mkRFT89Ktr4URDDfxeskKCqr8kD7RNgMsH", Privkey: "cRSXZBdjXpvXykdTyLT7iNneJHsroTZ5u5zrVzLuuR5XnmQwGZwk", Pubkey: "034f3ba1aeb2bba9481663986af5e480bd96e20a4b0ea66215d9aee914175786f7"},
		Addr{Address: "mmUcpuVGLST4oGsmouCChtTYuULPKGt49z", Privkey: "cRbPWnf9J1EV7xEXKv6D2R9zTKjVeMHfVFv47gEsKwM6FQdpvpws", Pubkey: "02372b895870068686bb748866702a8093685796ae1c177ce5adeb0db92f0572e7"},
	}

	a0, a1, a2, a3, a4 = addrs[0], addrs[1], addrs[2], addrs[3], addrs[4]
)
