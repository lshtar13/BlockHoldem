package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {
	bc *Blockchain
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

	createBlockchainData := createBlockchainCmd.String("address", "", "address")

	switch os.Args[1] {
	case "createBlockchain":
		createBlockchainCmd.Parse(os.Args[2:])
	case "printChain":
		printChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainData == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) createBlockchain(data string) {
	bc, err := CreateBlockchain(data)
	if err != nil {
		log.Panic(err)
	}
	cli.bc = bc
	fmt.Println("Create new blockchain!")
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

	for {
		block, _ := bci.Next()
		fmt.Printf("Prev. hash : %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
