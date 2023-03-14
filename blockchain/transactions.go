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
	TxID 	string	`json:"txId"`
	Index	int		`json:"index"`
	Owner 	string	`json:"owner"`
}

type TxOutput struct {
	Owner string	`json:"owner"`
	Amount int		`json:"amount"`
}

type UTxOut struct {
	TxID 	string	`json:"txId"`
	Index	int		`json:"index"`
	Amount 	int	`json:"amount"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

func IsOnMempool(uTxOut *UTxOut) bool {
	exists := false
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxInputs {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
			}
		}
	}
	return exists
}

func makeCoinbaseTx(address string) *Tx {
	inputs := []*TxInput{
		{"", -1, "COINBASE"},
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
	uTxOutputs := Blockchain().UTxOutputsByAddress(from)

	for _, uTxOut := range uTxOutputs {
		if total >= amount {
			break
		}
		txInput := TxInput{uTxOut.TxID, uTxOut.Index, from}
		txInputs = append(txInputs, &txInput)
		total += uTxOut.Amount
	}
	
	if change := total - amount; change != 0 {
		changeTxOutput := &TxOutput{from, change}
		txOutputs = append(txOutputs, changeTxOutput)  
	}
	TxOutput := &TxOutput{to, amount}
	txOutputs = append(txOutputs, TxOutput)
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