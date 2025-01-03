package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/lshtar13/BlockHoldem/blockchain"
)

type Handler interface {
	Handle(bc *blockchain.Blockchain) error
}

var nodeAddr string
var miningAddr string
var knownNodes = []string{"localhost:3000"}
var blockTransitChan chan transit
var txTransitChan chan transit

func sendData(addr string, data []byte) error {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available ...\n", addr)
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("Fail to send")
	}
	return nil
}

func router(conn net.Conn, bc *blockchain.Blockchain) error {
	req, _ := io.ReadAll(conn)
	cmd := bytes2Command(req[:cmdLength])
	fmt.Printf("Received %s command\n", cmd)

	var payload Handler
	switch cmd {
	case "version":
		payload = NewVersion()
	case "getBlocks":
		payload = NewGetBlocks()
	case "inv":
		payload = NewInv()
	default:
		return fmt.Errorf("Unknown command: %s", cmd)
	}

	var err error
	var buf bytes.Buffer

	buf.Write(req[cmdLength:])
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&payload)
	if err != nil {
		return err
	}

	payload.Handle(bc)

	conn.Close()

	return nil
}

func StartServer(nodeID, minerAddr string) {
	nodeAddr = fmt.Sprintf("localhost:%s", nodeID)
	miningAddr = minerAddr
	ln, err := net.Listen(protocol, nodeAddr)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	blockTransitChan = make(chan transit)
	defer close(blockTransitChan)

	txTransitChan = make(chan transit)
	defer close(txTransitChan)

	bc, err := blockchain.NewBlockchain(nodeID)

	if nodeAddr != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			go router(conn, bc)
			go transitBlock()
			go transitTx()
		} else {
			fmt.Printf("Error : %s\n", err)
		}
	}
}
