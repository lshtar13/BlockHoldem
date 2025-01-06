package service

import (
	"github.com/lshtar13/BlockHoldem/blockchain"
	"github.com/lshtar13/BlockHoldem/network/common"
	"github.com/lshtar13/BlockHoldem/network/node"
)

type tx struct {
	AddFrom     string
	Transaction []byte
}

func (t *tx) Handle(bc *blockchain.Blockchain) error {
	transaction, err := blockchain.DeserializeTransaction(t.Transaction)
	if err != nil {
		return err
	}

	AddTx(transaction)
	if node.MySelf() == node.CentralNode() {
		BroadCastInv("tx", []string{node.MySelf(), t.AddFrom}, [][]byte{transaction.ID})
	}
	return nil
}

func sendTx(addr string, t *blockchain.Transaction) error {
	data := tx{node.MySelf(), t.Serialize()}
	payload := common.GobEncode(data)
	req := append(common.Command2Bytes("tx"), payload...)

	return SendData(addr, req)
}

func NewTx() *tx {
	return &tx{}
}
