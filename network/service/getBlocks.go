package service

import (
	"fmt"

	"github.com/lshtar13/blockchain/blockchain"
	"github.com/lshtar13/blockchain/network/common"
	"github.com/lshtar13/blockchain/network/node"
)

type getBlocks struct {
	AddrFrom string
}

func (gb *getBlocks) Handle(bc *blockchain.Blockchain) error {
	blocks := bc.GetBlockHashes()
	err := sendInv(gb.AddrFrom, "block", blocks)
	if err != nil {
		return err
	}

	return nil
}

func sendGetBlocks(addr string) error {
	fmt.Println("  Send GetBlocks...")
	payload := common.GobEncode(getBlocks{node.MySelf()})
	req := append(common.Command2Bytes("getBlocks"), payload...)

	return SendData(addr, req)
}

func NewGetBlocks() *getBlocks {
	return &getBlocks{}
}
