package network

import (
	"fmt"

	"github.com/lshtar13/BlockHoldem/blockchain"
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
			blockTransitChan <- transit{item, iv.AddrFrom}
		}
	case "tx":
		//
	}

	return nil
}

func sendInv(addr string, sort string, items [][]byte) error {
	payload := gobEncode(inv{nodeAddr, sort, items})
	req := append(command2Bytes("inv"), payload...)

	return sendData(addr, req)
}

func NewInv() *inv {
	return &inv{}
}
