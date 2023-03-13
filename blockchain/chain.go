package blockchain

import (
	"sync"

	"github.com/zzoopro/zzoocoin/db"
	"github.com/zzoopro/zzoocoin/utils"
)

type blockchain struct {	
	NewestHash 			string `json:"newestHash"`
	Height 				int `json:"height"`
	CurrentDifficulty 	int `json:"currentDifficulty"`
}

var (
	b *blockchain
	once sync.Once
)

const (
	defaultDifficulty int = 2
	epoch int = 5	
	blockTime int = 2
	allowedRange int = 2
)

func Blockchain() *blockchain {
	if b == nil { 
		once.Do(func() {
			b = &blockchain{				
				Height: 0,
			}
			chainData := db.ChainData()
			if chainData == nil {
				b.AddBlock()
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

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height % epoch == 0 {
		return b.calculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}
}
 
func (b *blockchain) calculateDifficulty() int {
	blocks := b.Blocks()
	newestBlock := blocks[0]
	lastCalculatedBlock := blocks[epoch - 1]
	actualMinute := (newestBlock.Timestamp / 60) - (lastCalculatedBlock.Timestamp / 60)
	expectedMinute := epoch * blockTime
	if actualMinute > (expectedMinute + allowedRange) {
		return b.CurrentDifficulty - 1
	} else if actualMinute < (expectedMinute - allowedRange) {
		return b.CurrentDifficulty + 1
	}  
	return b.CurrentDifficulty	
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

func (b *blockchain) allTxOutPuts() []*TxOutput {
	blocks := b.Blocks()
	var txOutputs []*TxOutput 
	for _, block := range blocks {
		for _, tx := range block.Transactions {
			 txOutputs = append(txOutputs, tx.TxOutputs...)
		}
	}
	return txOutputs
}

func (b *blockchain) TxOutputsByAddress(address string) []*TxOutput {
	allTxOutPuts := b.allTxOutPuts()
	var txOutsByAddress []*TxOutput
	for _, txOutPut := range allTxOutPuts {
		if txOutPut.Owner == address {
			txOutsByAddress = append(txOutsByAddress, txOutPut)
		}		
	}
	return txOutsByAddress
}

func (b *blockchain) BalanceByAddress(address string) int {
	txOutsByAddress := b.TxOutputsByAddress(address)
	var total int
	for _, txOut := range txOutsByAddress {
		total += txOut.Amount
	}
	return total
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height + 1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}