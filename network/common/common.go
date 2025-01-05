package common

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

const protocol = "tcp"
const CmdLength = 12

func Protocol() string {
	return protocol
}

func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func Command2Bytes(cmd string) []byte {
	var result [CmdLength]byte

	for i, c := range cmd {
		result[i] = byte(c)
	}

	return result[:]
}

func Bytes2Command(data []byte) string {
	var cmd []byte

	for _, b := range data {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}

	return fmt.Sprintf("%s", cmd)
}
