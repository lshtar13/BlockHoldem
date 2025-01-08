package cli

import (
	"fmt"
	"log"

	"github.com/lshtar13/BlockHoldem/blockchain"
)

func (cli *CLI) printChain() {
	bc, err := blockchain.NewBlockchain(cli.nodeID)
	if err != nil {
		log.Panic(err)
	}
	cli.bc = bc
	bci := cli.bc.Iterator()

	for {
		block, _ := bci.Next()
		block.Print()
		fmt.Println()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
