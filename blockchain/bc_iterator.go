package blockchain

import (
	"github.com/boltdb/bolt"
)

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
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
