package blockchain

import (
	"fmt"

	"github.com/zzoopro/zzoocoin/db"
	"github.com/zzoopro/zzoocoin/utils"
)

type Block struct {
	Data string `json:"data"`
	Hash string `json:"hash"`
	PrevHash string `json:"prev_hash,omitempty"`
	Height int `json:"height"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}


func createBlock(data string, prevHash string, height int ) *Block {
	block := &Block{
		Data: data,
		Hash: utils.Hash(data, prevHash, fmt.Sprint(height)),
		PrevHash: prevHash,
		Height: height,
	}	
	block.persist()
	return block
}

