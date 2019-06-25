package bip44

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCoinType(t *testing.T) {
	coinType, err := GetCoinType("ETH")
	assert.NoError(t, err)
	assert.EqualValues(t, coinType, 0x3c)
	t.Log(coinType)
}
