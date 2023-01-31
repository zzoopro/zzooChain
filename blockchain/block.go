package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"github.com/zzoopro/zzoocoin/db"
	"github.com/zzoopro/zzoocoin/utils"
)

type Block struct {
	Data 		string 	`json:"data"`
	Hash 		string 	`json:"hash"`
	PrevHash 	string 	`json:"prev_hash,omitempty"`
	Height 		int 	`json:"height"`
	Difficulty 	int 	`json:"difficulty"`
	Nonce 		int 	`json:"nonce"`
}

var (
	ErrNotFound = errors.New("Block not found.")
	difficulty = 2
)

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}


func createBlock(data string, prevHash string, height int ) *Block {
	block := &Block{
		Data: data,
		Hash: "",
		PrevHash: prevHash,
		Height: height,
		Difficulty: difficulty,
		Nonce: 0,
	}	
	block.mine()
	block.persist()
	return block
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	} 
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		blockAsString := fmt.Sprint(b)
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(blockAsString)))		
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}