package router

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"github.com/lshtar13/blockchain/blockchain"
	"github.com/lshtar13/blockchain/network/common"
	"github.com/lshtar13/blockchain/network/service"
)

type Handler interface {
	Handle(bc *blockchain.Blockchain) error
}

func Route(conn net.Conn, bc *blockchain.Blockchain) error {
	defer conn.Close()
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
	case "block":
		payload = service.NewBlock()
	case "tx":
		payload = service.NewTx()
	case "getData":
		payload = service.NewGetData()
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}

	var err error
	var buf bytes.Buffer

	fmt.Println("  Decoding ...")
	buf.Write(req[common.CmdLength:])
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(payload)
	if err != nil {
		return err
	}

	payload.Handle(bc)

	return nil
}
