package blockchain

import (
	"github.com/dgraph-io/badger/v3"
)

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
			// no blockchain exists yet
			genesis := Genesis()
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
		var encodedBlock []byte
		err = item.Value(func(val []byte) error {
			encodedBlock = val
			return nil
		})
		block = DeSerialize(encodedBlock)
		return err
	})
	HandleError(err)
	return block
}
