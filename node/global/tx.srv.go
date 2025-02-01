package global

import (
	"context"
	"io"

	"github.com/lshtar13/blockchain/chain"
	"github.com/lshtar13/blockchain/node/protos"
	"google.golang.org/grpc"
)

type MemPool map[string]*chain.Transaction

func (m MemPool) Add(tx *chain.Transaction) (result int) {
	key := string(tx.ID[:])
	if m[key] != nil {
		result = 1
	} else {
		result = 0
	}
	m[key] = tx

	return
}

func (m MemPool) Remove(id []byte) {
	delete(m, string(id[:]))
}

func (m MemPool) IsExist(id []byte) bool {
	_, isExist := m[string(id[:])]
	return isExist
}

func (m MemPool) Attain(id []byte) *chain.Transaction {
	return m[string(id[:])]
}

type TxSrv struct {
	UnimplementedTxSrvServer
	bc   *chain.Blockchain
	pool MemPool
}

func NewTxSrv(bc *chain.Blockchain) *TxSrv {
	return &TxSrv{bc: bc, pool: make(MemPool)}
}

func (srv *TxSrv) ReqTx(req *TxReq, stream TxSrv_ReqTxServer) error {
	for _, hash := range req.Hash {
		tx, err := srv.bc.FindTransaction(hash)
		if err != nil {
			continue
		}

		if err := stream.Send(protos.NewTx(&tx)); err != nil {
			return err
		}
	}

	return nil
}

func ReqTx(cc *grpc.ClientConn, ctx context.Context, hashes [][]byte) ([]*chain.Transaction, error) {
	client := NewTxSrvClient(cc)
	stream, err := client.ReqTx(ctx, &TxReq{Hash: hashes})
	if err != nil {
		return nil, err
	}

	results := []*chain.Transaction{}
	for {
		tx, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		results = append(results, protos.ToTransaction(tx))
	}

	return results, nil
}

func (srv *TxSrv) SendTx(_ context.Context, tx *protos.Tx) (*protos.Ack, error) {
	transaction := protos.ToTransaction(tx)
	srv.pool.Add(transaction)

	// delete miner.go and invoke minning process herererererer
	return &protos.Ack{}, nil
}

func SendTx(cc *grpc.ClientConn, ctx context.Context, transaction *chain.Transaction) error {
	client := NewTxSrvClient(cc)
	_, err := client.SendTx(ctx, protos.NewTx(transaction))

	return err
}
