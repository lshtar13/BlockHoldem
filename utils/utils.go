package utils

import (
	"bytes"
	"encoding/binary"
)

type Set map[string]bool

func (s Set) Add(key string) (int, bool) {
	result := s[key]
	s[key] = true

	return s.Length(), result
}

func (s Set) Remove(key string) {
	delete(s, key)
}

func (s Set) IsExist(key string) bool {
	_, isExist := s[key]
	return isExist
}

func (s Set) Attain(key string) bool {
	return s[key]
}

func (s Set) Length() int {
	return len(s)
}

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
