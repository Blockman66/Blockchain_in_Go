package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(num int64) []byte {                                        //临时方案
	var buf = new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)

	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

func ReverseByte(data []byte) {
	for i,j := 0, len(data)-1; i<j; i,j = i+1,j-1{
          data[i],data[j] = data[j],data[i]
	}
}