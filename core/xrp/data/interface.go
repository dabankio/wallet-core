package data

import (
	"io"
)

type Hashable interface {
	GetType() string
	Prefix() HashPrefix
	GetHash() *Hash256
}

type SignerAgent interface {
	Hashable
	InitialiseForSigning()
	InitialiseForMultiSigning()
	SigningPrefix() HashPrefix
	GetPublicKey() *PublicKey
	GetSignature() *VariableLength
	AddSignature(*Signer)
}

type Router interface {
	Hashable
	SuppressionId() Hash256
}

type Storer interface {
	Hashable
	Ledger() uint32
	NodeType() NodeType
	NodeId() *Hash256
}

type LedgerEntry interface {
	Storer
	GetLedgerEntryType() LedgerEntryType
	GetLedgerIndex() *Hash256
	GetPreviousTxnId() *Hash256
	Affects(Account) bool
}

type Transaction interface {
	SignerAgent
	GetTransactionType() TransactionType
	GetBase() *TxBase
	PathSet() PathSet
}

type Wire interface {
	Unmarshal(Reader) error
	Marshal(io.Writer) error
}
