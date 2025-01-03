package network

import (
	"github.com/lshtar13/BlockHoldem/blockchain"
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
	payload := gobEncode(getBlocks{nodeAddr})
	req := append(command2Bytes("getBlocks"), payload...)

	return sendData(addr, req)
}

func NewGetBlocks() *getBlocks {
	return &getBlocks{}
}
