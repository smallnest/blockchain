package blockchain

import (
	"encoding/binary"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// Store 定义了存储的通用接口。
type Store interface {
	Get(height uint64) (*Block, error)
	Add(height uint64, block *Block) error
	GetBatch(height uint64, count int) ([]*Block, error)
	Exist(height uint64) (bool, error)
	Close() error
}

func int2Bytes(height uint64) []byte {
	var data = make([]byte, 0, 8)
	binary.BigEndian.PutUint64(data, height)
	return data
}

var _ Store = &LevelDBStore{}

// LevelDBStore 基于leveldb实现的Store
type LevelDBStore struct {
	db *leveldb.DB
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
func (s *LevelDBStore) Get(height uint64) (*Block, error) {
	key := int2Bytes(height)
	data, err := s.db.Get(key, nil)
	if err != nil {
		return nil, err
	}

	var block = &Block{}
	_, err = block.Unmarshal(data)
	return block, err
}

// Add 增加一个区块.
func (s *LevelDBStore) Add(height uint64, block *Block) error {
	key := int2Bytes(height)
	data, err := block.Marshal(nil)
	if err != nil {
		return err
	}

	return s.db.Put(key, data, nil)
}

// GetBatch 得到一批数据.
func (s *LevelDBStore) GetBatch(height uint64, count int) ([]*Block, error) {
	var blocks = make([]*Block, 0, count)

	key := int2Bytes(height)
	iter := s.db.NewIterator(nil, nil)
	defer iter.Release()

	var err error
	for ok := iter.Seek(key); ok; ok = iter.Next() {
		value := iter.Value()
		var block = &Block{}
		_, err = block.Unmarshal(value)
		if err != nil {
			return blocks, err
		}
		blocks = append(blocks, block)
	}

	return blocks, iter.Error()
}

// Exist 检查key是否存在.
func (s *LevelDBStore) Exist(height uint64) (bool, error) {
	key := int2Bytes(height)
	return s.db.Has(key, nil)
}

// Close 关闭db.
func (s *LevelDBStore) Close() error {
	return s.db.Close()
}
