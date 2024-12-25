package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	MaxNonce = math.MaxInt64
)

const targetBits = 20

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(MaxBlockSize-targetBits))

	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) prepareData(_nonce int) []byte {
	timestamp, _ := IntToHex(pow.block.Timestamp)
	hardness, _ := IntToHex(int64(targetBits))
	nonce, _ := IntToHex(int64(_nonce))
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			timestamp, hardness, nonce,
		}, []byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hash [32]byte
	var hashInt big.Int

	nonce := 0
	for nonce < MaxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce += 1
		}
	}

	fmt.Printf("\n\n")

	return nonce, hash[:]

}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}
