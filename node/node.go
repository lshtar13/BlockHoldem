package node

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/lshtar13/blockchain/chain"
	"github.com/lshtar13/blockchain/node/global"
	"google.golang.org/grpc"
)

// var memPool map[string]*chain.Transaction
type Srvs struct {
	Blk    *global.BlkSrv
	Tx     *global.TxSrv
	Inv    *global.InvSrv
	Vers   *global.VersSrv
	Ledger *global.LedgerSrv
}

type Node struct {
	id     string
	addr   string
	port   int
	srvs   Srvs
	miner  Miner
	ledger *global.Ledger
	bc     *chain.Blockchain
	// wg unify
	wg  *sync.WaitGroup
	lis net.Listener
}

const (
	BlkType = iota
	TxType
)

func (n *Node) Sync(addr string, wg *sync.WaitGroup) {
	defer wg.Done()

	cc, err := grpc.NewClient(addr)
	if err != nil {
		return
	}
	defer cc.Close()

	myHeight := n.bc.GetBestHeight()
	if hisHeight, err := global.ReqVers(cc, context.Background(), myHeight); err != nil || hisHeight <= myHeight {
		return
	}

	invs, err := global.ReqInv(cc, context.Background())
	if err != nil {
		return
	}

	blk2Req := [][]byte{}
	for _, inv := range invs {
		switch inv.Type {
		case BlkType:
			if _, err := n.bc.GetBlock(inv.Hash); err != nil {
				blk2Req = append(blk2Req, inv.Hash)
			}
		case TxType:
		default:
			return
		}
	}

	if blks, err := global.ReqBlk(cc, context.Background(), blk2Req); err == nil {
		for _, blk := range blks {
			n.bc.AddBlock(blk)
		}
	}
}

func (n *Node) PreService() error {
	// load dial book
	var err error
	n.ledger, err = global.NewLedger(n.id)
	if err != nil {
		return fmt.Errorf("error while creating ledger:%v", err)
	}

	n.ledger.Update()
	wg := new(sync.WaitGroup)
	for central := range n.ledger.Central {
		wg.Add(1)
		go n.Sync(central, wg)
	}
	wg.Wait()

	return nil
}

func (n *Node) Global() {
	defer n.wg.Done()

	server := grpc.NewServer()
	n.srvs.Blk = global.NewBlkSrv(n.bc)
	global.RegisterBlkSrvServer(server, n.srvs.Blk)
	n.srvs.Tx = global.NewTxSrv(n.bc)
	global.RegisterTxSrvServer(server, n.srvs.Tx)
	n.srvs.Inv = global.NewInvSrv(n.bc)
	global.RegisterInvSrvServer(server, n.srvs.Inv)
	n.srvs.Vers = global.NewVersSrv(n.bc)
	global.RegisterVersSrvServer(server, n.srvs.Vers)
	n.srvs.Ledger = global.NewLedgerSrv(n.ledger)
	global.RegisterLedgerSrvServer(server, n.srvs.Ledger)
	if err := server.Serve(n.lis); err != nil {
		log.Fatalf("error while serving:%v\n", err)
	}
}

func (n *Node) Local() {
	defer n.wg.Done()

	server := grpc.NewServer()
	// register Server
	if err := server.Serve(n.lis); err != nil {
		log.Fatalf("error while serving:%v\n", err)
	}
}

func (n *Node) Start(isMiner bool, mineCap int) error {
	var err error
	if bc, err := chain.NewBlockchain(n.id); err != nil {
		return fmt.Errorf("error while getting new blockchain: %v", err)
	} else {
		n.bc = bc
	}

	if n.lis, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", n.port)); err != nil {
		log.Fatalf("err while listening in port %d: %v\n", n.port, err)
	}
	defer n.lis.Close()

	if isMiner {
		n.miner = *NewMiner(n.bc, mineCap)
		go n.miner.Mine()
		defer n.miner.Stop()
	}

	n.wg = new(sync.WaitGroup)

	n.wg.Add(1)
	go n.Global()

	// todo: pre-service
	if err := n.PreService(); err != nil {
		return fmt.Errorf("error while preservicing:%v", err)
	}
	defer n.ledger.Save2FIle()

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
