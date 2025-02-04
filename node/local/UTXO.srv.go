package local

import (
	context "context"

	"github.com/lshtar13/blockchain/base58"
	"github.com/lshtar13/blockchain/chain"
	"github.com/lshtar13/blockchain/wallet"
	"google.golang.org/grpc"
)

type UTXOSrv struct {
	UnimplementedUTXOSrvServer
	utxoSet *chain.UTXOSet
}

func (srv *UTXOSrv) ReqBalance(_ context.Context, req *BalanceReq) (*Balance, error) {
	addr := req.Addr
	balance := 0
	pubKeyHash := base58.Base58Decode([]byte(addr))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-wallet.AddressChecksumLen]
	UTXOs := srv.utxoSet.FindUTXO(pubKeyHash)

	for _, utxo := range UTXOs {
		balance += utxo.Value
	}

	return &Balance{Balance: int64(balance)}, nil
}

func ReqBalance(cc *grpc.ClientConn, ctx context.Context, addr string) (int, error) {
	client := NewUTXOSrvClient(cc)
	balance, err := client.ReqBalance(ctx, &BalanceReq{Addr: addr})

	return int(balance.Balance), err
}

func NewUtxoSrv(bc *chain.Blockchain) *UTXOSrv {
	return &UTXOSrv{utxoSet: &chain.UTXOSet{Blockchain: bc}}
}
