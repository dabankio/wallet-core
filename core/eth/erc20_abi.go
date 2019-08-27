// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package geth

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

// ERC20InterfaceABI is the input ABI used to generate the binding from.
const ERC20InterfaceABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenOwner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenOwner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"tokenOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// ERC20InterfaceABIHelper tool for contract abi
type ERC20InterfaceABIHelper struct {
	abi abi.ABI
}

// NewERC20InterfaceABIHelper constructor
func NewERC20InterfaceABIHelper() *ERC20InterfaceABIHelper {
	parsed, _ := abi.JSON(strings.NewReader(ERC20InterfaceABI))
	return &ERC20InterfaceABIHelper{parsed}
}

// PackedAllowance is a free data retrieval call binding the contract method 0xdd62ed3e.
// Solidity: function allowance(address tokenOwner, address spender) constant returns(uint256 remaining)
func (_ERC20Interface *ERC20InterfaceABIHelper) PackedAllowance(tokenOwner *ETHAddress, spender *ETHAddress) ([]byte, error) {
	return _ERC20Interface.abi.Pack("allowance", tokenOwner.address, spender.address)
}

// UnpackAllowance is a free data retrieval call binding the contract method 0xdd62ed3e.
// Solidity: function allowance(address tokenOwner, address spender) constant returns(uint256 remaining)
func (_ERC20Interface *ERC20InterfaceABIHelper) UnpackAllowance(output []byte) (*BigInt, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Interface.abi.Unpack(out, "allowance", output)
	return &BigInt{*ret0}, err
}

// PackedBalanceOf is a free data retrieval call binding the contract method 0x70a08231.
// Solidity: function balanceOf(address tokenOwner) constant returns(uint256 balance)
func (_ERC20Interface *ERC20InterfaceABIHelper) PackedBalanceOf(tokenOwner *ETHAddress) ([]byte, error) {
	return _ERC20Interface.abi.Pack("balanceOf", tokenOwner.address)
}

// UnpackBalanceOf is a free data retrieval call binding the contract method 0x70a08231.
// Solidity: function balanceOf(address tokenOwner) constant returns(uint256 balance)
func (_ERC20Interface *ERC20InterfaceABIHelper) UnpackBalanceOf(output []byte) (*BigInt, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Interface.abi.Unpack(out, "balanceOf", output)
	return &BigInt{*ret0}, err
}

// PackedTotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Interface *ERC20InterfaceABIHelper) PackedTotalSupply() ([]byte, error) {
	return _ERC20Interface.abi.Pack("totalSupply")
}

// UnpackTotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Interface *ERC20InterfaceABIHelper) UnpackTotalSupply(output []byte) (*BigInt, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20Interface.abi.Unpack(out, "totalSupply", output)
	return &BigInt{*ret0}, err
}

// PackedApprove is a paid mutator transaction binding the contract method 0x095ea7b3.
// Solidity: function approve(address spender, uint256 tokens) returns(bool success)
func (_ERC20Interface *ERC20InterfaceABIHelper) PackedApprove(spender *ETHAddress, tokens *BigInt) ([]byte, error) {
	return _ERC20Interface.abi.Pack("approve", spender.address, tokens.bigint)
}

// PackedTransfer is a paid mutator transaction binding the contract method 0xa9059cbb.
// Solidity: function transfer(address to, uint256 tokens) returns(bool success)
func (_ERC20Interface *ERC20InterfaceABIHelper) PackedTransfer(to *ETHAddress, tokens *BigInt) ([]byte, error) {
	return _ERC20Interface.abi.Pack("transfer", to.address, tokens.bigint)
}

// PackedTransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
// Solidity: function transferFrom(address from, address to, uint256 tokens) returns(bool success)
func (_ERC20Interface *ERC20InterfaceABIHelper) PackedTransferFrom(from *ETHAddress, to *ETHAddress, tokens *BigInt) ([]byte, error) {
	return _ERC20Interface.abi.Pack("transferFrom", from.address, to.address, tokens.bigint)
}
