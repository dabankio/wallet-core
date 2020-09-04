package internal

import (
	"encoding/json"

	"github.com/dabankio/gobbc"
	"github.com/pkg/errors"
)

const (
	SymbolBBC = "BBC"
	SymbolMKF = "MKF"
)

var knownSymbols = []string{
	SymbolBBC, SymbolMKF,
}

func isKnownSymbol(symbol string) error {
	for _, s := range knownSymbols {
		if s == symbol {
			return nil
		}
	}
	return errors.Errorf("Unknown symbol %s", symbol)
}
func SymbolSerializer(symbol string) gobbc.Serializer {
	switch symbol {
	case SymbolBBC:
		return gobbc.BBCSerializer
	case SymbolMKF:
		return gobbc.MKFSerializer
	default:
		return unknownSymbolSerializer(symbol)
	}
}

type unknownSymbolSerializer string

func (s unknownSymbolSerializer) Serialize(gobbc.RawTransaction) ([]byte, error) {
	return nil, errors.Errorf("unable to serialize, unknown symbol %s", s)
}
func (s unknownSymbolSerializer) Deserialize([]byte) (gobbc.RawTransaction, error) {
	return gobbc.RawTransaction{}, errors.Errorf("unable to deserialize, unknown symbol %s", s)
}

func DecodeSymbolTx(symbol, txData string) (string, error) {
	tx, err := gobbc.DecodeRawTransaction(SymbolSerializer(symbol), txData, false)
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(tx)
	if err != nil {
		return "", errors.Wrapf(err, "unable to marshal json, %#v", tx)
	}
	return string(b), nil
}
