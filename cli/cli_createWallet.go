package cli

import (
	"fmt"

	"github.com/lshtar13/BlockHoldem/wallet"
)

func (cli *CLI) createWallet() {
	wallets, _ := wallet.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Println("Your new address: ", address)
}
