package service

import (
	"github.com/lshtar13/BlockHoldem/blockchain"
	"github.com/lshtar13/BlockHoldem/network/common"
	"github.com/lshtar13/BlockHoldem/network/node"
)

type version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

func (ver *version) Handle(bc *blockchain.Blockchain) error {
	var err error
	myBestHeight := bc.GetBestHeight()
	hisBestHeight := ver.BestHeight

	if myBestHeight < hisBestHeight {
		err = sendGetBlocks(ver.AddrFrom)
	} else if myBestHeight > hisBestHeight {
		err = sendVersion(ver.AddrFrom, bc)
	}

	if err != nil {
		return err
	}

	node.Insert(ver.AddrFrom)

	return nil
}

func sendVersion(addr string, bc *blockchain.Blockchain) error {
	bestHeight := bc.GetBestHeight()
	payload := common.GobEncode(version{node.Version(), bestHeight, node.MySelf()})

	request := append(common.Command2Bytes("version"), payload...)

	return SendData(addr, request)
}

func NewVersion() *version {
	return &version{}
}
