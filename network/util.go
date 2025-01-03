package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func command2Bytes(cmd string) []byte {
	var result [cmdLength]byte

	for i, c := range cmd {
		result[i] = byte(c)
	}

	return result[:]
}

func bytes2Command(data []byte) string {
	var cmd []byte

	for _, b := range data {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}

	return fmt.Sprintf("%s", cmd)
}
