package main

import "bytes"

type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte                                                         //两种Hash
}

func (in *TXInput) UsesKey(pubKeyHash []byte) bool {                      //用原始公钥两种Hash后印证in中pubkey

	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
