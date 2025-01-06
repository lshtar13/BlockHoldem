package cli

import (
	"fmt"
	"log"

	"github.com/lshtar13/BlockHoldem/base58"
	"github.com/lshtar13/BlockHoldem/blockchain"
	"github.com/lshtar13/BlockHoldem/wallet"
)

func (cli *CLI) getBalance(address string) {
	bc, _ := blockchain.NewBlockchain(cli.nodeID)
	UTXOSet := blockchain.UTXOSet{Blockchain: bc}
	defer bc.DB.Close()

	if !wallet.ValidateAddress(address) {
		log.Panic("Wrong address!")
	}

	balance := 0
	pubKeyHash := base58.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-wallet.AddressChecksumLen]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
