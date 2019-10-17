package omni

import (
	"github.com/dabankio/omnicli"
	"github.com/dabankio/omnicli/btcjson"
)

func importAddrs(cli *omnicli.Cli, presetAddrs []omnicli.Addr) error {
	for _, add := range presetAddrs {
		e := cli.Importprivkey(btcjson.ImportPrivKeyCmd{
			PrivKey: add.Privkey,
		})
		if e != nil {
			return e
		}

		// e = omnicli.CliImportaddress(btcjson.ImportAddressCmd{
		// 	Address: add.Address,
		// })
		// if e != nil {
		// 	return e
		// }
	}
	return nil
}
