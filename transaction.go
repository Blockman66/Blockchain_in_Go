package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}


func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func (tx *Transaction) Serialized() []byte{           //全TX序列化
	var encoded bytes.Buffer

	en := gob.NewEncoder(&encoded)
	err := en.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
    return encoded.Bytes()
}

func (tx *Transaction) Hash()[]byte{                 //copy后，ID取0 序列化取Hash

	var hash [32]byte
	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(tx.Serialized())
	return hash[:]
}

func (tx *Transaction) TrimmedCopy() Transaction{                              //返回一个副本，其中vin的sig和key全部取空
	var inputs []TXInput
	var outputs []TXOutput

	for _,in := range tx.Vin{
		inputs = append(inputs,TXInput{in.Txid,in.Vout,nil,nil})
	}
	outputs = tx.Vout

	txcopy := Transaction{tx.ID,inputs,outputs}
	return  txcopy
}

func (tx *Transaction) Sign (privKey ecdsa.PrivateKey, prevTXs map[string]Transaction){         //签名函数
	if tx.IsCoinbase(){                                         //是否元交易
		return
	}

	for _,in := range tx.Vin{                                      // 检验所有输入引用的输出所在交易的ID是否为空（交易有效性判断？）
		if prevTXs[hex.EncodeToString(in.Txid)].ID == nil{
			log.Panic("ERROR: Previous transaction is not correct!")
		}
	}
	txcopy := tx.TrimmedCopy()                                        //创造拷贝
	for  inId, in := range txcopy.Vin{                               //  注意：遍历拷贝（输入中没有sig和key）！
		prevTX := prevTXs[hex.EncodeToString(in.Txid)]                  //声明该输入所引用的过往交易
		txcopy.Vin[inId].Signature = nil                              //双重检验
		txcopy.Vin[inId].PubKey = prevTX.Vout[in.Vout].PubKeyHash    // 用当前输入所引用的过往交易中的输出的key给当前输入的key赋值
		txcopy.ID = tx.Hash()                                         //对只有单个输入中key非空的拷贝求交易ID
		txcopy.Vin[inId].PubKey = nil                                   //再次将key归0

		r,s,err :=ecdsa.Sign(rand.Reader,&privKey,txcopy.ID)              //用私钥对此ID签名（对此输入签名）
		if err !=nil{
			log.Panic(err)
		}
        tx.Vin[inId].Signature = append(r.Bytes(),s.Bytes()...)          //保存签名于原交易（只剩key为空）

	}
}

func (tx Transaction) Verify(prevTXs map[string]Transaction) bool{                             //认证函数
	if tx.IsCoinbase(){                                                       //元交易直接成功
		return true
	}

	for _,in := range tx.Vin{
		if prevTXs[hex.EncodeToString(in.Txid)].ID == nil{
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txcopy := tx.TrimmedCopy()                                                     //创造copy
	curve := elliptic.P256()

	for inId,in := range tx.Vin{                                             //注意：遍历原交易！！
		prevTX := prevTXs[hex.EncodeToString(in.Txid)]
		txcopy.Vin[inId].Signature = nil
		txcopy.Vin[inId].PubKey = prevTX.Vout[in.Vout].PubKeyHash
		txcopy.ID = txcopy.Hash()
		txcopy.Vin[inId].PubKey = nil                                         //同Sign，用副本求每个输入对应的总交易ID，（不涉及本交易的sig和key）

		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.Signature)                                           //  //分离原tx中的range量的sig
		r.SetBytes(in.Signature[:sigLen/2])
		s.SetBytes(in.Signature[sigLen/2:])

		x :=big.Int{}                                         //分离原tx中的range量的key
		y := big.Int{}
		keyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:keyLen/2])
		y.SetBytes(in.PubKey[keyLen/2:])

		rawPubkey := &ecdsa.PublicKey{curve,&x,&y}
		if ecdsa.Verify(rawPubkey,txcopy.ID,&r,&s) == false{       //用公钥对签名标（ID）的签名（r，s）进行验证，ID来自副本，公钥和签名来自原交易。
			return false
		}
	}
	return true
}

func (tx Transaction) String() string{                   //返回一个交易全部信息的string
	var lines []string
	lines = append(lines,fmt.Sprintf("------ Transaction %x",tx.ID))

	for i, input := range tx.Vin{
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Vout))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Vout{
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines,"\n")
}

func NewCoinbaseTX (to, data string) *Transaction {                 //创建并返回创世TX
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.ID = tx.Hash()
	return &tx
}


func NewUTOXTransaction(from, to string, amount int, bc *Blockchain) *Transaction {

	var inputs []TXInput
	var outputs []TXOutput

	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(wallet.PublicKey)
	acc, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)

	if acc<amount{
		log.Panic("Not enough fund!")
	}

	for txid, outs := range validOutputs{
		txID,err := hex.DecodeString(txid)
		if err!=nil{
			log.Panic(err)
		}
		for _,out := range outs{
			input := TXInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs,input)
		}
	}

	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc>amount{
		outputs = append(outputs,*NewTXOutput(acc-amount,from))
	}

	tx := Transaction{nil,inputs,outputs}
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}
