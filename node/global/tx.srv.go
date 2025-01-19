package global

import (
	"context"
	"io"

	"github.com/lshtar13/blockchain/chain"
	"google.golang.org/grpc"
)

func newTxInput(in *chain.TXInput) *TxIn {
	return &TxIn{
		Txid:   in.Txid,
		Vout:   int64(in.Vout),
		Sig:    in.Signature,
		PubKey: in.Pubkey,
	}
}

func toTxInput(in *TxIn) chain.TXInput {
	return chain.TXInput{
		Txid:      in.GetTxid(),
		Vout:      int(in.GetVout()),
		Signature: in.GetSig(),
		Pubkey:    in.GetPubKey(),
	}
}

func newTxOutput(out *chain.TXOutput) *TxOut {
	return &TxOut{
		Value:      int64(out.Value),
		PubKeyHash: out.PubKeyHash,
	}
}

func toTxOutput(out *TxOut) chain.TXOutput {
	return chain.TXOutput{
		Value:      int(out.GetValue()),
		PubKeyHash: out.GetPubKeyHash(),
	}
}

func newTx(tx *chain.Transaction) *Tx {
	vin, vout := []*TxIn{}, []*TxOut{}

	for _, in := range tx.Vin {
		vin = append(vin, newTxInput(&in))
	}

	for _, out := range tx.Vout {
		vout = append(vout, newTxOutput(&out))
	}

	return &Tx{
		Id:   tx.ID,
		Vin:  vin,
		Vout: vout,
	}
}

func toTransaction(tx *Tx) *chain.Transaction {
	vin, vout := []chain.TXInput{}, []chain.TXOutput{}

	for _, in := range tx.Vin {
		vin = append(vin, toTxInput(in))
	}

	for _, out := range tx.Vout {
		vout = append(vout, toTxOutput(out))
	}

	return &chain.Transaction{
		ID:   tx.GetId(),
		Vin:  vin,
		Vout: vout,
	}
}

type TxSrv struct {
	UnimplementedTxSrvServer
	BC *chain.Blockchain
}

func (srv *TxSrv) ReqTx(req *TxReq, stream TxSrv_ReqTxServer) error {
	for _, hash := range req.Hash {
		tx, err := srv.BC.FindTransaction(hash)
		if err != nil {
			continue
		}

		if err := stream.Send(newTx(&tx)); err != nil {
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

		results = append(results, toTransaction(tx))
	}

	return results, nil
}

func SendTx(cc *grpc.ClientConn, ctx context.Context, transaction *chain.Transaction) error {
	client := NewTxSrvClient(cc)
	_, err := client.SendTx(ctx, newTx(transaction))

	return err
}
