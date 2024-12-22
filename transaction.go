package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

type TXOutput struct {
	Value        int
	ScriptPubkey string
}

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

func (in *TXInput) CanUnlockOutputWith(scriptPubkey string) bool {
	return in.ScriptSig == scriptPubkey
}

func (out *TXOutput) CanBeUnlockedWith(scriptSig string) bool {
	return out.ScriptPubkey == scriptSig
}

func (tx *Transaction) SetID() error {
	var encoded bytes.Buffer
	var hash [32]byte
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]

	return err
}

// to - address, data - signature
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}
