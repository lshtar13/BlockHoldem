package router

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"github.com/lshtar13/BlockHoldem/blockchain"
	"github.com/lshtar13/BlockHoldem/network/common"
	"github.com/lshtar13/BlockHoldem/network/service"
)

type Handler interface {
	Handle(bc *blockchain.Blockchain) error
}

func Router(conn net.Conn, bc *blockchain.Blockchain) error {
	req, _ := io.ReadAll(conn)
	cmd := common.Bytes2Command(req[:common.CmdLength])
	fmt.Printf("Received %s command\n", cmd)

	var payload Handler
	switch cmd {
	case "version":
		payload = service.NewVersion()
	case "getBlocks":
		payload = service.NewGetBlocks()
	case "inv":
		payload = service.NewInv()
	default:
		return fmt.Errorf("Unknown command: %s", cmd)
	}

	var err error
	var buf bytes.Buffer

	buf.Write(req[common.CmdLength:])
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(&payload)
	if err != nil {
		return err
	}

	payload.Handle(bc)

	conn.Close()

	return nil
}
