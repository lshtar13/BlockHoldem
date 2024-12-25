package cli

import (
	"fmt"
	"log"

	"github.com/lshtar13/BlockHoldem/blockchain"
)

func (cli *CLI) createBlockchain(data string) {
	bc, err := blockchain.CreateBlockchain(data)
	if err != nil {
		log.Panic(err)
	}
	cli.bc = bc
	fmt.Println("Create new blockchain!")
}
