package node

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/lshtar13/blockchain/chain"
	"github.com/lshtar13/blockchain/node/global"
	"google.golang.org/grpc"
)

// var memPool map[string]*chain.Transaction

type Node struct {
	id      string
	addr    string
	port    int
	isMiner bool
	book    []string
	bc      *chain.Blockchain
	wg      *sync.WaitGroup
	lis     *net.Listener
}

func (n *Node) PreService() {
	// load dial book
}

func (n *Node) Global() {
	defer n.wg.Done()

	server := grpc.NewServer()
	global.RegisterBlkSrvServer(server, &global.BlkSrv{BC: n.bc})
	global.RegisterTxSrvServer(server, &global.TxSrv{BC: n.bc})
	global.RegisterInvSrvServer(server, &global.InvSrv{BC: n.bc})
	global.RegisterVersSrvServer(server, &global.VersSrv{BC: n.bc})
	err := server.Serve(*n.lis)
	if err != nil {
		log.Fatalf("error while serving:%v\n", err)
	}
}

func (n *Node) Local() {
	defer n.wg.Done()

	server := grpc.NewServer()
	global.RegisterBlkSrvServer(server, &global.BlkSrv{BC: n.bc})
	global.RegisterTxSrvServer(server, &global.TxSrv{BC: n.bc})
	global.RegisterInvSrvServer(server, &global.InvSrv{BC: n.bc})
	global.RegisterVersSrvServer(server, &global.VersSrv{BC: n.bc})
	err := server.Serve(*n.lis)
	if err != nil {
		log.Fatalf("error while serving:%v\n", err)
	}
}

func (n *Node) Start() error {
	bc, err := chain.NewBlockchain(n.id)
	if err != nil {
		return fmt.Errorf("error while getting new blockchain: %v", err)
	}
	n.bc = bc

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", n.port))
	if err != nil {
		log.Fatalf("err while listening in port %d: %v\n", n.port, err)
	}
	n.lis = &lis
	defer lis.Close()

	n.wg = new(sync.WaitGroup)

	n.wg.Add(1)
	go n.Global()

	// todo: pre-service
	n.preservice()

	n.wg.Add(1)
	go n.Local()

	n.wg.Wait()
	return nil
}

func NewNode(id string, addr string, port int) (*Node, error) {
	if lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port)); err != nil {
		return nil, fmt.Errorf("port %d not available", port)
	} else {
		lis.Close()
		return &Node{id: id, addr: addr, port: port}, nil
	}
}
