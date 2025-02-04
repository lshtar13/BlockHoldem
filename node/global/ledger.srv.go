package global

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/lshtar13/blockchain/chain"
	"github.com/lshtar13/blockchain/utils"
	"google.golang.org/grpc"
)

const (
	EnvName   = "LEDGER_PATH"
	WritePerm = 0644
	AddedKey  = "Added"
)

type Ledger struct {
	NodeID  string    `json:"NodeID"`
	Central utils.Set `json:"Central"`
	Miners  utils.Set `json:"Miners"`
}

var genLedger = Ledger{NodeID: "genesis", Central: make(utils.Set), Miners: make(utils.Set)}

func GetName(nodeID string) string {
	path := os.Getenv(EnvName)
	if path == "" {
		log.Fatalf("no such env: %s\n", EnvName)
	}

	return fmt.Sprintf("%s/%s.json", path, nodeID)
}

// wrapper
func sendBlock(addr string, block *chain.Block, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	cc, err := grpc.NewClient(addr)
	if err != nil {
		errCh <- fmt.Errorf("error while connecting to %s:%v", addr, err)
		return
	}
	defer cc.Close()

	err = SendBlk(cc, context.Background(), block)
	if err != nil {
		errCh <- fmt.Errorf("error while sending block to %s:%v", addr, err)
		return
	}

	errCh <- nil
}

func sendTransaction(addr string, transaction *chain.Transaction, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	cc, err := grpc.NewClient(addr)
	if err != nil {
		errCh <- fmt.Errorf("error while connecting to %s:%v", addr, err)
		return
	}
	defer cc.Close()

	err = SendTx(cc, context.Background(), transaction)
	if err != nil {
		errCh <- fmt.Errorf("error while sending transaction to %s:%v", addr, err)
		return
	}

	errCh <- nil
}

func (l *Ledger) Propagate(data interface{}) error {
	wg := new(sync.WaitGroup)
	errCh := make(chan error)
	defer close(errCh)

	switch d := data.(type) {
	case *chain.Block:
		for addr := range l.Central {
			wg.Add(1)
			go sendBlock(addr, d, errCh, wg)
		}

	case *chain.Transaction:
		for addr := range l.Miners {
			wg.Add(1)
			go sendTransaction(addr, d, errCh, wg)
		}

	default:
		return fmt.Errorf("%v is neither block nor transaction", d)
	}

	wg.Wait()

	return nil
}

func (l *Ledger) ReqNodes(addr string, newCentralCh chan<- []string, newMinerCh chan<- []string, wg *sync.WaitGroup) {
	defer wg.Done()
	cc, err := grpc.NewClient(addr)
	if err != nil {
		return
	}
	defer cc.Close()

	client := NewLedgerSrvClient(cc)
	newNodes, err := client.ReqNodes(context.Background(), &NodeReq{Type: NodeType_Central, Nodes: l.Central})
	if err == nil {
		newCentralCh <- newNodes.Nodes
	}

	newNodes, err = client.ReqNodes(context.Background(), &NodeReq{Type: NodeType_Miner, Nodes: l.Central})
	if err == nil {
		newMinerCh <- newNodes.Nodes
	}
}

func (l *Ledger) Update() {
	newCentralsCh := make(chan []string)
	newMinersCh := make(chan []string)
	wg := new(sync.WaitGroup)
	for central := range l.Central {
		wg.Add(1)
		go l.ReqNodes(central, newCentralsCh, newMinersCh, wg)
	}

	wg.Wait()
	close(newCentralsCh)
	close(newMinersCh)

	select {
	case newNodes := <-newCentralsCh:
		for _, newNode := range newNodes {
			l.Central.Add(newNode)
		}
	case newNodes := <-newMinersCh:
		for _, newNode := range newNodes {
			l.Miners.Add(newNode)
		}
	}
}

func (l Ledger) Save2FIle() error {
	data, err := json.Marshal(l)
	if err != nil {
		return fmt.Errorf("error while marshal:%v", err)
	}

	name := GetName(l.NodeID)
	err = os.WriteFile(name, data, WritePerm)
	if err != nil {
		return fmt.Errorf("error while writing to %s:%v", name, err)
	}

	return nil
}

func NewLedger(nodeID string) (*Ledger, error) {
	// get data from file
	log.Printf("get new ledger ...\n")
	name := GetName(nodeID)
	data, err := os.ReadFile(name)

	// if data is empty
	if os.IsNotExist(err) && len(data) == 0 {
		// get data from genesis file
		genName := GetName("genesis")
		data, err = os.ReadFile(genName)
		if err != nil {
			return nil, fmt.Errorf("error while reading %s: %v", genName, err)
		}

		// copy
		err = os.WriteFile(name, data, WritePerm)
		if err != nil {
			return nil, fmt.Errorf("error while writing to %s: %v", name, err)
		}

		log.Printf("set default file ...\n")
	} else if err != nil {
		return nil, fmt.Errorf("error while reading %s: %v", name, err)
	}
	log.Printf("read ledger file (%s) ...\n", name)

	result := &Ledger{}
	err = json.Unmarshal(data, result)
	if err != nil {
		log.Println(data)
		return nil, fmt.Errorf("error while unmarshal %s: %v", name, err)
	}
	result.NodeID = nodeID

	return result, nil
}

type LedgerSrv struct {
	UnimplementedLedgerSrvServer
	ledger *Ledger
}

func NewLedgerSrv(ledger *Ledger) *LedgerSrv {
	return &LedgerSrv{ledger: ledger}
}

func (srv *LedgerSrv) ReqNodes(_ context.Context, req *NodeReq) (*NodeRet, error) {
	results := []string{}

	var set utils.Set
	switch req.Type {
	case NodeType_Central:
		set = srv.ledger.Central
	case NodeType_Miner:
		set = srv.ledger.Miners
	default:
		return nil, fmt.Errorf("invlaid req node type:%v", req.Type)
	}

	for node := range set {
		if _, isExist := req.Nodes[node]; !isExist {
			results = append(results, node)
		}
	}

	for node := range req.Nodes {
		set.Add(node)
	}

	return &NodeRet{Nodes: results}, nil
}

func ReqNodes(cc *grpc.ClientConn, ctx context.Context, Type NodeType, nodes map[string]bool) ([]string, error) {
	client := NewLedgerSrvClient(cc)
	ret, err := client.ReqNodes(ctx, &NodeReq{Nodes: nodes, Type: Type})
	if err != nil {
		return nil, fmt.Errorf("error while requesting nodes:%v", err)
	}
	return ret.Nodes, nil
}
