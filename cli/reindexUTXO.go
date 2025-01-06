package cli

import (
	"log"

	"github.com/lshtar13/BlockHoldem/blockchain"
)

func (cli *CLI) reindexUTXO() {
	bc, err := blockchain.NewBlockchain(cli.nodeID)
	if err != nil {
		log.Panic(err)
	}

	UTXOset := blockchain.UTXOSet{Blockchain: bc}
	UTXOset.Reindex()
}
