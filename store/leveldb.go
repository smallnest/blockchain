package store

import (
	"github.com/smallnest/blockchain"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var _ blockchain.Store = &LevelDBStore{}

// LevelDBStore 基于leveldb实现的Store
type LevelDBStore struct {
	db *leveldb.DB
}

func convertLevelDBError(err error) error {
	if err == nil {
		return nil
	}
	switch err {
	case errors.ErrNotFound:
		return blockchain.ErrNotFound
	default:
		return err
	}
}

// NewLevelDBStore 新建一个leveldb store.
func NewLevelDBStore(dataFile string) (*LevelDBStore, error) {
	o := &opt.Options{
		Filter: filter.NewBloomFilter(10),
	}
	db, err := leveldb.OpenFile(dataFile, o)
	if err != nil {
		return nil, err
	}

	store := &LevelDBStore{
		db: db,
	}
	return store, nil
}

// Get 查找指定的区块链
func (s *LevelDBStore) Get(height uint64) (*blockchain.Block, error) {
	key := blockchain.Int2Bytes(height)
	data, err := s.db.Get(key, nil)
	if err != nil {
		return nil, convertLevelDBError(err)
	}

	var block = &blockchain.Block{}
	_, err = block.Unmarshal(data)
	return block, err
}

// Add 增加一个区块.
func (s *LevelDBStore) Add(height uint64, block *blockchain.Block) error {
	key := blockchain.Int2Bytes(height)
	data, err := block.Marshal(nil)
	if err != nil {
		return err
	}

	return s.db.Put(key, data, nil)
}

// GetBatch 得到一批数据.
func (s *LevelDBStore) GetBatch(height uint64, count int) ([]*blockchain.Block, error) {
	var blocks = make([]*blockchain.Block, 0, count)

	key := blockchain.Int2Bytes(height)
	iter := s.db.NewIterator(nil, nil)
	defer iter.Release()

	var err error
	for ok := iter.Seek(key); ok; ok = iter.Next() {
		value := iter.Value()
		var block = &blockchain.Block{}
		_, err = block.Unmarshal(value)
		if err != nil {
			return blocks, err
		}
		blocks = append(blocks, block)
	}

	return blocks, convertLevelDBError(iter.Error())
}

// Exist 检查key是否存在.
func (s *LevelDBStore) Exist(height uint64) (bool, error) {
	key := blockchain.Int2Bytes(height)
	return s.db.Has(key, nil)
}

// Close 关闭db.
func (s *LevelDBStore) Close() error {
	return s.db.Close()
}
