package node

import (
	"github.com/lshtar13/blockchain/chain"
	"github.com/lshtar13/blockchain/utils"
)

type Miner struct {
	bc   *chain.Blockchain
	addr string
	cap  int
	ch   chan *chain.Transaction
	set  utils.Set
	txs  []*chain.Transaction
}

func (m *Miner) Add(tx *chain.Transaction) {
	m.ch <- tx
}

func (m *Miner) Stop() {
	close(m.ch)
}

func (m *Miner) HandleNewBlock(block *chain.Block) {
	// hererererererer
}

func (m *Miner) Mine() {
	for tx := range m.ch {
		key := string(tx.ID[:])
		ntx, _ := m.set.Add(key)
		m.txs = append(m.txs, tx)

		if ntx >= m.cap {
			cbTx := chain.NewCoinbaseTX(m.addr, "")
			block := m.bc.MineBlock(append(m.txs, cbTx))
			m.HandleNewBlock(block)
			m.txs = []*chain.Transaction(nil)
		}
	}
}

func NewMiner(bc *chain.Blockchain, cap int) *Miner {
	return &Miner{bc: bc, cap: cap, ch: make(chan *chain.Transaction), set: make(utils.Set)}
}
