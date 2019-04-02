package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "wallet.dat"

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallets() (*Wallets, error) {                                                  //

	wallets := Wallets{make(map[string]*Wallet)}
	err := wallets.LoadFromFile()

	return &wallets, err
}

func (ws *Wallets) CreateWallet() string { //在钱包库增加一个钱包，用字符串形式的addr标记
	wallet := NewWallet()
	address := string(wallet.GetAddress()[:]) //fmt.Sprintf("%s", wallet.GetAddress())

	ws.Wallets[address] = wallet
	return address
}

func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) GetAddress() []string{
	var addresses []string

	for address := range ws.Wallets{
		addresses = append(addresses, address)
	}
	return addresses
}

func (ws *Wallets) LoadFromFile() error {                     //loads wallets from the file
	if _, err := os.Stat(walletFile); os.IsNotExist(err) { //判断钱包数据库是否存在
		return err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var temp Wallets
	gob.Register(elliptic.P256())                              //???
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&temp)
	if err != nil {
		log.Panic(err)
	}

	ws.Wallets = temp.Wallets
	return nil
}

func (ws Wallets) SaveToFile () {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encode := gob.NewEncoder(&content)
	err := encode.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}