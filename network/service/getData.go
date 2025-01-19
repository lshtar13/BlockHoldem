package service

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/lshtar13/blockchain/blockchain"
	"github.com/lshtar13/blockchain/network/common"
	"github.com/lshtar13/blockchain/network/node"
)

type getData struct {
	AddrFrom string
	Type     string
	ID       []byte
}

func (gd *getData) Handle(bc *blockchain.Blockchain) error {
	fmt.Printf("  Handling getData %s\n", gd.Type)
	switch gd.Type {
	case "block":
		block, err := bc.GetBlock([]byte(gd.ID))
		if err != nil {
			log.Panic(err)
		}
		sendBlock(gd.AddrFrom, block)
	case "tx":
		txID := hex.EncodeToString(gd.ID)
		tx := GetTx(txID)

		sendTx(gd.AddrFrom, tx)
	}

	return nil
}

func sendGetData(addr string, sort string, id []byte) error {
	fmt.Println("  Sedn GetData...")
	payload := common.GobEncode(getData{node.MySelf(), sort, id})
	req := append(common.Command2Bytes("getData"), payload...)

	return SendData(addr, req)
}

func NewGetData() *getData {
	return &getData{}
}
