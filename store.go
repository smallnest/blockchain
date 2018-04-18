package blockchain

import (
	"encoding/binary"
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

// Store 定义了存储的通用接口。
type Store interface {
	Get(height uint64) (*Block, error)
	Add(height uint64, block *Block) error
	GetBatch(height uint64, count int) ([]*Block, error)
	Exist(height uint64) (bool, error)
	Close() error
}

func Int2Bytes(height uint64) []byte {
	var data = make([]byte, 8)
	binary.BigEndian.PutUint64(data, height)
	return data
}
