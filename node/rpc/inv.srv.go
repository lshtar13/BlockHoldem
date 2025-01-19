package rpc

import (
	"context"

	"github.com/lshtar13/blockchain/chain"
)

type InvSrv struct {
	UnimplementedInvSrvServer
	bc *chain.Blockchain
}

func (srv *InvSrv) ReqInv(_ context.Context, _ *InvReq) (*InvRet, error) {
	invRet := []*Inv{}
	blocks := srv.bc.GetBlockHashes()
	for _, hash := range blocks {
		invRet = append(invRet, &Inv{Type: InvType_BlkInv, Hash: hash})
	}

	return &InvRet{Invs: invRet}, nil
}
