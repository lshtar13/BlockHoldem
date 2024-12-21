package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(tgt int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, tgt)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
