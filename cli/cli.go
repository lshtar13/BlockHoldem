package cli

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/lshtar13/BlockHoldem/blockchain"
)

type CLI struct {
	bc *blockchain.Blockchain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createBlockchain -address ADDRESS - create new blockchain")
	fmt.Println("  printChain - print all the blocks of the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()
	createBlockchainCmd := flag.NewFlagSet("createBlockchain", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createWallet", flag.ExitOnError)

	createAddress := createBlockchainCmd.String("address", "", "address")
	getBalanceAddress := getBalanceCmd.String("address", "", "address")
	sendFrom := sendCmd.String("from", "", "from")
	sendTo := sendCmd.String("to", "", "to")
	sendAmount := sendCmd.String("amount", "", "amount")

	switch os.Args[1] {
	case "createBlockchain":
		createBlockchainCmd.Parse(os.Args[2:])
	case "printChain":
		printChainCmd.Parse(os.Args[2:])
	case "getBalance":
		getBalanceCmd.Parse(os.Args[2:])
	case "send":
		sendCmd.Parse(os.Args[2:])
	case "createWallet":
		createWalletCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createBlockchainCmd.Parsed() {
		if *createAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if sendCmd.Parsed() {
		for _, data := range []*string{sendFrom, sendTo, sendAmount} {
			if *data == "" {
				sendCmd.Usage()
				os.Exit(1)
			}
		}
		amount, _ := strconv.Atoi(*sendAmount)
		cli.send(*sendFrom, *sendTo, amount)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}
}
