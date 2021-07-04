package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Hash     []byte // derived from data & prev hash
	Data     []byte
	PrevHash []byte
}

type BlockChain struct {
	blocks []*Block
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

func Genesis() *Block {
	block := CreateBlock("Genesis", []byte{})
	return block
}

func InitBlockchain() *BlockChain {
	genesisBlock := Genesis()
	blockChain := &BlockChain{[]*Block{genesisBlock}}
	return blockChain
}

func main() {
	blockChain := InitBlockchain()
	blockChain.AddBlock("hello")
	blockChain.AddBlock("worlds")
	for _, v := range blockChain.blocks {
		fmt.Println(v.Hash)
	}
}
