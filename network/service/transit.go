package service

type transit struct {
	data   []byte
	addrTo string
}

var blockTransitChan chan transit
var txTransitChan chan transit

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

func PutBlockTransit(data []byte, addr string) {
	blockTransitChan <- transit{data, addr}
}

func PutTxTransit(data []byte, addr string) {
	txTransitChan <- transit{data, addr}
}

func StartBlockTransit() {
	blockTransitChan = make(chan transit)
	transitBlock()
}

func StartTxTransit() {
	txTransitChan = make(chan transit)
	transitTx()
}

func EndBlockTransit() {
	close(blockTransitChan)
}

func EndTxTransit() {
	close(txTransitChan)
}
