package utils

import (
	"bytes"
	"encoding/binary"
)

func IntToHex(tgt int64) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, tgt)

	return buff.Bytes(), err
}

func ReverseBytes(arr []byte) {
	for i, l := 0, len(arr)-1; i < l; i, l = i+1, l-1 {
		arr[i], arr[l] = arr[l], arr[i]
	}
}
