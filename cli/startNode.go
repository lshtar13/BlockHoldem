package cli

import (
	"fmt"
	"log"

	"github.com/lshtar13/blockchain/node"
	"github.com/lshtar13/blockchain/wallet"
)

func (cli *CLI) startNode(port int, minerAddr string, mineCap int) {
	isMiner := false
	if len(minerAddr) > 0 {
		isMiner = true
		if wallet.ValidateAddress(minerAddr) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddr)
		} else {
			log.Panic("Wrong miner address!")
		}
	}

	if node, err := node.NewNode(cli.nodeID, minerAddr, port); err == nil {
		if err = node.Start(isMiner, mineCap); err != nil {
			log.Panic("err while starting node:", err)
		}
	} else {
		log.Panic("err while making node", err)
	}
}
