package geth

// 该文件对一些用到的go类型进行了封装，使得可以用gomobile导出给客户端

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

// NewAddressesWrap .
func NewAddressesWrap() *AddressesWrap {
	return &AddressesWrap{}
}

// AddressesWrap 地址数组包装
type AddressesWrap struct {
	wrap []common.Address
}

// AddOne 追加一个地址
func (w *AddressesWrap) AddOne(addr *ETHAddress) {
	w.wrap = append(w.wrap, addr.address)
}

// ------------------------------------------------------------------------------------------------------------------------------

// NewByte32ArrayWrap .
func NewByte32ArrayWrap() *Byte32ArrayWrap {
	return &Byte32ArrayWrap{}
}

// Byte32ArrayWrap wrap [][32]byte
type Byte32ArrayWrap struct {
	wrap [][32]byte
}

// AddOne len of bytes shoule be 32
func (w *Byte32ArrayWrap) AddOne(b []byte) error {
	if l := len(b); l != 32 {
		return fmt.Errorf("len of bytes should be 32, but got: %d", l)
	}
	b32 := new([32]byte)
	copy(b32[:], b)
	w.wrap = append(w.wrap, *b32)
	return nil
}

// ------------------------------------------------------------------------------------------------------------------------------

// NewUint8ArrayWrap .
func NewUint8ArrayWrap() *Uint8ArrayWrap {
	return &Uint8ArrayWrap{}
}

// Uint8ArrayWrap wrap []uint8
type Uint8ArrayWrap struct {
	wrap []uint8
}

// AddOne .
func (w *Uint8ArrayWrap) AddOne(n uint8) {
	w.wrap = append(w.wrap, n)
}

// ------------------------------------------------------------------------------------------------------------------------------

// NewSizedByteArray 创建一个定长字节数组
func NewSizedByteArray(b []byte) *SizedByteArray {
	return &SizedByteArray{wrap: b}
}

// SizedByteArray 固定长度的字节数组(固定长度意味着修改时长度需要与原来的一致)
type SizedByteArray struct {
	wrap []byte
}

// Set 如果长度与原来的不一致则报错
func (w *SizedByteArray) Set(b []byte) error {
	if err := w.requireLen(len(b)); err != nil {
		return err
	}
	w.wrap = b
	return nil
}

// Get return byte array
func (w *SizedByteArray) Get() []byte {
	return w.wrap
}

func (w *SizedByteArray) requireLen(length int) error {
	if l := len(w.wrap); l != length {
		return fmt.Errorf("require len to be %d, actual %d", length, l)
	}
	return nil
}

func (w *SizedByteArray) bytes32() (b32 [32]byte, err error) {
	if err = w.requireLen(32); err != nil {
		return
	}
	copy(b32[:], w.wrap)
	return
}
