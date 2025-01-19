package cli

import (
	"log"

	"github.com/lshtar13/blockchain/chain"
)

func (cli *CLI) reindexUTXO() {
	bc, err := chain.NewBlockchain(cli.nodeID)
	if err != nil {
		log.Panic(err)
	}

	UTXOset := chain.UTXOSet{Blockchain: bc}
	UTXOset.Reindex()
}
