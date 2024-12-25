package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
)

type Wallets struct {
	Wallets   map[string]*Wallet
	Marshaled *MarshaledWallets
}

type MarshaledWallets struct {
	Wallets map[string][]byte
}

func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	marshaled := MarshaledWallets{}

	wallets.Wallets = make(map[string]*Wallet)
	wallets.Marshaled = &marshaled
	wallets.Marshaled.Wallets = make(map[string][]byte)

	err := wallets.LoadFromFile()

	return &wallets, err
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.GetAddress()

	ws.Wallets[address] = wallet
	marshaled := wallet.MarshalWallet()
	ws.Marshaled.Wallets[address] = marshaled

	return address
}

func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(ws.Marshaled)
	if err != nil {
		log.Panic(err)
	}

	for address, marshaled := range ws.Marshaled.Wallets {
		ws.Wallets[address] = UnmarshalWallet(marshaled)
	}

	return nil
}

func (ws Wallets) SaveToFile() {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(*ws.Marshaled)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletFile, buf.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
