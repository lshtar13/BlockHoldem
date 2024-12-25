package cli

import (
	"fmt"

	"github.com/lshtar13/BlockHoldem/blockchain"
)

func (cli *CLI) send(from, to string, amount int) {
	bc, _ := blockchain.NewBlockchain(from)
	defer bc.DB.Close()

	tx := blockchain.NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")
}
