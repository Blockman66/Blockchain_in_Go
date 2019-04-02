package main

import (
	"bytes"
)

type TXOutput struct {
	Value        int
	PubKeyHash  []byte                                     //两种Hashh
}

func (out TXOutput) Lock(address []byte){                               //解码地址取得公钥，并为输出加密（赋值）
	PubKeyHash := Base58Decode(address)
	PubKeyHash = PubKeyHash[1:len(PubKeyHash)-4]

	out.PubKeyHash = PubKeyHash
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {                //用公钥验证输出（公钥对比）
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func NewTXOutput(value int, address string) *TXOutput {                       //创建新输出！！
	txo := &TXOutput{value,nil}
	txo.Lock([]byte(address))

	return  txo
}

