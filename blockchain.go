package main

import (
	"github.com/boltdb/bolt"
)

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

func (bc *Blockchain) AddBlock(data string) error {
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		bc.tip = b.Get([]byte("l"))

		return nil
	})

	newBlock := NewBlock(data, bc.tip)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return err
	})

	return err
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.db}
}

func (bci *BlockchainIterator) Next() (*Block, error) {
	var block *Block

	err := bci.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(bci.currentHash)
		block, _ = DeserializeBlock(encodedBlock)

		return nil
	})

	bci.currentHash = block.PrevBlockHash

	return block, err
}

func NewBlockchain() (*Blockchain, error) {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
			return err
		} else {
			tip = b.Get([]byte("l"))
			return nil
		}

	})

	return &Blockchain{tip, db}, err
}
