package blockchain

import (
	"bytes"

	wlt "github.com/lshtar13/BlockHoldem/wallet"
)

type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	Pubkey    []byte
}

func (in *TXInput) UseKey(pubKeyHash []byte) bool {
	lockingHash := wlt.HashPubkey(in.Pubkey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
