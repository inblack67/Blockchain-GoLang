package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash     []byte // derived from data & prev hash
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func Genesis() *Block {
	block := CreateBlock("Genesis", []byte{})
	return block
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	HandleError(err)
	return res.Bytes()
}

func (b *Block) DeSerialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	HandleError(err)
	return &block
}

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
