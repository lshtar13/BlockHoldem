package cli

import (
	"fmt"

	"github.com/lshtar13/BlockHoldem/wallet"
)

func (cli *CLI) createWallet() {
	wallets, _ := wallet.NewWallets(cli.nodeID)
	address := wallets.CreateWallet()
	wallets.SaveToFile(cli.nodeID)

	fmt.Println("Your new address: ", address)
}
