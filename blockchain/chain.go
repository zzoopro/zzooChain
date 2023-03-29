package blockchain

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/zzoopro/zzoocoin/db"
	"github.com/zzoopro/zzoocoin/utils"
)

type blockchain struct {	
	NewestHash 			string `json:"newestHash"`
	Height 				int `json:"height"`
	CurrentDifficulty 	int `json:"currentDifficulty"`
	m 					sync.Mutex
}


type storage interface {
	FindBlock(hash string) []byte
	SaveBlock(hash string, data []byte)
	SaveChain(data []byte)
	LoadChain() []byte
	DeleteAllBlocks()
}

var (
	b *blockchain
	once sync.Once
)

var dbStorage storage = db.DB{}

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
		chainData := dbStorage.LoadChain()
		if chainData == nil {
			b.AddBlock()
		} else {
			b.restore(chainData)
		}
	})
	return b
}

func Status(b *blockchain, rw http.ResponseWriter) {
	b.m.Lock()
	defer b.m.Unlock()
	utils.HandleErr(json.NewEncoder(rw).Encode(b))
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func persistBlockchain(b *blockchain) {
	dbStorage.SaveChain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock() *Block {
	block := createBlock(b.NewestHash, b.Height + 1, GetDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
	return block
}

func Blocks(b *blockchain) []*Block {
	b.m.Lock()
	defer b.m.Unlock()
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
	b.m.Lock()
	defer b.m.Unlock()
	b.CurrentDifficulty = newBlocks[0].Difficulty
	b.Height = len(newBlocks)
	b.NewestHash = newBlocks[0].Hash
	persistBlockchain(b)
	dbStorage.DeleteAllBlocks()
	for _, block := range newBlocks {
		persistBlock(block)
	}
}

func (b *blockchain) AddPeerBlock(newBlock *Block) {
	b.m.Lock()
	m.m.Lock()
	defer b.m.Unlock()
	defer m.m.Unlock()

	b.Height = newBlock.Height
	b.CurrentDifficulty = newBlock.Difficulty
	b.NewestHash = newBlock.Hash

	persistBlock(newBlock)
	persistBlockchain(b)
	
	for _, tx := range newBlock.Transactions {
		if _, ok := m.Txs[tx.Id]; ok {
			delete(m.Txs, tx.Id)
		}
	}
}