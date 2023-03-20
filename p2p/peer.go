package p2p

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var Peers map[string]*peer = make(map[string]*peer)

type peer struct {
	connection *websocket.Conn
}
 
func initPeer(connection *websocket.Conn, address, port string) *peer {
	p := &peer{
		connection,
	}
	key := fmt.Sprintf("%s:%s", address, port)
	Peers[key] = p
	return p
}