package blockchain

import (
	"reflect"
	"sync"
	"testing"

	"github.com/zzoopro/zzoocoin/utils"
)

type fakeDB struct {
	fakeLoadChain func() []byte
	fakeFindBlock func() []byte
}

func (f fakeDB) LoadChain() []byte {
	return f.fakeLoadChain()
}
func (f fakeDB) FindBlock(hash string) []byte {
	return f.fakeFindBlock()
}
func (fakeDB) SaveBlock(hash string, data []byte) {}
func (fakeDB) SaveChain(data []byte) {}
func (fakeDB) DeleteAllBlocks() {}

func TestBlockchain(t *testing.T){
	t.Run("Should create blockchain", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeLoadChain: func() []byte {return nil},
		}
		bc := Blockchain()
		if bc.Height != 1 {
			t.Error("Blockchain() should create a blockchain.")
		}
	})

	t.Run("Should restore blockchain", func(t *testing.T) {
		once = *new(sync.Once)
		dbStorage = fakeDB{
			fakeLoadChain: func() []byte {
				bc := &blockchain{Height: 2, NewestHash: "xxx", CurrentDifficulty: 1}
				return utils.ToBytes[any](bc)
			},
		}
		bc := Blockchain()
		if bc.Height != 2 {
			t.Errorf("Blockchain() should restore a blockchain with a height of %d, got %d", 2, bc.Height)
		}
	})
}

func TestBlocks (t *testing.T){
	blocks := []*Block{
		{PrevHash: "x"},
		{PrevHash: ""},
	}
	fakeBlocks := 0
	dbStorage = fakeDB{
		fakeFindBlock: func() []byte {
			defer func ()  {
				fakeBlocks++	
			}()
			return utils.ToBytes[any](blocks)
		},
	}
	bc := &blockchain{}
	blocksResult := Blocks(bc)
	if reflect.TypeOf(blocksResult) != reflect.TypeOf([]*Block{}) {
		t.Error("Blocks() should return a slice of blocks.")
	}
}

func TestFindTx(t *testing.T){
	t.Run("Tx not found.", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Height: 2,					
					Transactions: []*Tx{},
				}
				return utils.ToBytes[any](b)
			},
		}
		tx := FindTx(&blockchain{NewestHash: "x"}, "test")
		if tx != nil {
			t.Error("Tx should be not found.")
		}
	})

	t.Run("Tx should be found.", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Height: 2,					
					Transactions: []*Tx{
						{Id: "test"},
					},
				}
				return utils.ToBytes[any](b)
			},
		}
		tx := FindTx(&blockchain{NewestHash: "x"}, "test")
		if tx == nil {
			t.Error("Tx should be found.")
		}
	})
}

func TestGetDifficulty(t *testing.T){
	blocks := []*Block{
		{PrevHash: "x"},
		{PrevHash: "x"},
		{PrevHash: "x"},
		{PrevHash: "x"},
		{PrevHash: ""},
	}
	fakeBlock := 0
	dbStorage = fakeDB{
		fakeFindBlock: func() []byte {
			defer func ()  {
				fakeBlock++
			}()
			return utils.ToBytes[any](blocks[fakeBlock])
		},
	}
	type test struct {
		height 	int
		want 	int
	}
	tests := []test{
		{height: 0, want: defaultDifficulty},
		{height: 2, want: defaultDifficulty},
		{height: 5, want: 3 },
	}
	for _, tc := range tests {
		bc := &blockchain{Height: tc.height}
		got := GetDifficulty(bc)
		if got != tc.want {
			t.Errorf("getDifficulty() should return %d got %d", tc.want, got)
		}
	}
}

func TestAddPeerBlock(t *testing.T) {
	bc := &blockchain{
		Height: 1,
		CurrentDifficulty: 1,
		NewestHash: "xx",
	}
	m.Txs["test"] = &Tx{}
	newBlock := &Block{
		Difficulty: 2,
		Hash: "test",
		Transactions: []*Tx{
			{Id: "test"},
		},
	}
	bc.AddPeerBlock(newBlock)
	if bc.CurrentDifficulty != 2 || bc.Height != 2 || bc.NewestHash != "test" {
		t.Error("AddPeerBlock should mutate the blockchain.")
	}
}

func TestReplace(t *testing.T){
	bc := &blockchain{
		Height: 1,
		CurrentDifficulty: 1,
		NewestHash: "xx",
	}
	blocks := []*Block{
		{Difficulty: 2, Hash: "test"},
		{Difficulty: 2, Hash: "test"},		
	}
	bc.Replace(blocks)
	if bc.CurrentDifficulty != 2 || bc.Height != 2 || bc.NewestHash != "test" {
		t.Error("Replace() should mutate the blockchain.")
	}
}