package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

type Blockchain struct {
	tip []byte
	DB  *bolt.DB
}

const dbFile_f = "blockchain_%s.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "Let there be light"

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.DB}
}

func (bc *Blockchain) GetBlock(blkID []byte) (*Block, error) {
	bci := bc.Iterator()
	for {
		b, err := bci.Next()
		if err != nil {
			return nil, err
		}

		if bytes.Equal(b.Hash, blkID) {
			return b, nil
		}

		if len(b.PrevBlockHash) == 0 {
			break
		}
	}

	return nil, fmt.Errorf("cannot find block with id : %x", blkID)
}

func (bc *Blockchain) GetBestHeight() int {
	var bestHeight int
	err := bc.DB.View(func(btx *bolt.Tx) error {
		b := btx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte("l"))

		data := b.Get(lastHash)
		lastBlock, _ := DeserializeBlock(data)

		bestHeight = lastBlock.Height

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return bestHeight
}

func (bc *Blockchain) GetBlockHashes() [][]byte {
	var hashes [][]byte
	bci := bc.Iterator()

	for {
		b, _ := bci.Next()
		hashes = append(hashes, b.Hash)
		if len(b.PrevBlockHash) == 0 {
			break
		}
	}

	return hashes
}

func (bc *Blockchain) FindUTXO() map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block, _ := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			// 사용한 output들 수집
			// 최근의 블록 부터 검사하기 때문에 가능 ...
			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

func (bc *Blockchain) AddBlock(b *Block) error {
	err := bc.DB.Update(func(btx *bolt.Tx) error {
		bucket := btx.Bucket([]byte(blocksBucket))
		blockInDB := bucket.Get(b.Hash)

		if blockInDB != nil {
			return nil
		}

		err := bucket.Put(b.Hash, b.Serialize())
		if err != nil {
			log.Panic(err)
		}

		lastHash := bucket.Get([]byte("l"))
		lastBlockData := bucket.Get(lastHash)
		lastBlock, err := DeserializeBlock(lastBlockData)
		if err != nil {
			log.Panic(err)
		}

		if lastBlock.Height < b.Height {
			err := bucket.Put([]byte("l"), b.Hash)
			if err != nil {
				log.Panic(err)
			}
			bc.tip = b.Hash
		}

		return nil
	})

	return err
}

func (bc *Blockchain) MineBlock(txs []*Transaction) *Block {
	var lastHash []byte
	var lastHeight int

	for _, tx := range txs {
		if !bc.VerifyTransaction(tx) {
			log.Panic("Error: Invalid transaction")
		}
	}

	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		data := b.Get(lastHash)
		lastBlock, _ := DeserializeBlock(data)

		lastHeight = lastBlock.Height

		return nil
	})

	if err != nil {
		log.Panic("Error: extracting lasthash from blockchain")
	}

	newBlock := NewBlock(txs, bc.tip, lastHeight+1)

	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return err
	})
	if err != nil {
		log.Panic(err)
	}

	return newBlock
}

func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block, _ := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}

func NewBlockchain(nodeID string) (*Blockchain, error) {
	dbFile := fmt.Sprintf(dbFile_f, nodeID)
	if !dbExists(dbFile) {
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

func CreateBlockchain(address, nodeID string) (*Blockchain, error) {
	dbFile := fmt.Sprintf(dbFile_f, nodeID)
	if dbExists(dbFile) {
		fmt.Println("blockchain already exists ...")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(btx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)

		b, err := btx.CreateBucket([]byte(blocksBucket))
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

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
