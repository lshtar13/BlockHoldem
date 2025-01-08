package node

import "fmt"

var nodeAddr string
var knownNodes = []string{"localhost:3000"}

const nodeVersion = 1

func Version() int {
	return nodeVersion
}

func SetNodeAddr(id string) {
	nodeAddr = fmt.Sprintf("localhost:%s", id)
}

func MySelf() string {
	return nodeAddr
}

func IsCentral() bool {
	return nodeAddr == knownNodes[0]
}

func CentralNode() string {
	return knownNodes[0]
}

func Insert(node string) {
	if !isKnown(node) {
		knownNodes = append(knownNodes, node)
	}
}

func Erase(node string) {
	var updatedNodes []string

	for _, _node := range knownNodes {
		if _node != node {
			updatedNodes = append(updatedNodes, _node)
		}
	}

	knownNodes = updatedNodes
}

func isKnown(node string) bool {
	for _, _node := range knownNodes {
		if _node == node {
			return true
		}
	}

	return false
}

func Disjoint(nodes []string) []string {
	var disjoint []string

NODE:
	for _, node := range knownNodes {
		for _, avoid := range nodes {
			if avoid == node {
				continue NODE
			}
		}

		disjoint = append(disjoint, node)
	}

	return disjoint
}
