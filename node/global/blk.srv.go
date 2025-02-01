package global

import (
	"context"
	"io"

	"github.com/lshtar13/blockchain/chain"
	"github.com/lshtar13/blockchain/node/protos"
	"google.golang.org/grpc"
)

type BlkSrv struct {
	UnimplementedBlkSrvServer
	bc *chain.Blockchain
}

func NewBlkSrv(bc *chain.Blockchain) *BlkSrv {
	return &BlkSrv{bc: bc}
}

func (srv *BlkSrv) ReqBlk(req *BlkReq, stream BlkSrv_ReqBlkServer) error {
	for _, hash := range req.Hash {
		blk, err := srv.bc.GetBlock(hash)
		if err != nil {
			continue
		}
		if err := stream.Send(protos.NewBlk(blk)); err != nil {
			return err
		}
	}

	return nil
}

func (srv *BlkSrv) SendBlk(_ context.Context, blk *protos.Blk) (*protos.Ack, error) {
	block := protos.ToBlock(blk)
	err := srv.bc.AddBlock(block)
	if err != nil {
		return nil, err
	}

	return &protos.Ack{}, nil
}

func ReqBlk(cc *grpc.ClientConn, ctx context.Context, hashes [][]byte) ([]*chain.Block, error) {
	client := NewBlkSrvClient(cc)
	stream, err := client.ReqBlk(ctx, &BlkReq{Hash: hashes})
	if err != nil {
		return nil, err
	}

	results := []*chain.Block{}
	for {
		blk, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		results = append(results, protos.ToBlock(blk))
	}

	return results, nil
}

func SendBlk(cc *grpc.ClientConn, ctx context.Context, block *chain.Block) error {
	client := NewBlkSrvClient(cc)
	_, err := client.SendBlk(ctx, protos.NewBlk(block))

	return err
}
