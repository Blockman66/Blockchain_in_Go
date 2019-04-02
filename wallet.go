package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
	"golang.org/x/crypto/ripemd160"
)

const addressChecksumLen = 4
const version = byte(0x00)

type Wallet struct{
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

func NewWallet() *Wallet{                                                     //创建新钱包
	private, public := newKeyPair()
	wallet := Wallet{private,public}

	return  &wallet
}

func newKeyPair()(ecdsa.PrivateKey,[]byte){                                 //生成密钥对（依赖系统级随机函数）

	curve := elliptic.P256()
	private ,err := ecdsa.GenerateKey(curve,rand.Reader)
	if err !=nil{
		log.Panic(err)
	}
	pubkey := append(private.X.Bytes(),private.Y.Bytes()...)

	return *private,pubkey
}

func  (w Wallet)GetAddress()[]byte{                                               //通过pubkey得到address
	publickeyhash := HashPubKey(w.PublicKey)

    versionedpayload := append([]byte{version},publickeyhash...)
    checksum := checksum(versionedpayload)

    fullpayload := append(versionedpayload,checksum...)
    address := Base58Encode(fullpayload)

    return address
}

func ValidateAddress(address string) bool {                                   //checksum地址
	pubkeyHash := Base58Decode([]byte(address))

	actualChecksum := pubkeyHash[len(pubkeyHash)-addressChecksumLen:]
	pubkeyHash = pubkeyHash[:len(pubkeyHash)-addressChecksumLen]

	return bytes.Compare(checksum(pubkeyHash),actualChecksum) == 0
}

func HashPubKey(pubkey []byte)[]byte{                                 //对pubkey进行两种Hash

	pubkeySHA256 := sha256.Sum256(pubkey)
	RIPEMD160Hasher := ripemd160.New()
	RIPEMD160Hasher.Write(pubkeySHA256[:])
	pubkeyRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return pubkeyRIPEMD160
}

func checksum(payload []byte)[]byte{                                   //两次Hash取前四位
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])

	return second[:addressChecksumLen]
}

