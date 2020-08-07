// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ERC20InterfaceABI is the input ABI used to generate the binding from.
// const ERC20InterfaceABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenOwner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenOwner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"tokenOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// ERC20InterfaceBin is the compiled bytecode used for deploying new contracts.
// const ERC20InterfaceBin = `0x`

// PackedDeployERC20Interface deploys a new Ethereum contract, binding an instance of ERC20Interface to it.
func PackedDeployERC20Interface() ([]byte, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20InterfaceABI))
	if err != nil {
		return nil, err
	}
	arguments, err := parsed.Constructor.Inputs.Pack()
	if err != nil {
		return nil, err
	}
	bytecode := common.FromHex(ERC20InterfaceBin)
	return append(bytecode, arguments...), nil
}

type ERC20InterfaceABIHelper struct {
	abi abi.ABI
}

func NewERC20InterfaceABIHelper() *ERC20InterfaceABIHelper {
	parsed, _ := abi.JSON(strings.NewReader(ERC20InterfaceABI))
	return &ERC20InterfaceABIHelper{parsed}
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
// Solidity: function allowance(address tokenOwner, address spender) constant returns(uint256 remaining)
func (_ERC20Interface *ERC20InterfaceABIHelper) PackedAllowance(tokenOwner common.Address, spender common.Address) ([]byte, error) {
	return _ERC20Interface.abi.Pack("allowance", tokenOwner, spender)
}

// UnpackAllowance is a free data retrieval call binding the contract method 0xdd62ed3e.
// Solidity: function allowance(address tokenOwner, address spender) constant returns(uint256 remaining)
func (_ERC20Interface *ERC20InterfaceABIHelper) UnpackAllowance(output []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Interface.abi.Unpack(out, "allowance", output)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
// Solidity: function balanceOf(address tokenOwner) constant returns(uint256 balance)
func (_ERC20Interface *ERC20InterfaceABIHelper) PackedBalanceOf(tokenOwner common.Address) ([]byte, error) {
	return _ERC20Interface.abi.Pack("balanceOf", tokenOwner)
}

// UnpackBalanceOf is a free data retrieval call binding the contract method 0x70a08231.
// Solidity: function balanceOf(address tokenOwner) constant returns(uint256 balance)
func (_ERC20Interface *ERC20InterfaceABIHelper) UnpackBalanceOf(output []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Interface.abi.Unpack(out, "balanceOf", output)
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Interface *ERC20InterfaceABIHelper) PackedTotalSupply() ([]byte, error) {
	return _ERC20Interface.abi.Pack("totalSupply")
}

// UnpackTotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Interface *ERC20InterfaceABIHelper) UnpackTotalSupply(output []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Interface.abi.Unpack(out, "totalSupply", output)
	return *ret0, err
}

// PackedApprove is a paid mutator transaction binding the contract method 0x095ea7b3.
// Solidity: function approve(address spender, uint256 tokens) returns(bool success)
func (_ERC20Interface *ERC20InterfaceABIHelper) PackedApprove(spender common.Address, tokens *big.Int) ([]byte, error) {
	return _ERC20Interface.abi.Pack("approve", spender, tokens)
}

// PackedTransfer is a paid mutator transaction binding the contract method 0xa9059cbb.
// Solidity: function transfer(address to, uint256 tokens) returns(bool success)
func (_ERC20Interface *ERC20InterfaceABIHelper) PackedTransfer(to common.Address, tokens *big.Int) ([]byte, error) {
	return _ERC20Interface.abi.Pack("transfer", to, tokens)
}

// PackedTransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
// Solidity: function transferFrom(address from, address to, uint256 tokens) returns(bool success)
func (_ERC20Interface *ERC20InterfaceABIHelper) PackedTransferFrom(from common.Address, to common.Address, tokens *big.Int) ([]byte, error) {
	return _ERC20Interface.abi.Pack("transferFrom", from, to, tokens)
}
