package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/inblack67/blockchain-golang/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage => ")
	fmt.Println(" add -block BLOCK_DATA<string> - add a block to the chain")
	fmt.Println(" print - Prints the block to the chain")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit() // exit after closing all the goroutines properly => unlike via os.exit => badgerDB needs this to grabage collect the keys and values
	}

}

func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Block added")
}

func (cli *CommandLine) printChain() {
	iterator := cli.blockchain.Iterator()
	for {
		block := iterator.Next()
		// fmt.Printf("%x \n", block.PrevHash)
		fmt.Printf("%s \n", block.Data)
		// fmt.Printf("%x \n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Println("POW: ", pow.Validate())
		if string(block.PrevHash) == "" {
			// genesis block has it's prevHash as empty string
			break
		}
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()
	addBlockFlag := flag.NewFlagSet("add", flag.ExitOnError)
	addBlockData := addBlockFlag.String("block", "", "Block data")
	printBlockFlag := flag.NewFlagSet("print", flag.ExitOnError)

	switch os.Args[1] {
	case "add":
		err := addBlockFlag.Parse(os.Args[2:])
		blockchain.HandleError(err)
	case "print":
		err := printBlockFlag.Parse(os.Args[2:])
		blockchain.HandleError(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockFlag.Parsed() {
		if *addBlockData == "" {
			addBlockFlag.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printBlockFlag.Parsed() {
		cli.printChain()
	}

}

func main() {
	defer os.Exit(0)
	chain := blockchain.InitBlockchain()
	defer chain.Database.Close() // defer only runs if goroutines exit gracefully

	cli := CommandLine{chain}
	cli.run()
}
