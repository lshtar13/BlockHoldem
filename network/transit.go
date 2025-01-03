package network

type transit struct {
	data   []byte
	addrTo string
}

func transitBlock() {
	for t := range blockTransitChan {
		sendGetData(t.addrTo, "block", t.data)
	}
}

func transitTx() {
	for t := range blockTransitChan {
		sendGetData(t.addrTo, "tx", t.data)
	}
}
