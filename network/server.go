package network

import (
	"fmt"
	"log"
	"net"

	"github.com/lshtar13/BlockHoldem/blockchain"
	"github.com/lshtar13/BlockHoldem/network/common"
	"github.com/lshtar13/BlockHoldem/network/node"
	"github.com/lshtar13/BlockHoldem/network/router"
	"github.com/lshtar13/BlockHoldem/network/service"
)

func StartServer(nodeID, minerAddr string) {
	node.SetNodeAddr(nodeID)
	fmt.Printf("Address : %s\n", node.MySelf())
	ln, err := net.Listen(common.Protocol(), node.MySelf())
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc, err := blockchain.NewBlockchain(nodeID)
	if err != nil {
		log.Panic(err)
	}

	service.PreService(bc)

	service.StartBlockTransit()
	defer service.EndBlockTransit()
	fmt.Println("Start Block Transit!")

	if minerAddr != "" {
		service.StartTxTransit()
		defer service.EndTxTransit()
		fmt.Println("Start Tx Transit!")

		service.StartMining(minerAddr, bc)
		defer service.EndMining()
		fmt.Println("Start Mining!")
	}

	// herererer : verbose!
	fmt.Printf("Start Server! : %s\n", ln.Addr().String())
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error : %s\n", err)
		} else {
			go router.Route(conn, bc)
		}
	}
}
