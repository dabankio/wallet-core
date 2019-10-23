// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Contains all the wrappers from the core/types package.

package eth

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// ETHTransaction represents a single Ethereum transaction.
type ETHTransaction struct {
	tx *types.Transaction
}

// NewETHTransaction creates a new ether transaction with the given properties.
func NewETHTransaction(nonce int64, to *ETHAddress, amount *BigInt, gasLimit int64, gasPrice *BigInt, data []byte) *ETHTransaction {
	return &ETHTransaction{types.NewTransaction(uint64(nonce), to.address, amount.bigint, uint64(gasLimit), gasPrice.bigint, common.CopyBytes(data))}
}

// NewETHTransactionForContractCreation 创建一个部署合约交易(amount 为 0)
func NewETHTransactionForContractCreation(nonce int64, gasLimit int64, gasPrice *BigInt, data []byte) *ETHTransaction {
	tx := types.NewContractCreation(uint64(nonce), big.NewInt(0), uint64(gasLimit), gasPrice.bigint, data)
	return &ETHTransaction{tx: tx}
}

// NewETHTransactionFromRLP parses a transaction from an RLP data dump.
func newETHTransactionFromRLP(data []byte) (*ETHTransaction, error) {
	tx := &ETHTransaction{
		tx: new(types.Transaction),
	}
	if err := rlp.DecodeBytes(common.CopyBytes(data), tx.tx); err != nil {
		return nil, err
	}
	return tx, nil
}

// EncodeRLP encodes a transaction into an hex RLP data dump.
func (tx *ETHTransaction) EncodeRLP() (string, error) {
	data, err := rlp.EncodeToBytes(tx.tx)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(data), nil
}

// newETHTransactionFromJSON parses a transaction from a JSON data dump.
func newETHTransactionFromJSON(data string) (*ETHTransaction, error) {
	tx := &ETHTransaction{
		tx: new(types.Transaction),
	}
	if err := json.Unmarshal([]byte(data), tx.tx); err != nil {
		return nil, err
	}
	return tx, nil
}

// EncodeJSON encodes a transaction into a JSON data dump.
func (tx *ETHTransaction) EncodeJSON() (string, error) {
	data, err := json.Marshal(tx.tx)
	return string(data), err
}

// GetData .
func (tx *ETHTransaction) GetData() []byte { return tx.tx.Data() }

// GetGas .
func (tx *ETHTransaction) GetGas() int64 { return int64(tx.tx.Gas()) }

// GetGasPrice .
func (tx *ETHTransaction) GetGasPrice() *BigInt { return &BigInt{tx.tx.GasPrice()} }

// GetValue .
func (tx *ETHTransaction) GetValue() *BigInt { return &BigInt{tx.tx.Value()} }

// GetNonce .
func (tx *ETHTransaction) GetNonce() int64 { return int64(tx.tx.Nonce()) }

// GetHash .
func (tx *ETHTransaction) GetHash() *Hash { return &Hash{tx.tx.Hash()} }

// GetCost .
func (tx *ETHTransaction) GetCost() *BigInt { return &BigInt{tx.tx.Cost()} }

// GetTo .
func (tx *ETHTransaction) GetTo() *ETHAddress {
	if to := tx.tx.To(); to != nil {
		return &ETHAddress{*to}
	}
	return nil
}
