package main

import (
	"fmt"

	"github.com/inblack67/blockchain-golang/blockchain"
)

func main() {
	blockChain := blockchain.InitBlockchain()
	blockChain.AddBlock("hello")
	blockChain.AddBlock("worlds")
	for _, v := range blockChain.Blocks {
		fmt.Println(v.Hash)
	}
}
