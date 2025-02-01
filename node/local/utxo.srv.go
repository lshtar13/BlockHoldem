package local

import "github.com/lshtar13/blockchain/chain"

type UtxoSrv struct {
	UnimplementedUtxoSrvServer
	utxoSet *chain.UTXOSet
}

func NewUtxoSrv(bc *chain.Blockchain) *UtxoSrv {
	return &UtxoSrv{utxoSet: &chain.UTXOSet{Blockchain: bc}}
}
