package blockchain

import (
	"errors"
	"sync"
	"time"

	"github.com/zzoopro/zzoocoin/utils"
	"github.com/zzoopro/zzoocoin/wallet"
)

const (
	minerReward int = 50	
)

type mempool struct {
	Txs 	map[string]*Tx
	m 		sync.Mutex
}

var (
	memOnce sync.Once
)
var m *mempool

func Mempool() *mempool {
	memOnce.Do(func() {
		m = &mempool{
			Txs: make(map[string]*Tx),
		}
	})
	return m
}

// errors
var (
	ErrNoMoney = errors.New("Not enough money.")
)

type Tx struct {
	Id string				`json:"id"`
	Timestamp int			`json:"timestamp"`
	TxInputs []*TxInput		`json:"txInputs"`
	TxOutputs []*TxOutput	`json:"txOutputs"`
}

type TxInput struct {
	TxID 		string	`json:"txId"`
	Index		int		`json:"index"`
	Signature 	string	`json:"signature"`
}

type TxOutput struct {
	Address string	`json:"address"`
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

func (t *Tx) sign() {
	for _, txInput := range t.TxInputs {
		txInput.Signature = wallet.Sign(t.Id, wallet.Wallet())
	}
}

func validate(tx *Tx) bool {
	valid := true
	for _, txIn := range tx.TxInputs {
		prevTx := FindTx(Blockchain(), txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		address := prevTx.TxOutputs[txIn.Index].Address
		valid = wallet.Verify(txIn.Signature, tx.Id, address)
		if !valid {
			break
		}
	}
	return valid
}

func IsOnMempool(uTxOut *UTxOut) bool {	
	exists := false

	Outer:
	for _, tx := range Mempool().Txs {		
		for _, input := range tx.TxInputs {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
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
	if BalanceByAddress(Blockchain(), from) < amount {
		return nil, ErrNoMoney
	}
	var txInputs []*TxInput
	var txOutputs []*TxOutput
	total := 0
	uTxOutputs := UTxOutputsByAddress(Blockchain(), from)

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
	tx.sign()
	valid := validate(tx)
	if !valid {
		return nil, ErrNoMoney
	}
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) (*Tx, error) {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return nil, err
	}
	m.Txs[tx.Id] = tx
	return tx, nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	var txs []*Tx	
	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	txs = append(txs, coinbase)
	m.Txs = make(map[string]*Tx)
	return txs
}

func (m *mempool) AddPeerTx(tx *Tx) {
	m.m.Lock()
	defer m.m.Unlock()
	m.Txs[tx.Id] = tx
}