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
	ln, err := net.Listen(common.Protocol(), node.MySelf())
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc, err := blockchain.NewBlockchain(nodeID)

	service.PreService(bc)

	service.StartBlockTransit(bc)
	defer service.EndBlockTransit()

	if minerAddr != "" {
		service.StartTxTransit()
		defer service.EndTxTransit()

		service.StartMining(minerAddr, bc)
		defer service.EndMining()
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			go router.Router(conn, bc)
		} else {
			fmt.Printf("Error : %s\n", err)
		}
	}
}
