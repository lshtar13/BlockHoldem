package cli

import (
	"fmt"
	"log"

	"github.com/lshtar13/BlockHoldem/blockchain"
)

func (cli *CLI) createBlockchain(data string) {
	bc, err := blockchain.CreateBlockchain(data)
	defer bc.DB.Close()
	if err != nil {
		log.Panic(err)
	}

	UTXOSet := blockchain.UTXOSet{Blockchain: bc}
	UTXOSet.Reindex()

	cli.bc = bc
	fmt.Println("Create new blockchain!")
}
