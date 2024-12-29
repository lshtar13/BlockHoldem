package cli

import (
	"fmt"

	"github.com/lshtar13/BlockHoldem/base58"
	"github.com/lshtar13/BlockHoldem/blockchain"
	"github.com/lshtar13/BlockHoldem/wallet"
)

func (cli *CLI) getBalance(address string) {
	bc, _ := blockchain.NewBlockchain(address)
	UTXOSet := blockchain.UTXOSet{Blockchain: bc}
	defer bc.DB.Close()

	// validation ...

	balance := 0
	pubKeyHash := base58.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-wallet.AddressChecksumLen]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
