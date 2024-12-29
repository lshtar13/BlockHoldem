package cli

import (
	"fmt"

	"github.com/lshtar13/BlockHoldem/blockchain"
)

func (cli *CLI) send(from, to string, amount int) {
	bc, _ := blockchain.NewBlockchain(from)
	utxoset := blockchain.UTXOSet{Blockchain: bc}
	defer bc.DB.Close()

	tx := blockchain.NewUTXOTransaction(from, to, amount, &utxoset)
	cbTx := blockchain.NewCoinbaseTX(from, "")
	block := bc.MineBlock([]*blockchain.Transaction{cbTx, tx})
	utxoset.Update(block)

	fmt.Println("Success!")
}
