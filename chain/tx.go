package chain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	wlt "github.com/lshtar13/blockchain/wallet"
)

const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func (tx *Transaction) Print(nindent int) {
	indent := strings.Repeat(" ", nindent)
	fmt.Printf("%sID : %x\n", indent, tx.ID)
	fmt.Printf("%sVin : \n", indent)
	for _, vin := range tx.Vin {
		vin.Print(nindent * 2)
	}

	fmt.Printf("%sVout : \n", indent)
	for _, out := range tx.Vout {
		out.Print(nindent * 2)
	}
}

func NewUTXOTransaction(nodeID, from, to string, amount int, utxoset *UTXOSet) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	wallets, _ := wlt.NewWallets(nodeID)
	wallet := wallets.GetWallet(from)
	fromPubKeyHash := wlt.HashPubkey(wallet.PublicKey)

	acc, validOutputs := utxoset.FindSpendableOutputs(fromPubKeyHash, amount)
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
	tx.ID = tx.Hash()
	utxoset.Blockchain.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}

func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func DeserializeTransaction(data []byte) (*Transaction, error) {
	var tx Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&tx)

	return &tx, err
}

func (tx *Transaction) Hash() []byte {
	serialized := tx.Serialize()
	hash := sha256.Sum256(serialized)

	return hash[:]
}

func (tx *Transaction) IsCoinbase() bool {
	return tx.Vin[0].Vout == -1
}

// sender, receiver, amount
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}

	return Transaction{tx.ID, inputs, outputs}
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].Pubkey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].Pubkey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Signature = signature
	}
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].Pubkey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].Pubkey = nil

		r, s := big.Int{}, big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x, y := big.Int{}, big.Int{}
		keyLen := len(vin.Pubkey)
		x.SetBytes(vin.Pubkey[:(keyLen / 2)])
		y.SetBytes(vin.Pubkey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if !ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) {
			return false
		}
	}

	return true
}

// to - address
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, []byte(data), []byte(to)}
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.ID = tx.Hash()

	return &tx
}
