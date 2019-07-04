{{- /* （模板注释，不会生成到最终文件中）该模板生成合约的abi工具类，可以帮助打包、解包合约函数数据，修改自 https://github.com/ethereum/go-ethereum/blob/master/accounts/abi/bind/template.go */ -}}
// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package {{.Package}}

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

{{range $contract := .Contracts}}
	// {{.Type}}ABI is the input ABI used to generate the binding from.
	const {{.Type}}ABI = "{{.InputABI}}"

	{{if issignalexists "s_gen_bin"}}
	{{if .InputBin}}
		// {{.Type}}Bin is the compiled bytecode used for deploying new contracts.
		const {{.Type}}Bin = `{{.InputBin}}`

		// PackedDeploy{{.Type}} deploys a new Ethereum contract, binding an instance of {{.Type}} to it.
		func PackedDeploy{{.Type}}({{range $idx, $e := .Constructor.Inputs}}{{if gt $idx 0}},{{end}} {{.Name}} {{bindmobiletypego .Type true}}{{end}}) ([]byte, error) {
		  parsed, err := abi.JSON(strings.NewReader({{.Type}}ABI))
		  if err != nil {
		    return nil, err
		  }
		arguments, err := parsed.Constructor.Inputs.Pack({{range $idx, $e := .Constructor.Inputs}}{{if gt $idx 0}},{{end}} {{.Name}}{{if iswrapgotype .Type}}.{{gotypewrapfield .Type}}{{end}}{{end}})
		if err != nil {
			return nil, err
		}
		bytecode := common.FromHex({{.Type}}Bin)
		return append(bytecode, arguments...), nil
		}
	{{end}}
	{{end}}

	// {{.Type}}ABIHelper tool for contract abi
	type {{.Type}}ABIHelper struct {
		abi abi.ABI
	}

	// New{{.Type}}ABIHelper constructor
	func New{{.Type}}ABIHelper() *{{.Type}}ABIHelper {
		parsed, _ := abi.JSON(strings.NewReader({{.Type}}ABI))
		return &{{.Type}}ABIHelper{parsed}
	}

	{{range .Calls}}
		// Packed{{.Normalized.Name}} is a free data retrieval call binding the contract method 0x{{printf "%x" .Original.Id}}.
		// Solidity: {{.Original.String}}
		func (_{{$contract.Type}} *{{$contract.Type}}ABIHelper) Packed{{.Normalized.Name}}({{range $idx, $e := .Normalized.Inputs}}{{if gt $idx 0}},{{end}} {{.Name}} {{bindmobiletypego .Type true}} {{end}}) ([]byte, error) {
			return _{{$contract.Type}}.abi.Pack("{{.Original.Name}}" {{range .Normalized.Inputs}}, {{.Name}}{{if iswrapgotype .Type}}.{{gotypewrapfield .Type}}{{end}} {{end}})
		}

		// Unpack{{.Normalized.Name}} is a free data retrieval call binding the contract method 0x{{printf "%x" .Original.Id}}.
		// Solidity: {{.Original.String}}
		func (_{{$contract.Type}} *{{$contract.Type}}ABIHelper) Unpack{{.Normalized.Name}}(output []byte) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type}};{{end}} },{{else}}{{range .Normalized.Outputs}}{{bindmobiletypego .Type true}},{{end}}{{end}} error) {
			{{if .Structured}}ret := new(struct{
				{{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type}}
				{{end}}
			}){{else}}var (
				{{range $i, $_ := .Normalized.Outputs}}ret{{$i}} = new({{bindtype .Type}})
				{{end}}
			){{end}}
			out := {{if .Structured}}ret{{else}}{{if eq (len .Normalized.Outputs) 1}}ret0{{else}}&[]interface{}{
				{{range $i, $_ := .Normalized.Outputs}}ret{{$i}},
				{{end}}
			}{{end}}{{end}}
			err := _{{$contract.Type}}.abi.Unpack(out, "{{.Original.Name}}", output)
			return {{if .Structured}}*ret,{{else}}{{range $i, $_ := .Normalized.Outputs}}{{if iswrapgotype .Type}}{{bindmobiletypego .Type false}}{*ret{{$i}}}{{else}}*ret{{$i}}{{end}},{{end}}{{end}} err
		}
	{{end}}

	{{range .Transacts}}
		// Packed{{.Normalized.Name}} is a paid mutator transaction binding the contract method 0x{{printf "%x" .Original.Id}}.
		// Solidity: {{.Original.String}}
		func (_{{$contract.Type}} *{{$contract.Type}}ABIHelper) Packed{{.Normalized.Name}}({{range $idx, $e := .Normalized.Inputs}}{{if gt $idx 0}},{{end}} {{.Name}} {{bindmobiletypego .Type true}} {{end}}) ([]byte, error) {
			return _{{$contract.Type}}.abi.Pack("{{.Original.Name}}" {{range .Normalized.Inputs}}, {{.Name}}{{if iswrapgotype .Type}}.{{gotypewrapfield .Type}}{{end}}{{end}})
		}
	{{end}}

{{end}}