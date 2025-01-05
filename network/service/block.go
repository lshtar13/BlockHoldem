package service

import (
	"fmt"

	"github.com/lshtar13/BlockHoldem/blockchain"
	"github.com/lshtar13/BlockHoldem/network/common"
	"github.com/lshtar13/BlockHoldem/network/node"
)

type block struct {
	AddrFrom string
	Block    []byte
}

func (blk *block) Handle(bc *blockchain.Blockchain) error {
	b, err := blockchain.DeserializeBlock(blk.Block)
	if err != nil {
		return err
	}
	fmt.Println("Recevied a new block!")
	err = bc.AddBlock(b)
	if err == nil {
		utxoSet := blockchain.UTXOSet{Blockchain: bc}
		utxoSet.Update(b)
		fmt.Printf("Added block %x\n", b.Hash)
	}

	return err
}

func sendBlock(addr string, b *blockchain.Block) error {
	data := block{node.MySelf(), b.Serialize()}
	payload := common.GobEncode(data)
	req := append(common.Command2Bytes("block"), payload...)

	return SendData(addr, req)
}

func NewBlock() *block {
	return &block{}
}
