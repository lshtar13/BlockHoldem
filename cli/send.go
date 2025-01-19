package cli

import (
	"fmt"

	"github.com/lshtar13/blockchain/chain"
)

func (cli *CLI) send(from, to string, amount int) {
	bc, _ := chain.NewBlockchain(cli.nodeID)
	utxoset := chain.UTXOSet{Blockchain: bc}
	defer bc.DB.Close()

	tx := chain.NewUTXOTransaction(cli.nodeID, from, to, amount, &utxoset)
	cbTx := chain.NewCoinbaseTX(from, "")
	block := bc.MineBlock([]*chain.Transaction{cbTx, tx})
	utxoset.Update(block)

	fmt.Println("Success!")
}
