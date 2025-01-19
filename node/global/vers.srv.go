package global

import (
	"context"

	"github.com/lshtar13/blockchain/chain"
	"google.golang.org/grpc"
)

type VersSrv struct {
	UnimplementedVersSrvServer
	BC *chain.Blockchain
}

func (srv *VersSrv) ReqVers(_ context.Context, req *VersReq) (*VersRet, error) {
	bestHeight := int64(srv.BC.GetBestHeight())
	if req.BestHeight > bestHeight {
		// send reqInv ...
	}

	return &VersRet{Version: 1, BestHeight: bestHeight}, nil
}

func ReqVers(cc *grpc.ClientConn, ctx context.Context, height int) (int, error) {
	client := NewVersSrvClient(cc)
	vers, err := client.ReqVers(ctx, &VersReq{Version: 1, BestHeight: int64(height)})
	if err != nil {
		return 0, nil
	}

	return int(vers.GetBestHeight()), nil
}
