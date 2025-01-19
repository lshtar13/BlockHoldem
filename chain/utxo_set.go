package chain

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

type UTXOSet struct {
	Blockchain *Blockchain
}

const utxoBucket = "utxo"

func (u *UTXOSet) Reindex() {
	db := u.Blockchain.DB
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(bucketName)
		_, err := tx.CreateBucket(bucketName)

		return err
	})
	if err != nil {
		log.Panic(err)
	}

	UTXO := u.Blockchain.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				return err
			}
			err = b.Put(key, outs.Serialize())
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (u *UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.DB

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := DeserializeOutputs(v)

			for outIdx, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return accumulated, unspentOutputs
}

func (u *UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	db := u.Blockchain.DB

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeserializeOutputs(v)

			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

func (u *UTXOSet) CountTxs() int {
	db := u.Blockchain.DB
	cnt := 0

	err := db.View(func(btx *bolt.Tx) error {
		bucket := btx.Bucket([]byte(utxoBucket))
		cursor := bucket.Cursor()

		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			cnt++
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return cnt
}

func (u *UTXOSet) Update(block *Block) {
	db := u.Blockchain.DB

	err := db.Update(func(btx *bolt.Tx) error {
		var err error
		b := btx.Bucket([]byte(utxoBucket))

		for _, tx := range block.Transactions {
			if tx.IsCoinbase() {
				goto VOUT
			}

			for _, vin := range tx.Vin {
				updatedOuts := TXOutputs{}
				outsBytes := b.Get(vin.Txid)
				fmt.Printf("TxId : %x\n", vin.Txid)
				fmt.Printf("outsBytes : %x, %d\n", outsBytes, len(outsBytes))
				outs := DeserializeOutputs(outsBytes)

				for outIdx, out := range outs.Outputs {
					if outIdx != vin.Vout {
						updatedOuts.Outputs = append(updatedOuts.Outputs, out)
					}
				}

				err = b.Put(vin.Txid, updatedOuts.Serialize())
				if err != nil {
					return err
				}
			}
		VOUT:
			newOutputs := TXOutputs{}
			newOutputs.Outputs = append(newOutputs.Outputs, tx.Vout...)

			err = b.Put(tx.ID, newOutputs.Serialize())
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}
