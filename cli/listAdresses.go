package cli

import (
	"fmt"
	"log"

	"github.com/lshtar13/blockchain/wallet"
)

func (cli *CLI) listAddresses() {
	wallets, err := wallet.NewWallets(cli.nodeID)
	if err != nil {
		log.Panic(err)
	}

	for _, address := range wallets.GetAddresses() {
		fmt.Println(address)
	}
}
