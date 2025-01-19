package service

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/lshtar13/blockchain/blockchain"
	"github.com/lshtar13/blockchain/network/common"
	"github.com/lshtar13/blockchain/network/node"
)

func SendData(addr string, data []byte) error {
	conn, err := net.Dial(common.Protocol(), addr)
	if err != nil {
		fmt.Printf("%s is not available ...\n", addr)
		node.Erase(addr)
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("fail to send")
	} else {
		fmt.Printf("send data to %s\n", addr)
	}
	return nil
}

func BroadCastInv(sort string, disjoint []string, items [][]byte) {
	for _, node := range node.Disjoint(disjoint) {
		sendInv(node, sort, items)
	}
}

func PreService(bc *blockchain.Blockchain) {
	fmt.Println("PreService:")
	if !node.IsCentral() {
		sendVersion(node.CentralNode(), bc)
	}
	fmt.Println("done(Preservice)")
}
