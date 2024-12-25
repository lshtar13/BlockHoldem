package main

import "bytes"

type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	Pubkey    []byte
}

func (in *TXInput) UseKey(pubKeyHash []byte) bool {
	lockingHash := HashPubkey(in.Pubkey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
