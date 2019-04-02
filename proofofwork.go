package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const targetBit = 24

const maxNonce = math.MaxInt64


type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {                                //block + target = pow
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBit))

	var pow = &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte{                             //全部序列化

	data := bytes.Join(
		[][]byte{
			pow.block.HashTransactions(),
			pow.block.PrevBlockHash,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBit)),
			IntToHex(int64(nonce)),
		},
			[]byte{})
	return data
}

func (pow *ProofOfWork) Run() (int,[]byte){           //返回nonce 和 达标Hash
	var hashInt big.Int
	var hash  [32]byte
	nonce := 0

	fmt.Printf("Mining a new blockchain")

	for nonce<maxNonce{

		data:= pow.prepareData(nonce)
		hash = sha256.Sum256(data)                                               // [32]byte
        hashInt.SetBytes(hash[:])                                                   //DEX

        if hashInt.Cmp(pow.target) ==-1{
        	fmt.Printf("%x\n\n",hash)                        //???
        	break
		}else{
         nonce++
		}
	}
	return nonce,hash[:]
}

func (pow *ProofOfWork)Validate() bool{                           // 验证block的Hash与nonce的合理性
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash:= sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid:= (hashInt.Cmp(pow.target)==-1)
	return isValid
}