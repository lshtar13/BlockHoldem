package global

import (
	"context"

	"github.com/lshtar13/blockchain/chain"
	"google.golang.org/grpc"
)

type InvSrv struct {
	UnimplementedInvSrvServer
	bc *chain.Blockchain
}

func NewInvSrv(bc *chain.Blockchain) *InvSrv {
	return &InvSrv{bc: bc}
}

func (srv *InvSrv) ReqInv(_ context.Context, _ *InvReq) (*InvRet, error) {
	invRet := []*Inv{}
	blocks := srv.bc.GetBlockHashes()
	for _, hash := range blocks {
		invRet = append(invRet, &Inv{Type: InvType_BlkInv, Hash: hash})
	}

	return &InvRet{Invs: invRet}, nil
}

const (
	BLOCK = iota + 1
	TRANSACTION
)

type Inventory struct {
	Type int
	Hash []byte
}

func ReqInv(cc *grpc.ClientConn, ctx context.Context) ([]Inventory, error) {
	client := NewInvSrvClient(cc)
	invs, err := client.ReqInv(ctx, &InvReq{})
	if err != nil {
		return nil, err
	}

	result := []Inventory{}
	for _, inv := range invs.Invs {
		result = append(result, Inventory{int(inv.GetType()), inv.GetHash()})
	}

	return result, nil
}
