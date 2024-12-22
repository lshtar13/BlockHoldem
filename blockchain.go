package main

import (
	"fmt"
	"os"

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
const genesisCoinbaseData = "Let there be light"

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func (bc *Blockchain) AddBlock(txs []*Transaction) error {
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		bc.tip = b.Get([]byte("l"))

		return nil
	})

	newBlock := NewBlock(txs, bc.tip)

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

func NewBlockchain(address string) (*Blockchain, error) {
	if !dbExists() {
		fmt.Printf("No existing blockchain ...")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Blockchain{tip, db}, nil
}

func CreateBlockchain(address string) (*Blockchain, error) {
	if dbExists() {
		fmt.Println("blockchain already exists ...")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)

		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			return err
		}

		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			return err
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			return err
		}

		tip = genesis.Hash

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Blockchain{tip, db}, err
}
