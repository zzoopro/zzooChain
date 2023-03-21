package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/zzoopro/zzoocoin/blockchain"
	"github.com/zzoopro/zzoocoin/utils"
)

type MessageKind int

const (
	MessageNewestBlock 		   MessageKind = iota
	MessageAllBlocksRequest 	
	MessageAllBlocksResponse 	
	
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
	}
}