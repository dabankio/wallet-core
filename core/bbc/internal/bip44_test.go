package internal

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtendedKey_Child(t *testing.T) {
	as := assert.New(t)

	// TODO 跟多轮次的随机种子

	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	as.Nil(err)

	mk, err := NewMaster(seed)
	as.Nil(err)

	for i := 0; i < 9999999; i++ {
		// for i := 0; i < 999; i++ {
		if i == 255 {
			continue
		}
		if i%100000 == 0 {
			fmt.Println(i)
		}
		child, err := mk.Child(uint32(i))
		as.Nil(err)
		as.Len(child.key, 32)
		// fmt.Println(len(child.key))
	}
}
