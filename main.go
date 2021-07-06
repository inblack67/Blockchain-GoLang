package main

import (
	"fmt"

	"github.com/inblack67/blockchain-golang/blockchain"
)

func main() {
	chain := blockchain.InitBlockchain()
	chain.AddBlock("hello")
	chain.AddBlock("worlds")
	for _, v := range chain.Blocks {
		fmt.Printf("%x \n", v.PrevHash)
		fmt.Printf("%s \n", v.Data)
		fmt.Printf("%x \n", v.Hash)
		pow := blockchain.NewProof(v)
		fmt.Println("POW: ", pow.Validate())
	}
}
