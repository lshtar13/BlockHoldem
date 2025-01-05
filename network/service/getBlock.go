package service

import (
	"github.com/lshtar13/BlockHoldem/blockchain"
	"github.com/lshtar13/BlockHoldem/network/common"
	"github.com/lshtar13/BlockHoldem/network/node"
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
	payload := common.GobEncode(getBlocks{node.MySelf()})
	req := append(common.Command2Bytes("getBlocks"), payload...)

	return SendData(addr, req)
}

func NewGetBlocks() *getBlocks {
	return &getBlocks{}
}
