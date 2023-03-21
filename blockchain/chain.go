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
	return b
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func persistBlockchain(b *blockchain) {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height + 1, GetDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
}

func Blocks(b *blockchain) []*Block {
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
 
func CalculateDifficulty(b *blockchain) int {
	blocks := Blocks(b)
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

func GetDifficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height % epoch == 0 {
		return CalculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

func UTxOutputsByAddress(b *blockchain, address string) []*UTxOut {
	var uTxOutputs []*UTxOut
	creatorTxs := make(map[string]bool)
	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxInputs {
				if input.Signature == "COINBASE" {
					break
				}
				if FindTx(b, input.TxID).TxOutputs[input.Index].Address == address {
					creatorTxs[input.TxID] = true
				}
			}
			for index, output := range tx.TxOutputs {
				if output.Address == address {
					if _, ok := creatorTxs[tx.Id]; !ok {
						uTxOutput := &UTxOut{tx.Id, index, output.Amount}
						if !IsOnMempool(uTxOutput) {
							uTxOutputs = append(uTxOutputs, uTxOutput)
						}						
					}
				}				
			}
		}
	}
	return uTxOutputs
}

func BalanceByAddress(b *blockchain, address string) int {
	txOutsByAddress := UTxOutputsByAddress(b, address)
	var total int
	for _, txOut := range txOutsByAddress {
		total += txOut.Amount
	}
	return total
}

func Txs(b *blockchain) []*Tx {
	var txs []*Tx
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transactions...)
	}
	return txs
}

func FindTx(b *blockchain, txId string) *Tx {
	for _, tx := range Txs(b) {
		if tx.Id == txId {
			return tx
		}
	}
	return nil
}

func (b *blockchain) Replace(newBlocks []*Block) {
	b.CurrentDifficulty = newBlocks[0].Difficulty
	b.Height = len(newBlocks)
	b.NewestHash = newBlocks[0].Hash
	persistBlockchain(b)
	db.EmptyBlocks()
	for _, block := range newBlocks {
		persistBlock(block)
	}
}