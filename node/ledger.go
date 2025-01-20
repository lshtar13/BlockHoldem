package node

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	EnvName   = "LEDGER_PATH"
	WritePerm = 0644
)

type Ledger struct {
	nodeID     string
	central    []string
	miners     []string
	ordinaries []string
}

func GetName(nodeID string) string {
	path := os.Getenv(EnvName)
	if path == "" {
		log.Fatalf("no such env: %s\n", EnvName)
	}

	return fmt.Sprintf("%s/%s.json", path, nodeID)
}

func (l Ledger) Save2FIle() error {
	data, err := json.Marshal(l)
	if err != nil {
		return fmt.Errorf("error while marshal:%v", err)
	}

	name := GetName(l.nodeID)
	err = os.WriteFile(name, data, WritePerm)
	if err != nil {
		return fmt.Errorf("error while writing to %s:%v", name, err)
	}

	return nil
}

func NewLedger(nodeID string) (*Ledger, error) {
	name := GetName(nodeID)
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("error while reading %s: %v", name, err)
	}

	if len(data) == 0 {
		genName := GetName("genesis")
		data, err = os.ReadFile(genName)
		if err != nil {
			return nil, fmt.Errorf("error while reading %s: %v", genName, err)
		}

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
	result.nodeID = nodeID

	return result, nil
}
