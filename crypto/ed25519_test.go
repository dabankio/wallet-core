package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEd25519(t *testing.T) {
	privk, pubk := "de4e126fd76aa11bdbbf2b861dcba0748deb376056ef26b99872a11813f9b402", "07f57f6448b28127be08bebcdcd31074d0bfc49b0bf2096bfa1bf6e708f9c94f"

	for i, msg := range [][]byte{
		[]byte("helloworld"),
		[]byte("de4e126fd76aa11bdbbf2b861dcba0748deb376056ef26b99872a11813f9b402:123"),
	} {

		priv, _ := HexDecodeThenReverse(privk)
		sig := Ed25519sign(priv, msg)

		pub, _ := HexDecodeThenReverse(pubk)
		require.True(t, Ed25519verify(pub, msg, sig), fmt.Sprintf("idx: %d", i))
	}

}
