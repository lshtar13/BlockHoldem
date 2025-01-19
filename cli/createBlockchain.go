package cli

import (
	"fmt"
	"log"

	"github.com/lshtar13/blockchain/chain"
)

func (cli *CLI) createBlockchain(address string) {
	bc, err := chain.CreateBlockchain(address, cli.nodeID)
	defer bc.DB.Close()
	if err != nil {
		log.Panic(err)
	}

	UTXOSet := chain.UTXOSet{Blockchain: bc}
	UTXOSet.Reindex()

	cli.bc = bc
	fmt.Println("Create new blockchain!")
}
