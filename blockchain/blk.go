package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

const MaxBlockSize = 256

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
	Height        int
}

func (b *Block) Print() {
	fmt.Printf("============ Block %x ============\n", b.Hash)
	fmt.Printf("Height: %d\n", b.Height)
	fmt.Printf("Nonce: %d\n", b.Nonce)
	fmt.Printf("Prev. block: %x\n", b.PrevBlockHash)
	fmt.Println("Transactions : ")
	for _, t := range b.Transactions {
		t.Print(2)
	}
}

func (b *Block) HashTransactions() []byte {
	var txs [][]byte

	for _, tx := range b.Transactions {
		txs = append(txs, tx.Serialize())
	}
	mTree := NewMerKleTree(txs)

	return mTree.RootNode.Data
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0, height}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce
	return block
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{}, 0)
}

func DeserializeBlock(d []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)

	return &block, err
}
