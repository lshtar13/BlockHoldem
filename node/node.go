package node

import (
	"fmt"
	"log"
	"net"

	"github.com/lshtar13/blockchain/chain"
	"github.com/lshtar13/blockchain/node/global"
	"google.golang.org/grpc"
)

type CConn *grpc.ClientConn

var Addr string
var IsMiner bool
var memPool map[string]*chain.Transaction

func StartGlobal(port int, addr string, nodeID string) {
	Addr = addr

	bc, err := chain.NewBlockchain(nodeID)
	if err != nil {
		log.Fatalf("error while getting new blockchain: %v\n", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("err while listening in port %d: %v\n", port, err)
	}
	defer lis.Close()

	server := grpc.NewServer()
	global.RegisterBlkSrvServer(server, &global.BlkSrv{BC: bc})
	global.RegisterTxSrvServer(server, &global.TxSrv{BC: bc})
	global.RegisterInvSrvServer(server, &global.InvSrv{BC: bc})
	global.RegisterVersSrvServer(server, &global.VersSrv{BC: bc})
	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("error while serving:%v\n", err)
	}
}
