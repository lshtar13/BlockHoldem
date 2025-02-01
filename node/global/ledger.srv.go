package global

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/lshtar13/blockchain/utils"
	"google.golang.org/grpc"
)

const (
	EnvName   = "LEDGER_PATH"
	WritePerm = 0644
	AddedKey  = "Added"
)

type Ledger struct {
	NodeID  string
	Central utils.Set
	Miners  utils.Set
}

func GetName(nodeID string) string {
	path := os.Getenv(EnvName)
	if path == "" {
		log.Fatalf("no such env: %s\n", EnvName)
	}

	return fmt.Sprintf("%s/%s.json", path, nodeID)
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
	name := GetName(nodeID)
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("error while reading %s: %v", name, err)
	}

	// if data is empty
	if len(data) == 0 {
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
	}

	result := &Ledger{}
	err = json.Unmarshal(data, result)
	if err != nil {
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
