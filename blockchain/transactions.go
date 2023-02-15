package blockchain

import (
	"time"

	"github.com/zzoopro/zzoocoin/utils"
)

const (
	minerReward int = 50
)

type Tx struct {
	Id string				`json:"id"`
	Timestamp int			`json:"timestamp"`
	TxInputs []*TxInput		`json:"txInputs"`
	TxOutputs []*TxOutput	`json:"txOutputs"`
}

type TxInput struct {
	Owner string
	Amount int
}

type TxOutput struct {
	Owner string
	Amount int
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

func makeCoinbaseTx(address string) *Tx {
	inputs := []*TxInput{
		{"COINBASE", minerReward},
	}
	outputs := []*TxOutput{
		{address, minerReward},
	}
	tx := Tx{
		Id: "",
		Timestamp: int(time.Now().Unix()),
		TxInputs: inputs,
		TxOutputs: outputs,
	}
	tx.getId()
	return &tx
}