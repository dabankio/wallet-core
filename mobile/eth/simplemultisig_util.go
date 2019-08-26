package geth

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lomocoin/wallet-core/core/eth"
)

// SimpleMultiSigExecuteSignResult .
type SimpleMultiSigExecuteSignResult struct {
	R, S *SizedByteArray
	V    int8
}

// ToHex 转为hex string方便传输
func (r *SimpleMultiSigExecuteSignResult) ToHex() string {
	bytes := r.R.wrap
	bytes = append(bytes, r.S.wrap...)
	bytes = append(bytes, byte(r.V))
	return hexutil.Encode(bytes)
}

// NewSimpleMultiSigExecuteSignResultFromHex decode from hex string
func NewSimpleMultiSigExecuteSignResultFromHex(h string) (*SimpleMultiSigExecuteSignResult, error) {
	if l := len(h); l != 132 {
		return nil, fmt.Errorf("length of hexed result should be 132,start with 0x, got: %d", l)
	}
	b, err := hexutil.Decode(h)
	if err != nil {
		return nil, err
	}
	ret := new(SimpleMultiSigExecuteSignResult)
	ret.R = &SizedByteArray{wrap: b[:32]}
	ret.S = &SizedByteArray{wrap: b[32:64]}
	ret.V = int8(b[64])
	return ret, nil
}

// UtilSimpleMultiSigExecuteSign 签名简单多签执行数据
func UtilSimpleMultiSigExecuteSign(chainID int64, signerPrivkHex string, hexedMultisigAddr, hexedDestinationAddr, hexedExecutor string, nonce int64, value, gasLimit *BigInt, data []byte) (*SimpleMultiSigExecuteSignResult, error) {
	v, r, s, err := eth.SimpleMultiSigExecuteSign(chainID, signerPrivkHex, hexedMultisigAddr, hexedDestinationAddr, hexedExecutor, uint64(nonce), value.bigint, gasLimit.bigint, data)
	if err != nil {
		return nil, err
	}
	return &SimpleMultiSigExecuteSignResult{
		V: int8(v),
		R: &SizedByteArray{r[:]},
		S: &SizedByteArray{s[:]},
	}, nil
}
