package network

import "github.com/lshtar13/BlockHoldem/blockchain"

type version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

func sendVersion(addr string, bc *blockchain.Blockchain) error {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(version{nodeVersion, bestHeight, nodeAddr})

	request := append(command2Bytes("version"), payload...)

	return sendData(addr, request)
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

	if !nodeIsKnown(ver.AddrFrom) {
		knownNodes = append(knownNodes, ver.AddrFrom)
	}

	return nil
}

func NewVersion() *version {
	return &version{}
}
