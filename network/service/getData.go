package service

import (
	"encoding/hex"

	"github.com/lshtar13/BlockHoldem/blockchain"
	"github.com/lshtar13/BlockHoldem/network/common"
	"github.com/lshtar13/BlockHoldem/network/node"
)

type getData struct {
	AddrFrom string
	Type     string
	ID       []byte
}

func (gd *getData) Handle(bc *blockchain.Blockchain) error {
	switch gd.Type {
	case "block":
		block, err := bc.GetBlock([]byte(gd.ID))
		if err != nil {
			sendBlock(gd.AddrFrom, block)
		}
	case "tx":
		txID := hex.EncodeToString(gd.ID)
		tx := GetTx(txID)

		sendTx(gd.AddrFrom, tx)
	}

	return nil
}

func sendGetData(addr string, sort string, id []byte) error {
	payload := common.GobEncode(getBlocks{node.MySelf()})
	req := append(common.Command2Bytes("getData"), payload...)

	return SendData(addr, req)
}

func NewGetData() *getData {
	return &getData{}
}
