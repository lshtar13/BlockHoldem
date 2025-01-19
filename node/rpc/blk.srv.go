package rpc

import (
	"log"

	"github.com/lshtar13/blockchain/chain"
)

type BlkSrv struct {
	UnimplementedBlkSrvServer
	bc *chain.Blockchain
}

func (srv *BlkSrv) ReqBlk(req *BlkReq, stream BlkSrv_ReqBlkServer) error {
	for _, hash := range req.Hash {
		_blk, err := srv.bc.GetBlock(hash)
		if err != nil {
			log.Fatalf("no such blk, %x\n", err)
		}

		stream <- _blk.Rpc()
	}
}
