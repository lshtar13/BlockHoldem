package cli

import (
	"fmt"
	"log"

	"github.com/lshtar13/BlockHoldem/network"
	"github.com/lshtar13/BlockHoldem/wallet"
)

func (cli *CLI) startNode(minerAddr string) {
	if len(minerAddr) > 0 {
		if wallet.ValidateAddress(minerAddr) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddr)
		} else {
			log.Panic("Wrong miner address!")
		}
	}
	network.StartServer(cli.nodeID, minerAddr)
}
