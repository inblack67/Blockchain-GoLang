package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

// - Public database that is distributed across multiple different peers.
// - Does not rely on trust
// - 49% nodes => corrupted data => db will be able to fix itself
// - Composed of blocks => each block contains the data we want to pass around the db as well a hash

// store blockchain data
//  two entities => blocks (stored with metadata, which describes all the blocks of the chain) and chainState object => stores the state of A chain as unspent transactions, and some metadata,
// with bitcoin => each block has seperate file on the disk => not necessary for smaller blockchain

const (
	dbPath   = "./assets/blocks"
	LastHash = "lastHash"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte

	// read only txn
	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(LastHash))
		HandleError(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})
	HandleError(err)
	newBlock := CreateBlock(data, lastHash)

	// read write txn
	err = bc.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		HandleError(err)
		err = txn.Set([]byte(LastHash), newBlock.Hash)
		bc.LastHash = newBlock.Hash
		return err
	})

	HandleError(err)
}

func InitBlockchain() *BlockChain {
	var lastHash []byte
	options := badger.DefaultOptions(dbPath)
	options.Dir = dbPath
	options.ValueDir = dbPath
	db, err := badger.Open(options)
	HandleError(err)
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(LastHash)); err == badger.ErrKeyNotFound {

			fmt.Println("No existing blockchain found")

			// no blockchain exists yet
			genesis := Genesis()

			fmt.Println("Genesis proved")

			err = txn.Set(genesis.Hash, genesis.Serialize())
			HandleError(err)
			err = txn.Set([]byte(LastHash), genesis.Hash)
			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte(LastHash))
			HandleError(err)
			lastHash, err = item.ValueCopy(lastHash)
			return err
		}
	})
	HandleError(err)
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	iterator := &BlockChainIterator{bc.LastHash, bc.Database}
	return iterator
}

// iterating backwards
func (iterator *BlockChainIterator) Next() *Block {
	var block *Block
	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		HandleError(err)
		err = item.Value(func(val []byte) error {
			fmt.Println("====")
			block = DeSerialize(val)
			return nil
		})
		return err
	})
	HandleError(err)
	iterator.CurrentHash = block.PrevHash
	return block
}
