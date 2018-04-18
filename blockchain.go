package blockchain

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"strings"
	"sync"
	"time"
)

// Block 代表区块链中的一块.
type Block struct {
	// 本区块在区块链中的高度
	Height uint64 `json:"height,omitempty"`
	// 本区块产生的时间戳
	Timestamp int64 `json:"timestamp,omitempty"`
	// 数据Data的哈希值
	Hash string `json:"hash,omitempty"`
	// 上一个区块中的Data的哈希值
	PrevHash string `json:"prev_hash,omitempty"`
	// 难度系数
	Difficulty int32 `json:"difficulty"`
	// 随机数
	Nonce int64 `json:"nonce"`
	// 本区块中的数据
	Data []byte `json:"data,omitempty"`
}

// Blockchain 是一条完整的区块链
type Blockchain struct {
	Blocks []*Block
	sync.RWMutex
	Store      Store
	Difficulty int32
	PrefixZero string
}

// LoadFromStore 从文件中加载blockchain.
func (bc *Blockchain) LoadFromStore() error {
	var i uint64
	for {
		block, err := bc.Store.Get(i)
		if err != nil {
			if err == ErrNotFound {
				return nil
			}
			return err
		}
		bc.AddBlock(block)
		i++
	}
}

// GenerateGenesisBlock 初始化创世块.
func (bc *Blockchain) GenerateGenesisBlock() {
	genesisBlock := &Block{
		Height:    0,
		Timestamp: time.Now().Unix(),
		Hash:      hash(&Block{}),
		PrevHash:  "",
		Data:      []byte{},
	}

	bc.Lock()
	bc.AddBlock(genesisBlock)
	bc.Unlock()
}

// AddBlock 在区块链上增加一个区块.
func (bc *Blockchain) AddBlock(block *Block) {
	bc.Blocks = append(bc.Blocks, block)
	bc.Store.Add(block.Height, block)
}

// generateBlock 为数据Data创建一个新的区块
func (bc *Blockchain) generateBlock(prevBlock *Block, data []byte) *Block {
	var newBlock = &Block{}
	newBlock.Height = prevBlock.Height + 1
	newBlock.Timestamp = time.Now().Unix()
	newBlock.PrevHash = prevBlock.Hash
	newBlock.Data = data
	newBlock.Hash = hash(newBlock)

	newBlock.Difficulty = bc.Difficulty

	for i := int64(0); ; i++ {
		newBlock.Nonce = i
		if !validateHash(hash(newBlock), bc.PrefixZero) {
			continue
		} else {
			newBlock.Hash = hash(newBlock)
			break
		}
	}

	if newBlock.Height > 1 && newBlock.Height%3600 == 0 {
		bc.adjustDifficulty()
	}

	return newBlock
}

// hash 计算哈希值.
func hash(block *Block) string {
	h := sha256.New()
	binary.Write(h, binary.BigEndian, block.Height)
	binary.Write(h, binary.BigEndian, block.Timestamp)
	binary.Write(h, binary.BigEndian, block.PrevHash)
	binary.Write(h, binary.BigEndian, block.Nonce)
	binary.Write(h, binary.BigEndian, block.Data)
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// validateBlock 校验块是否合法.
func validateBlock(newBlock, prevBlock *Block) bool {
	if prevBlock.Height+1 != newBlock.Height {
		return false
	}

	if prevBlock.Hash != newBlock.PrevHash {
		return false
	}

	if hash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func validateHash(hash string, prefixZero string) bool {
	return strings.HasPrefix(hash, prefixZero)
}

func (bc *Blockchain) adjustDifficulty() {
	last := len(bc.Blocks) - 1

	tookMs := (bc.Blocks[last].Timestamp - bc.Blocks[last-3600].Timestamp) / 1e6
	tookMs = tookMs / 3600
	if tookMs > 2000 {
		bc.Difficulty--
		bc.PrefixZero = strings.Repeat("0", int(bc.Difficulty))
	} else if tookMs < 500 {
		bc.Difficulty++
		bc.PrefixZero = strings.Repeat("0", int(bc.Difficulty))
	}
}
