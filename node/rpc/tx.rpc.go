package rpc

import "github.com/lshtar13/blockchain/chain"

func newTxInput(in *chain.TXInput) *TxInput {
	return &TxInput{
		Txid:   in.Txid,
		Vout:   int64(in.Vout),
		Sig:    in.Signature,
		PubKey: in.Pubkey,
	}
}

func newTxOutput(out *chain.TXOutput) *TxOutput {
	return &TxOutput{
		Value:      int64(out.Value),
		PubKeyHash: out.PubKeyHash,
	}
}

func newTx(tx *chain.Transaction) *Tx {
	vin, vout := []*TxInput{}, []*TxOutput{}

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
