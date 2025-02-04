package node

import (
	"fmt"

	"github.com/lshtar13/blockchain/chain"
	"github.com/lshtar13/blockchain/node/global"
	"github.com/lshtar13/blockchain/utils"
)

type Miner struct {
	bc      *chain.Blockchain
	addr    string
	cap     int
	transit chan *chain.Transaction
	set     utils.Set
	txs     []*chain.Transaction
	ledger  *global.Ledger
}

func (m *Miner) Stop() {
	close(m.transit)
}

func (m *Miner) HandleNewBlock(block *chain.Block) error {
	err := m.ledger.Propagate(block)
	if err != nil {
		return fmt.Errorf("error while propagating a new block: %v", err)
	}

	UTXOSet := &chain.UTXOSet{Blockchain: m.bc}
	err = UTXOSet.Update(block)
	if err != nil {
		return fmt.Errorf("error while updating with a new block: %v", err)
	}

	return nil
}

func (m *Miner) Mine() {
	for tx := range m.transit {
		key := string(tx.ID[:])
		ntx, _ := m.set.Add(key)
		m.txs = append(m.txs, tx)

		if ntx >= m.cap {
			cbTx := chain.NewCoinbaseTX(m.addr, "")
			block := m.bc.MineBlock(append(m.txs, cbTx))
			m.HandleNewBlock(block)
			m.txs = []*chain.Transaction(nil)
			m.set = make(utils.Set)
		}
	}
}

func NewMiner(bc *chain.Blockchain, addr string, cap int, ledger *global.Ledger) *Miner {
	return &Miner{bc: bc, addr: addr, cap: cap, transit: make(chan *chain.Transaction), set: make(utils.Set), ledger: ledger}
}
