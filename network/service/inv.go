package service

import (
	"fmt"

	"github.com/lshtar13/BlockHoldem/blockchain"
	"github.com/lshtar13/BlockHoldem/network/common"
	"github.com/lshtar13/BlockHoldem/network/node"
)

type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

func (iv *inv) Handle(bc *blockchain.Blockchain) error {
	fmt.Printf("Recevied inventory with %d %s\n", len(iv.Items), iv.Type)
	switch iv.Type {
	case "block":
		for _, item := range iv.Items {
			PutTxTransit(item, iv.AddrFrom)
		}
	case "tx":
		for _, txID := range iv.Items {
			if GetTx(string(txID)).ID == nil {
				sendGetData(iv.AddrFrom, "tx", txID)
			}
		}
	}

	return nil
}

func sendInv(addr string, sort string, items [][]byte) error {
	payload := common.GobEncode(inv{node.MySelf(), sort, items})
	req := append(common.Command2Bytes("inv"), payload...)

	return SendData(addr, req)
}

func NewInv() *inv {
	return &inv{}
}
