package blockchain

type Block struct {
	Hash     []byte // derived from data & prev hash
	Data     []byte
	PrevHash []byte
	Nonce    int
}

type BlockChain struct {
	Blocks []*Block
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
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
