package blockchain

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/big"
)

// consensus algo / proof of work algo
// secusre blockchain by forcing the network to do work to add a block to the blockchain => computational power => signing => proof of work
// work must be hard to do but proving it must be easy

// take the data from the block
// create a counter (nonce) which starts at 0
// create a hash of the data plus the counter
// check the hash to see if it meets a set of requirements
// Requirements:-
// first few bytes must contain 0s

const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	pow := &ProofOfWork{b, target}
	return pow
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Fatal(err)
	}
	return buff.Bytes()
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

