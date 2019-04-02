package main          //ok

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	Hash          []byte
	PrevBlockHash []byte
	Nonce         int
}

func (b *Block) Serialize() []byte {                       //block编码
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func DeserializeBlock(b []byte) *Block {                       //block解码（暂时没用）
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}
	return &block
}

func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {               //创建新block并挖矿
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock(coinbase *Transaction) *Block {                //创建创世block并挖矿
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func(b *Block)HashTransactions()[]byte{                             //对一个block中的每一个TX，ID取0编码化取Hash，然后连接后再Hash
	var txHashes [][]byte
	var txHash    [32]byte

	for _,tx := range b.Transactions{
		txHashes = append(txHashes,tx.Hash())
	}
	txHash = sha256.Sum256(bytes.Join(txHashes,[]byte{}))
	return txHash[:]
}