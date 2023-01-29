package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)


type Block struct {
	Data string `json:"data"`
	Hash string `json:"hash"`
	PrevHash string `json:"prev_hash,omitempty"`
	Height int `json:"height"`
}

type blockchain struct {
	bolcks []*Block
}


var (
	b *blockchain
	once sync.Once
)	

var (
	ErrNotFound = errors.New("Block Not Found")
)

func (b *Block) setHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlockchain := len(GetBlockchain().bolcks)
	if totalBlockchain == 0 {
		return ""
	}
	return GetBlockchain().bolcks[len(b.bolcks) - 1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash(), len(GetBlockchain().bolcks) + 1}
	newBlock.setHash()
	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	b.bolcks = append(b.bolcks, createBlock(data)) 
}

func (b *blockchain) AllBlocks() []*Block {
	return b.bolcks
}

func (b *blockchain) FindBlock(height int) (*Block, error){
	if height > len(b.bolcks) {
		return nil, ErrNotFound
	}
	return b.bolcks[height - 1], nil
}

func GetBlockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis Block")
		})		
	}
	return b
}