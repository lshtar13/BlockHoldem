package blockchain

import (
	"bytes"

	"github.com/lshtar13/BlockHoldem/base58"
	wlt "github.com/lshtar13/BlockHoldem/wallet"
)

type TXOutput struct {
	Value      int
	PubKeyHash []byte
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
