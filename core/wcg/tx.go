package wcg

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
)

type Tx struct {
	Hex                   string   `json:"-"`
	Ctxt                  string   `json:"-"`
	Type                  byte     `json:"-"`
	version               byte     `json:"-"`
	subType               byte     `json:"-"`
	Timestamp             uint32   `json:"timestamp"`
	Deadline              uint16   `json:"deadline"`
	pubKey                [32]byte `json:"-"`
	recipientId           [8]byte  `json:"-"`
	Amount                uint64   `json:"amountNQT"`
	Fee                   uint64   `json:"feeNQT"`
	fullHash              [32]byte `json:"-"`
	sign                  [64]byte `json:"-"`
	flags                 uint32   `json:"-"`
	BlockHeight           uint32   `json:"ecBlockHeight"`
	BlockId               uint64   `json:"ecBlockId"`
	FinishHeight          uint32   `json:"phasingFinishHeight"`
	votingModelCode       byte     `json:"-"`
	Quorum                uint64   `json:"phasingQuorum"`
	minBalance            uint64   `json:"-"`
	whiteListLength       byte     `json:"-"`
	WhiteList             []string `json:"phasingWhitelist"`
	HoldingID             uint64   `json:"phasingHolding"`
	minBalanceModelCode   byte     `json:"-"`
	linkedFullHashsLength byte     `json:"-"`
	linkFullHashs         []byte   `json:"-"`
	hashedSecretLength    byte     `json:"-"`
	hashSecrets           []byte   `json:"-"`
	algorithm             byte     `json:"-"`
	SenderPublicKey       string   `json:"senderPublicKey"`
	TypeInt               int      `json:"type"`
	VersionInt            int      `json:"version"`
	SubtypeInt            int      `json:"subtype"`
	RecipientRS           string   `json:"recipientRS"`
	Recipient             string   `json:"recipient"`
	Signature             string   `json:"signature"`
	VersionPhasing        int      `json:"versionPhasing"`
	CTXT
	ApproveTx
}

type CTXT struct {
	Sender   string `json:"sender"`
	SenderRS string `json:"senderRS"`
	PrivKey  string `json:"privatekey"`
	PubKey   string `json:"pubkey"`
}

type ApproveTx struct {
	TxHashSize           byte     `json:"-"`
	TxHashs              []string `json:"transactionFullHashes"`
	RevealedSecretLength int32    `json:"-"`
	RevealedSecrets      []byte   `json:"-"`
}

func (t Tx) GetHex() string {
	return t.Hex
}

func (t Tx) GetSender() string {
	return t.SenderRS
}

func (t Tx) String() string {
	if j, err := json.Marshal(t); err != nil {
		return ""
	} else {
		return string(j)
	}
}

func (t *Tx) Parse() error {
	if t.Hex == "" {
		return nil
	}

	if b, err := hex.DecodeString(t.Hex); err != nil {
		return err
	} else {
		buff := bytes.NewBuffer(b)

		if err := binary.Read(buff, binary.LittleEndian, &t.Type); err != nil {
			return err
		}

		var ver byte
		if err := binary.Read(buff, binary.LittleEndian, &ver); err != nil {
			return err
		}
		t.version = ver >> 4
		t.subType = ver & 0x0f

		if err := binary.Read(buff, binary.LittleEndian, &t.Timestamp); err != nil {
			return err
		}

		if err := binary.Read(buff, binary.LittleEndian, &t.Deadline); err != nil {
			return err
		}

		if err := binary.Read(buff, binary.LittleEndian, &t.pubKey); err != nil {
			return err
		}
		t.SenderPublicKey = hex.EncodeToString(t.pubKey[:])

		if err := binary.Read(buff, binary.LittleEndian, &t.recipientId); err != nil {
			return err
		}

		w := WCG{}

		if accountid, err := w.GetAccountIdByRecipient(t.recipientId[:]); err != nil {
			return err
		} else {
			t.Recipient = accountid

			if recipientRs, err := w.GetAccountById(accountid); err == nil {
				t.RecipientRS = recipientRs
			}
		}

		if err := binary.Read(buff, binary.LittleEndian, &t.Amount); err != nil {
			return err
		}

		if err := binary.Read(buff, binary.LittleEndian, &t.Fee); err != nil {
			return err
		}

		if err := binary.Read(buff, binary.LittleEndian, &t.fullHash); err != nil {
			return err
		}

		if err := binary.Read(buff, binary.LittleEndian, &t.sign); err != nil {
			return err
		}
		t.Signature = hex.EncodeToString(t.sign[:])

		if int(t.version) > 0 {
			if err := binary.Read(buff, binary.LittleEndian, &t.flags); err != nil {
				return err
			}

			if err := binary.Read(buff, binary.LittleEndian, &t.BlockHeight); err != nil {
				return err
			}

			if err := binary.Read(buff, binary.LittleEndian, &t.BlockId); err != nil {
				return err
			}
		}

		if int(t.version) > 0 {
			var version byte
			if err := binary.Read(buff, binary.LittleEndian, &version); err != nil {
				return err
			}
			t.VersionPhasing = int(version)
		}

		if t.version == 0x01 && t.subType == 0x09 {
			//	approveTransaction
			if err := t.parseApproveTx(buff); err != nil {
				return err
			}
		} else if t.version == 0x01 && t.subType == 0x00 {
			// sendMoney
			if err := t.parseNewTx(buff); err != nil {
				return err
			}
		}
	}

	if t.Ctxt != "" {
		var ctxt CTXT
		if err := json.Unmarshal([]byte(t.Ctxt), &ctxt); err == nil {
			t.SenderRS = ctxt.Sender
			t.PrivKey = ctxt.PrivKey
		}
	}

	return nil
}

func (t *Tx) parseNewTx(buff *bytes.Buffer) error {

	if err := binary.Read(buff, binary.LittleEndian, &t.FinishHeight); err != nil {
		return err
	}

	if err := binary.Read(buff, binary.LittleEndian, &t.votingModelCode); err != nil {
		return err
	}

	if err := binary.Read(buff, binary.LittleEndian, &t.Quorum); err != nil {
		return err
	}

	if err := binary.Read(buff, binary.LittleEndian, &t.minBalance); err != nil {
		return err
	}

	if err := binary.Read(buff, binary.LittleEndian, &t.whiteListLength); err != nil {
		return err
	}

	wcg := WCG{}
	for i := 0; i < int(t.whiteListLength); i++ {
		account := make([]byte, 8)
		if err := binary.Read(buff, binary.LittleEndian, &account); err != nil {
			return err
		}
		accountid, _ := wcg.GetAccountIdByRecipient(account)
		t.WhiteList = append(t.WhiteList, accountid)
	}

	if err := binary.Read(buff, binary.LittleEndian, &t.HoldingID); err != nil {
		return err
	}

	if err := binary.Read(buff, binary.LittleEndian, &t.minBalanceModelCode); err != nil {
		return err
	}

	if err := binary.Read(buff, binary.LittleEndian, &t.linkedFullHashsLength); err != nil {
		return err
	}

	for i := 0; i < int(t.linkedFullHashsLength); i++ {
		var linkedFullHash byte
		if err := binary.Read(buff, binary.LittleEndian, &linkedFullHash); err != nil {
			return err
		}
		t.linkFullHashs = append(t.linkFullHashs, linkedFullHash)
	}

	if err := binary.Read(buff, binary.LittleEndian, &t.hashedSecretLength); err != nil {
		return err
	}
	for i := 0; i < int(t.hashedSecretLength); i++ {
		var hashSecret byte
		if err := binary.Read(buff, binary.LittleEndian, &hashSecret); err != nil {
			return err
		}
		t.hashSecrets = append(t.hashSecrets, hashSecret)
	}

	if err := binary.Read(buff, binary.LittleEndian, &t.algorithm); err != nil {
		return err
	}
	return nil
}

func (t *Tx) parseApproveTx(buff *bytes.Buffer) error {
	if err := binary.Read(buff, binary.LittleEndian, &t.TxHashSize); err != nil {
		return err
	}

	for i := 0; i < int(t.TxHashSize); i++ {
		hash := make([]byte, 32)
		if err := binary.Read(buff, binary.LittleEndian, &hash); err != nil {
			return err
		}
		t.TxHashs = append(t.TxHashs, hex.EncodeToString(hash))
	}

	if err := binary.Read(buff, binary.LittleEndian, &t.RevealedSecretLength); err != nil {
		return err
	}

	for i := 0; i < int(t.RevealedSecretLength); i++ {
		var revealedSecret byte
		if err := binary.Read(buff, binary.LittleEndian, &revealedSecret); err != nil {
			return err
		}
		t.RevealedSecrets = append(t.RevealedSecrets, revealedSecret)
	}

	return nil
}
