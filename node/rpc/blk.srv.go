package rpc

import (
	"github.com/lshtar13/blockchain/chain"
)

type BlkSrv struct {
	UnimplementedBlkSrvServer
	bc *chain.Blockchain
}

func newBlk(b *chain.Block) *Blk {
	txs := []*Tx{}
	for _, tx := range b.Transactions {
		txs = append(txs, newTx(tx))
	}

	return &Blk{
		Timestamp:   b.Timestamp,
		PrevBlkHash: b.PrevBlockHash,
		Hash:        b.Hash,
		Nonce:       int64(b.Nonce),
		Height:      int64(b.Height),
	}
}

func (srv *BlkSrv) ReqBlk(req *BlkReq, stream BlkSrv_ReqBlkServer) error {
	for _, hash := range req.Hash {
		_blk, err := srv.bc.GetBlock(hash)
		if err != nil {
			continue
		}

		stream <- newBlk(_blk)
	}

	return nil
}
