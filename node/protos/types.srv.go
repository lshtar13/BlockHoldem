package protos

import "github.com/lshtar13/blockchain/chain"

func NewTxInput(in *chain.TXInput) *TxIn {
	return &TxIn{
		Txid:   in.Txid,
		Vout:   int64(in.Vout),
		Sig:    in.Signature,
		PubKey: in.Pubkey,
	}
}

func ToTxInput(in *TxIn) chain.TXInput {
	return chain.TXInput{
		Txid:      in.GetTxid(),
		Vout:      int(in.GetVout()),
		Signature: in.GetSig(),
		Pubkey:    in.GetPubKey(),
	}
}

func NewTxOutput(out *chain.TXOutput) *TxOut {
	return &TxOut{
		Value:      int64(out.Value),
		PubKeyHash: out.PubKeyHash,
	}
}

func ToTxOutput(out *TxOut) chain.TXOutput {
	return chain.TXOutput{
		Value:      int(out.GetValue()),
		PubKeyHash: out.GetPubKeyHash(),
	}
}

func NewTx(tx *chain.Transaction) *Tx {
	vin, vout := []*TxIn{}, []*TxOut{}

	for _, in := range tx.Vin {
		vin = append(vin, NewTxInput(&in))
	}

	for _, out := range tx.Vout {
		vout = append(vout, NewTxOutput(&out))
	}

	return &Tx{
		Id:   tx.ID,
		Vin:  vin,
		Vout: vout,
	}
}

func ToTransaction(tx *Tx) *chain.Transaction {
	vin, vout := []chain.TXInput{}, []chain.TXOutput{}

	for _, in := range tx.Vin {
		vin = append(vin, ToTxInput(in))
	}

	for _, out := range tx.Vout {
		vout = append(vout, ToTxOutput(out))
	}

	return &chain.Transaction{
		ID:   tx.GetId(),
		Vin:  vin,
		Vout: vout,
	}
}

func NewBlk(b *chain.Block) *Blk {
	txs := []*Tx{}
	for _, tx := range b.Transactions {
		txs = append(txs, NewTx(tx))
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

func ToBlock(b *Blk) *chain.Block {
	txs := []*chain.Transaction{}
	for _, tx := range b.Txs {
		txs = append(txs, ToTransaction(tx))
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
