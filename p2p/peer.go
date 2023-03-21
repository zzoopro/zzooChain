package p2p

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type peers struct {
	v map[string]*peer
	m sync.Mutex
}

type peer struct {
	connection 	*websocket.Conn
	inbox 		chan []byte
	key 		string
	address 	string
	port		string
}

var Peers peers = peers{
	v: make(map[string]*peer),
}

func AllPeers(p *peers) []string {
	p.m.Lock()
	defer p.m.Unlock()
	var keys []string
	for key := range p.v {
		keys = append(keys, key)
	}
	return keys
}

func (p *peer) close() {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	p.connection.Close()
	delete(Peers.v, p.key)	
}

func (p *peer) read() {
	defer p.close()
	for {
		message := Message{}
		err := p.connection.ReadJSON(&message)
		if err != nil {
			break 
		}
		handleMsg(&message, p)
	}
}

func (p *peer) write() {
	defer p.close()
	for {
		message, ok := <-p.inbox
		if !ok {
			break
		}
		p.connection.WriteMessage(websocket.TextMessage, message)
	}
} 
 
func initPeer(connection *websocket.Conn, address, port string) *peer {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	key := fmt.Sprintf("%s:%s", address, port)	
	p := &peer{
		connection: connection,
		inbox:		make(chan []byte),
		key:		key,
		address: 	address,
		port: 		port,
	}

	go p.read()
	go p.write()
	
	Peers.v[key] = p		
	return p
}