package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zzoopro/zzoocoin/db"
	"github.com/zzoopro/zzoocoin/utils"
)

type Block struct {	
	Hash 			string 	`json:"hash"`
	PrevHash 		string 	`json:"prev_hash,omitempty"`
	Height 			int 	`json:"height"`
	Difficulty 		int 	`json:"difficulty"`
	Nonce 			int 	`json:"nonce"`
	Timestamp   	int     `json:"timestamp"`
	Transactions 	[]*Tx  	`json:"transactions"`
}

var (
	ErrNotFound = errors.New("Block not found.")	
)

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func createBlock(prevHash string, height int ) *Block {
	block := &Block{
		Hash: "",
		PrevHash: prevHash,
		Height: height,
		Difficulty: Blockchain().difficulty(),
		Nonce: 0,
	}	
	block.mine()
	block.Transactions = Mempool.TxToConfirm()
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
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		fmt.Printf("Target: %s\nHash: %s\nNonce: %d\n\n\n", target, hash, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}