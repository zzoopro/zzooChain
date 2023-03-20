package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/zzoopro/zzoocoin/utils"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, request *http.Request) {
	// Port:3000 will upgrade the request from :4000 
	openPort := request.URL.Query().Get("openPort")
	ip := utils.Splitter(request.RemoteAddr, ":", 0)	
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}

	connection, err := upgrader.Upgrade(rw, request, nil)	
	utils.HandleErr(err)
	initPeer(connection, ip, openPort)
} 

func AddPeer(address, port, openPort string) {
	// Port:4000 is requesting an upgrade from the Port:3000
	connection, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort),nil)
	utils.HandleErr(err)	
	initPeer(connection, address, port)	
}