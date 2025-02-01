package global

import (
	"context"
	"io"

	"github.com/lshtar13/blockchain/chain"
	"google.golang.org/grpc"
)

func newBlk(b *chain.Block) *Blk {
	txs := []*Tx{}
	for _, tx := range b.Transactions {
		txs = append(txs, newTx(tx))
	}

	return &Blk{
		Timestamp:   b.Timestamp,
		Txs:         txs,
		PrevBlkHash: b.PrevBlockHash,
		Hash:        b.Hash,
		Nonce:       int64(b.Nonce),
		Height:      int64(b.Height),
	}
}

func toBlock(b *Blk) *chain.Block {
	txs := []*chain.Transaction{}
	for _, tx := range b.Txs {
		txs = append(txs, toTransaction(tx))
	}

	return &chain.Block{
		Timestamp:     b.GetTimestamp(),
		Transactions:  txs,
		PrevBlockHash: b.GetPrevBlkHash(),
		Hash:          b.GetHash(),
		Nonce:         int(b.GetNonce()),
		Height:        int(b.GetHeight()),
	}
}

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
		if err := stream.Send(newBlk(blk)); err != nil {
			return err
		}
	}

	return nil
}

func (srv *BlkSrv) SendBlk(_ context.Context, blk *Blk) (*Ack, error) {
	block := toBlock(blk)
	err := srv.bc.AddBlock(block)
	if err != nil {
		return nil, err
	}

	return &Ack{}, nil
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

		results = append(results, toBlock(blk))
	}

	return results, nil
}

func SendBlk(cc *grpc.ClientConn, ctx context.Context, block *chain.Block) error {
	client := NewBlkSrvClient(cc)
	_, err := client.SendBlk(ctx, newBlk(block))

	return err
}
