package blockchain

import "github.com/dgraph-io/badger/v3"

// store blockchain data
//  two entities => blocks (stored with metadata, which describes all the blocks of the chain) and chainState object => stores the state of A chain as unspent transactions, and some metadata,
// with bitcoin => each block has seperate file on the disk => not necessary for smaller blockchain

const (
	dbPath = "./assets/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

func (bc *BlockChain) AddBlock(data string) {

}

func InitBlockchain() *BlockChain {
	var lastHash []byte
	options := badger.DefaultOptions(dbPath)
	options.Dir = dbPath
	options.ValueDir = dbPath
	db, err := badger.Open(options)
	HandleError(err)
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lastHash")); err == badger.ErrKeyNotFound {
			// no blockchain exists yet
			genesis := Genesis()
			err = txn.Set(genesis.Hash, genesis.Serialize())
			HandleError(err)
			err = txn.Set([]byte("lastHash"), genesis.Hash)
			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lastHash"))
			HandleError(err)
			lastHash, err = item.ValueCopy(lastHash)
			return err
		}
	})
	HandleError(err)
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}
