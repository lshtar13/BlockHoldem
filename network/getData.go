package network

import "github.com/lshtar13/BlockHoldem/blockchain"

type getData struct {
	AddrFrom string
	Type     string
	ID       []byte
}

func (g *getData) Handle(bc *blockchain.Blockchain) error {
	//herere

}

func sendGetData(addr string, sort string, id []byte) error {
	payload := gobEncode(getBlocks{nodeAddr})
	req := append(command2Bytes("getData"), payload...)

	return sendData(addr, req)
}

func NewGetData() *getData {
	return &getData{}
}
