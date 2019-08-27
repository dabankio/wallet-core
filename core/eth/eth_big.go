package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// A BigInt represents a signed multi-precision integer.
type BigInt struct {
	bigint *big.Int
}

// NewBigInt allocates and returns a new BigInt set to x.
func NewBigInt(x int64) *BigInt {
	return &BigInt{big.NewInt(x)}
}

// GetBytes returns the absolute value of x as a big-endian byte slice.
func (bi *BigInt) GetBytes() []byte {
	return bi.bigint.Bytes()
}

// String returns the value of x as a formatted decimal string.
func (bi *BigInt) String() string {
	return bi.bigint.String()
}

// GetInt64 returns the int64 representation of x. If x cannot be represented in
// an int64, the result is undefined.
func (bi *BigInt) GetInt64() int64 {
	return bi.bigint.Int64()
}

// SetBytes interprets buf as the bytes of a big-endian unsigned integer and sets
// the big int to that value.
func (bi *BigInt) SetBytes(buf []byte) {
	bi.bigint.SetBytes(common.CopyBytes(buf))
}

// SetInt64 sets the big int to x.
func (bi *BigInt) SetInt64(x int64) {
	bi.bigint.SetInt64(x)
}

// Sign returns:
//
//	-1 if x <  0
//	 0 if x == 0
//	+1 if x >  0
//
func (bi *BigInt) Sign() int {
	return bi.bigint.Sign()
}

// SetString sets the big int to x.
//
// The string prefix determines the actual conversion base. A prefix of "0x" or
// "0X" selects base 16; the "0" prefix selects base 8, and a "0b" or "0B" prefix
// selects base 2. Otherwise the selected base is 10.
func (bi *BigInt) SetString(x string, base int) {
	bi.bigint.SetString(x, base)
}
