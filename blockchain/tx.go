package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"

	wlt "github.com/lshtar13/BlockHoldem/wallet"
)

const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	wallets, _ := wlt.NewWallets()
	wallet := wallets.GetWallet(from)
	fromPubKeyHash := wlt.HashPubkey(wallet.PublicKey)

	acc, validOutputs := bc.FindSpendableOutputs(fromPubKeyHash, amount)
	if acc < amount {
		log.Panic("Error: Not enough funds")
	}

	// inputs
	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)

		for _, out := range outs {
			// txid, output index, signature
			input := TXInput{txID, out, wallet.PublicKey, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, NewTXOutput(amount, to))
	if acc > amount {
		// refund
		outputs = append(outputs, NewTXOutput(acc-amount, from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}

func (tx *Transaction) SetID() error {
	var encoded bytes.Buffer
	var hash [32]byte
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]

	return err
}

func (tx *Transaction) IsCoinbase() bool {
	return tx.Vin[0].Vout == -1
}

// to - address
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, []byte(data), []byte(to)}
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}
