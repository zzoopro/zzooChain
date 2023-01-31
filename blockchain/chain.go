package blockchain

import (
	"sync"

	"github.com/zzoopro/zzoocoin/db"
	"github.com/zzoopro/zzoocoin/utils"
)

type blockchain struct {	
	NewestHash string `json:"newestHash"`
	Height int `json:"height"`
}

var (
	b *blockchain
	once sync.Once
)

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			chainData := db.ChainData()
			if chainData == nil {
				b.AddBlock("Genesis")
			} else {
				b.restore(chainData)
			}
		})
	}	
	return b
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, err := FindBlock(hashCursor)
		utils.HandleErr(err)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height + 1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}