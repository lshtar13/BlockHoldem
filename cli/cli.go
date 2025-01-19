package cli

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/lshtar13/blockchain/chain"
)

type CLI struct {
	nodeID string
	bc     *chain.Blockchain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createBlockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createWallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getBalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listAddresses - Lists all addresses from the wallet file")
	fmt.Println("  printChain - Print all the blocks of the blockchain")
	fmt.Println("  reindexUTXO - Rebuilds the UTXO set")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT -mine - Send AMOUNT of coins from FROM address to TO. Mine on the same node, when -mine is set.")
	fmt.Println("  startNode -miner ADDRESS - Start a node with ID specified in NODE_ID env. var. -miner enables mining")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	cli.nodeID = os.Getenv("NODE_ID")
	if cli.nodeID == "" {
		fmt.Printf("NODE_ID env. var is not set!")
		os.Exit(1)
	}

	createBlockchainCmd := flag.NewFlagSet("createBlockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createWallet", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listAddresses", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexUTXO", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startNode", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "address")
	createAddress := createBlockchainCmd.String("address", "", "address")
	sendFrom := sendCmd.String("from", "", "from")
	sendTo := sendCmd.String("to", "", "to")
	sendAmount := sendCmd.String("amount", "", "amount")
	startNodeMinerAddress := startNodeCmd.String("miner", "", "miner")

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
	case "listAddresses":
		listAddressesCmd.Parse(os.Args[2:])
	case "reindexUTXO":
		reindexUTXOCmd.Parse(os.Args[2:])
	case "startNode":
		startNodeCmd.Parse(os.Args[2:])
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

	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}

	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO()
	}

	if startNodeCmd.Parsed() {
		cli.startNode(*startNodeMinerAddress)
	}
}
