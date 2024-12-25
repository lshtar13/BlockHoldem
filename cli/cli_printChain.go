package cli

import (
	"fmt"
	"strconv"

	"github.com/lshtar13/BlockHoldem/blockchain"
)

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

	for {
		block, _ := bci.Next()
		fmt.Printf("Prev. hash : %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
