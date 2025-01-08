package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"strings"

	"github.com/lshtar13/BlockHoldem/base58"
	wlt "github.com/lshtar13/BlockHoldem/wallet"
)

type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

func (out *TXOutput) Print(nindent int) {
	indent := strings.Repeat(" ", nindent)
	fmt.Printf("%sValue : %d\n", indent, out.Value)
	fmt.Printf("%sPubKeyHash : %x\n", indent, out.PubKeyHash)
}

func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := base58.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-wlt.AddressChecksumLen]
	out.PubKeyHash = pubKeyHash
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, pubKeyHash)
}

func NewTXOutput(amount int, to string) TXOutput {
	txoutput := TXOutput{amount, nil}
	txoutput.Lock([]byte(to))

	return txoutput
}

type TXOutputs struct {
	Outputs []TXOutput
}

func (outs *TXOutputs) Serialize() []byte {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

func DeserializeOutputs(data []byte) *TXOutputs {
	outs := TXOutputs{}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&outs)
	if err != nil {
		log.Panic(err)
	}

	return &outs
}
