package blockchain

import (
	"errors"
	"time"

	"github.com/zzoopro/zzoocoin/utils"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

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

func makeTx(from, to string, amount int) (*Tx, error) {
	if Blockchain().BalanceByAddress(from) < amount {
		return nil, errors.New("Not enough money.")
	}
	var txInputs []*TxInput
	var txOutputs []*TxOutput
	total := 0
	fromTxOutputs := Blockchain().TxOutputsByAddress(from)
	for _, output := range fromTxOutputs {
		if total >= amount {
			break
		}
		txInput := &TxInput{output.Owner, output.Amount}
		txInputs = append(txInputs, txInput)
		total += txInput.Amount 
	}
	change := total - amount
	if change != 0 {
		changeTxOutput := &TxOutput{from, change}
		txOutputs = append(txOutputs, changeTxOutput)
	}
	txOutput := &TxOutput{to, amount}
	txOutputs = append(txOutputs, txOutput)
	tx := &Tx{
		Id: "",
		Timestamp: int(time.Now().Unix()),
		TxInputs: txInputs,
		TxOutputs: txOutputs,
	}
	tx.getId()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("zzoo", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("zzoo")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}