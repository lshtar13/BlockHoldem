package service

import (
	"encoding/hex"

	"github.com/lshtar13/blockchain/blockchain"
	"github.com/lshtar13/blockchain/network/node"
)

var MiningAddr string
var txChan chan *blockchain.Transaction
var mempool = make(map[string]*blockchain.Transaction)

func IsMiner() bool {
	return len(MiningAddr) > 0
}

func AddTx(tx *blockchain.Transaction) {
	mempool[string(tx.ID)] = tx
	txChan <- tx
}

func GetTx(id string) *blockchain.Transaction {
	return mempool[id]
}

func mine(bc *blockchain.Blockchain) {
Collect:
	for range txChan {
		if len(mempool) < 3 {
			continue Collect
		}

		var txs []*blockchain.Transaction
		for id := range mempool {
			tx := mempool[id]
			if bc.VerifyTransaction(tx) {
				txs = append(txs, tx)
			}
		}

		if len(txs) == 0 {
			continue Collect
		}

		cbTx := blockchain.NewCoinbaseTX(MiningAddr, "")
		txs = append(txs, cbTx)

		newBlock := bc.MineBlock(txs)
		UTXOSet := blockchain.UTXOSet{Blockchain: bc}
		UTXOSet.Update(newBlock)

		for _, tx := range txs {
			txID := hex.EncodeToString(tx.ID)
			delete(mempool, txID)
		}

		BroadCastInv("block", []string{node.MySelf()}, [][]byte{newBlock.Hash})
	}
}

func StartMining(addr string, bc *blockchain.Blockchain) {
	MiningAddr = addr
	txChan = make(chan *blockchain.Transaction)
	go mine(bc)
}

func EndMining() {
	close(txChan)
}
