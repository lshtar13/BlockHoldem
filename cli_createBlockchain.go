package main

import (
	"fmt"
	"log"
)

func (cli *CLI) createBlockchain(data string) {
	bc, err := CreateBlockchain(data)
	if err != nil {
		log.Panic(err)
	}
	cli.bc = bc
	fmt.Println("Create new blockchain!")
}
