package main

import (
	"bytes"
	"encoding/binary"
)

func IntToHex(tgt int64) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, tgt)

	return buff.Bytes(), err
}
