package blockchain

import (
	"bytes"
	"fmt"
	"strings"

	wlt "github.com/lshtar13/BlockHoldem/wallet"
)

type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	Pubkey    []byte
}

func (in *TXInput) Print(nindent int) {
	indent := strings.Repeat(" ", nindent)
	fmt.Printf("%sTxid : %x\n", indent, in.Txid)
	fmt.Printf("%sVout : %d\n", indent, in.Vout)
}

func (in *TXInput) UseKey(pubKeyHash []byte) bool {
	lockingHash := wlt.HashPubkey(in.Pubkey)

	return bytes.Equal(lockingHash, pubKeyHash)
}
