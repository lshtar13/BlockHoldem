package rpc

import (
	"context"

	"github.com/lshtar13/blockchain/chain"
)

type VersSrv struct {
	UnimplementedVersSrvServer
	bc *chain.Blockchain
}

func (srv *VersSrv) ReqVers(_ context.Context, req *VersReq) (*VersRet, error) {
	bestHeight := int64(srv.bc.GetBestHeight())
	if req.BestHeight > bestHeight {
		// send reqInv ...
	}

	return &VersRet{Version: 1, BestHeight: bestHeight}, nil
}
