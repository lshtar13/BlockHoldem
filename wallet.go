package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type WalletMimic struct {
	PrivateKey struct {
		X *big.Int
		Y *big.Int
		D *big.Int
	}

	PublicKey []byte
}

const AddressChecksumLen = 4
const version = byte(0x00)
const walletFile = "wallet.dat"

func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, _ := ecdsa.GenerateKey(curve, rand.Reader)
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

func (w *Wallet) GetAddress() string {
	pubKeyHash := HashPubkey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return string(address)
}

func (w *Wallet) MarshalWallet() []byte {
	var buf bytes.Buffer
	mimic := WalletMimic{}
	mimic.PrivateKey.X = w.PrivateKey.X
	mimic.PrivateKey.Y = w.PrivateKey.Y
	mimic.PrivateKey.D = w.PrivateKey.D
	mimic.PublicKey = w.PublicKey

	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(mimic)
	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}

func UnmarshalWallet(unmarshaled []byte) *Wallet {
	wallet := Wallet{}
	mimic := WalletMimic{}

	decoder := gob.NewDecoder(bytes.NewReader(unmarshaled))
	err := decoder.Decode(&mimic)
	if err != nil {
		log.Panic(err)
	}

	wallet.PublicKey = mimic.PublicKey
	wallet.PrivateKey.Curve = elliptic.P256()
	wallet.PrivateKey.D = mimic.PrivateKey.D
	wallet.PrivateKey.Y = mimic.PrivateKey.Y
	wallet.PrivateKey.X = mimic.PrivateKey.X

	return &wallet
}

func HashPubkey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hahser := ripemd160.New()
	_, err := RIPEMD160Hahser.Write(publicSHA256[:])
	if err != nil {
		log.Panic("Error : hash pubkey")
	}
	publicRIPEMD160 := RIPEMD160Hahser.Sum(nil)

	return publicRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:AddressChecksumLen]
}
