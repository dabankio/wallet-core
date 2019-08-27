package eth

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func BenchmarkSimpleMultiSigExecuteSignResult_ToHex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b65 := new([65]byte)
		for j := 0; j < 65; j++ {
			b65[j] = byte(rand.Intn(255))
		}
		r := SimpleMultiSigExecuteSignResult{
			R: &SizedByteArray{wrap: b65[:32]},
			S: &SizedByteArray{wrap: b65[32:64]},
			V: int8(b65[64]),
		}
		he := r.ToHex()
		r1, err := NewSimpleMultiSigExecuteSignResultFromHex(he)
		if err != nil {
			b.Fatal("Decode failed", err)
		}
		if !reflect.DeepEqual(r, *r1) {
			fmt.Println(r)
			fmt.Println(r1)
			b.Fatal("Not Equal")
		}
	}
}
