package wcg

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/core"
	"github.com/pkg/errors"
)

const symbol = "WCG"

var (
	gexp     = []int{1, 2, 4, 8, 16, 5, 10, 20, 13, 26, 17, 7, 14, 28, 29, 31, 27, 19, 3, 6, 12, 24, 21, 15, 30, 25, 23, 11, 22, 9, 18, 1}
	glog     = []int{0, 0, 1, 18, 2, 5, 19, 11, 3, 29, 6, 27, 20, 8, 12, 23, 4, 10, 30, 17, 7, 22, 28, 26, 21, 25, 9, 16, 13, 14, 24, 15}
	cwmap    = []int{3, 2, 1, 0, 7, 6, 5, 4, 13, 14, 15, 16, 12, 8, 9, 10, 11}
	alphabet = []string{"2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
)

func New(seed []byte) (c *WCG, err error) {
	c = new(WCG)
	c.Symbol = symbol
	c.DerivationPath, err = bip44.GetCoinDerivationPath(symbol)
	if err != nil {
		err = errors.Wrap(err, "bip44.GetCoinDerivationPath err:")
		return
	}
	c.MasterKey, err = hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return
	}
	return
}

type WCG struct {
	core.CoinInfo
}

// create wcg address by path
func (w *WCG) DeriveAddress() (address string, err error) {
	publicKey, err := w.DerivePublicKey()
	if err != nil {
		return
	}

	accountid, err := w.GetAccountIdByPk(publicKey)
	if err != nil {
		return
	}

	address, err = w.GetAccountById(accountid)
	if err != nil {
		return
	}

	return
}

// get wcg`s accountid by publickey
func (w *WCG) GetAccountIdByPk(publicKey string) (accountid string, err error) {

	hexByte, err := hex.DecodeString(publicKey)
	if err != nil {
		return
	}
	h := sha256.New()
	h.Write(hexByte)
	ret := h.Sum(nil)[0:8]
	accountid, err = w.GetAccountIdByRecipient(ret)

	return
}

// get wcg`s accountid by recipientid
func (w *WCG) GetAccountIdByRecipient(recipient []byte) (accountid string, err error) {
	if len(recipient) < 8 {
		err = errors.New("recipient error")
		return
	}
	value := big.NewInt(0)
	tmp1 := big.NewInt(0)
	tmp2 := big.NewInt(0)

	for i := 7; i >= 0; i-- {
		tmp1.Mul(value, big.NewInt(256))
		tmp2.Add(tmp1, big.NewInt(int64(recipient[i])))
		value = tmp2
	}

	accountid = value.String()
	return
}

// get wcg`s address by accountid
func (w *WCG) GetAccountById(accountid string) (address string, err error) {
	accountLen := len(accountid)
	if accountLen == 20 && accountid[0] != '1' {
		err = errors.New("account error")
		return
	}

	inp := make([]int, accountLen)
	for i := range inp {
		inp[i], _ = strconv.Atoi(string(accountid[i]))
	}

	var out []int

	for {
		divide := 0
		newlen := 0

		for i := 0; i < accountLen; i++ {
			divide = divide*10 + inp[i]

			if divide >= 32 {
				inp[newlen] = divide >> 5
				newlen++
				divide &= 31
			} else if newlen > 0 {
				inp[newlen] = 0
				newlen++
			}
		}

		accountLen = newlen
		out = append(out, divide)

		if newlen <= 0 {
			break
		}
	}

	codeword := []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	var pos = len(out)
	for i := 0; i < 13; i++ {
		pos--
		if pos >= 0 {
			codeword[i] = out[i]
		} else {
			codeword[i] = 0
		}
	}

	w.encode(codeword)

	address = fmt.Sprintf("%s-", w.Symbol)
	for i := 0; i < 17; i++ {
		address += alphabet[codeword[cwmap[i]]]

		if i&3 == 3 && i < 13 {
			address += "-"
		}
	}
	return
}

func (w *WCG) derivePublicKey(privateKey string) (publicKey string) {
	h := sha256.New()
	h.Write([]byte(privateKey))
	secret := h.Sum(nil)
	var p, s, k [64]byte

	for i := range secret {
		k[i] = secret[i]
	}

	Keygen(&p, &s, &k)
	publicKey = hex.EncodeToString(p[0:32])
	return
}

// create wcg`s publickey by path
func (w *WCG) DerivePublicKey() (publicKey string, err error) {
	privateKey, err := w.DerivePrivateKey()
	if err != nil {
		return
	}
	return w.derivePublicKey(privateKey), nil
}

// create wcg`s privatekey by path
func (w *WCG) DerivePrivateKey() (privateKey string, err error) {
	childKey := w.MasterKey
	childKey, err = w.MasterKey.Child(w.DerivationPath[len(w.DerivationPath)-1])
	if err != nil {
		return
	}

	ECPrivateKey, err := childKey.ECPrivKey()
	if err != nil {
		return
	}
	buff := bytes.NewBuffer(crypto.FromECDSA(ECPrivateKey.ToECDSA()))
	var tmp []string
	words := WcgWords{}
	l := uint64(words.Length() / 8)
	for i := 0; i < 4; i++ {
		var x uint64
		_ = binary.Read(buff, binary.LittleEndian, &x)
		w1 := x % l
		w2 := (((x / l) >> 0) + w1) % l
		w3 := (((((x / l) >> 0) / l) >> 0) + w2) % l
		tmp = append(tmp, words.FindIndex(int(w1)), words.FindIndex(int(w2)), words.FindIndex(int(w3)))
	}
	privateKey = strings.Join(tmp, " ")

	return
}

func (w *WCG) DecodeTx(msg string) (tx string, err error) {
	t := Tx{Hex: msg}
	if err = t.Parse(); err != nil {
		return
	}

	tx = t.String()

	return
}

func (w *WCG) Sign(msg, privateKey string) (sig string, err error) {
	h := sha256.New()
	h.Write([]byte(privateKey))
	digest := h.Sum(nil)

	var p, s, k [64]byte
	copy(k[:], digest)
	Keygen(&p, &s, &k)

	h.Reset()
	mByte, err := hex.DecodeString(msg)
	if err != nil {
		return
	}
	h.Write(mByte)
	m := h.Sum(nil)

	h.Reset()
	h.Write(m)
	h.Write(s[0:32])
	x := h.Sum(nil)

	var yp, ys, yk [64]byte
	copy(yk[:], x)
	Keygen(&yp, &ys, &yk)

	h.Reset()
	h.Write(m)
	h.Write(yp[0:32])
	h1 := h.Sum(nil)

	var v, signh, signs [64]byte
	copy(signh[:], h1)
	copy(signs[:], s[:])
	Sign(&v, &signh, &yk, &signs)

	copy(v[:], append(v[0:32], h1...))

	txByte := []byte(msg)

	signByte := append(txByte[0:192], []byte(hex.EncodeToString(v[0:64]))...)
	signByte = append(signByte[:], txByte[320:]...)

	sig = string(signByte)
	return
}

func (w *WCG) VerifySignature(pubKey, msg, signature string) error {
	// TODO
	return core.ErrThisFeatureIsNotSupported
}

func (w *WCG) encode(codeword []int) {
	p := []int{0, 0, 0, 0}

	for i := 12; i >= 0; i-- {
		fb := codeword[i] ^ p[3]
		p[3] = p[2] ^ w.gmult(30, fb)
		p[2] = p[1] ^ w.gmult(6, fb)
		p[1] = p[0] ^ w.gmult(9, fb)
		p[0] = w.gmult(17, fb)
	}

	codeword[13] = p[0]
	codeword[14] = p[1]
	codeword[15] = p[2]
	codeword[16] = p[3]
}

func (w *WCG) gmult(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}

	var idx = (glog[a] + glog[b]) % 31

	return gexp[idx]
}
