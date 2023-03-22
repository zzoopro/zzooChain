package p2p

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zzoopro/zzoocoin/blockchain"
	"github.com/zzoopro/zzoocoin/utils"
)

type MessageKind int

const (
	MessageNewestBlock 		   MessageKind = iota
	MessageAllBlocksRequest 	
	MessageAllBlocksResponse 		
	MessageNewBlockNotify
	MessageNewTxNotify
	MessageNewPeerNotify
)

type Message struct {
	Kind 		MessageKind
	Payload 	[]byte
}


func makeMessage[T any](kind MessageKind, payload T) []byte {
	message := Message{
		Kind: kind,		
		Payload: utils.ToJSON[any](payload),
	}		
	return utils.ToJSON[any](message)
}

func sendNewestBlock(p *peer) {
	fmt.Printf("Sending newest block to %s\n", p.key)
	block, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)

	message := makeMessage[any](MessageNewestBlock, block)
	p.inbox <- message
}

func requestAllBlocks(p *peer) {
	message := makeMessage[any](MessageAllBlocksRequest, nil)
	p.inbox <- message
}

func sendAllBlocks(p *peer) {
	message := makeMessage[any](MessageAllBlocksResponse, blockchain.Blocks(blockchain.Blockchain()))
	p.inbox <- message
}

func notifyNewBlock(block *blockchain.Block, p *peer) {
	message := makeMessage[any](MessageNewBlockNotify, block)
	p.inbox <- message
}

func notifyNewTx(tx *blockchain.Tx, p *peer) {
	message := makeMessage[any](MessageNewTxNotify, tx)
	p.inbox <- message
}

func notifyNewPeer(p *peer, ip string) {
	message := makeMessage[any](MessageNewPeerNotify, ip)
	p.inbox <- message
}

func handleMsg(m *Message, p *peer) {
	switch m.Kind { 
	case MessageNewestBlock:
		fmt.Printf("Recived the newest block from %s\n", p.key)
		var newestBlock blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &newestBlock))
		if newestBlock.Height >= blockchain.Blockchain().Height {
			fmt.Printf("Requesting all blocks from %s\n", p.key)
			requestAllBlocks(p)
		} else {
			fmt.Printf("Sending newest block to %s\n", p.key)
			sendNewestBlock(p)
		}

	case MessageAllBlocksRequest:
		fmt.Printf("Sending all blocks to %s\n", p.key)
		sendAllBlocks(p)

	case MessageAllBlocksResponse:
		fmt.Printf("Recived all blocks from %s\n", p.key)
		var allBlocks []*blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &allBlocks))
		blockchain.Blockchain().Replace(allBlocks)

	case MessageNewBlockNotify:
		var payload *blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Blockchain().AddPeerBlock(payload)

	case MessageNewTxNotify:
		var payload *blockchain.Tx
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Mempool().AddPeerTx(payload)

	case MessageNewPeerNotify:
		var payload string
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		parts := strings.Split(payload, ":")
		AddPeer(parts[0], parts[1], parts[2], false)
	}
}