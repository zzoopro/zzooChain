package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zzoopro/zzoocoin/blockchain"
	"github.com/zzoopro/zzoocoin/utils"
)

type urlText string

func (u urlText) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost:%d%s", port, u)
	return []byte(url), nil
}

type errorResponse struct {
	ErrorMessage string `json:"error_message"`
}

type urlDescription struct {
	URL 		urlText `json:"url"`
	Method 		string `json:"method"`
	Payload 	string `json:"payload,omitempty"`
	Description string `json:"description"`	
}

type BalanceResponse struct {
	Address string 	`json:"address"`
	Balance  int	`json:"balance"`
}

type addTxPayload struct {
	To 		string 	`json:"to"`
	Amount 	int		`json:"amount"`
}

var port int

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL: urlText("/"),
			Method: "GET",
			Description: "See documentation",
		},
		{
			URL: urlText("/status"),
			Method: "GET",
			Description: "See blockchain status",
		},
		{
			URL: urlText("/blocks"),
			Method: "GET",
			Description: "See all block",			
		},
		{
			URL: urlText("/blocks"),
			Method: "POST",
			Description: "Add a block",			
		},
		{
			URL: urlText("/blocks/{hash}"),
			Method: "GET",
			Description: "See a block",			
		},
		{
			URL: urlText("/balance/{address}"),
			Method: "GET",
			Description: "get balance by address",
		},
		{
			URL: urlText("/transaction"),
			Method: "POST",
			Description: "add transaction to mempool",
		},
	}	
	json.NewEncoder(rw).Encode(data)
}

func handleBlocks(rw http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":		
		json.NewEncoder(rw).Encode(blockchain.Blockchain().Blocks())
	case "POST":
		blockchain.Blockchain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
	}
}

func handleBlock(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	hash := vars["hash"]	
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)	
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}	 
}

func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func handleStatus(rw http.ResponseWriter, request *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.Blockchain())
}

func handleBalance(rw http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	address := vars["address"]	
	isTotal := request.URL.Query().Get("total")
	
	switch isTotal {
	case "true":
		amount := blockchain.Blockchain().BalanceByAddress(address)
		utils.HandleErr(json.NewEncoder(rw).Encode(BalanceResponse{address, amount}))
	default:
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Blockchain().TxOutputsByAddress(address)))			
	}	
}

func handleMempool(rw http.ResponseWriter, request *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
}

func handleTransaction(rw http.ResponseWriter, request *http.Request) {
	var payload addTxPayload 
	utils.HandleErr(json.NewDecoder(request.Body).Decode(&payload))
	fmt.Println(payload)
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		json.NewEncoder(rw).Encode(errorResponse{"Not enough funds."})
	}
	rw.WriteHeader(http.StatusCreated)
}

func Start(portNum int){
	router := mux.NewRouter()
	port = portNum	

	router.Use(contentTypeMiddleware)
	
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/mempool", handleMempool).Methods("GET")
	router.HandleFunc("/status", handleStatus).Methods("GET")
	router.HandleFunc("/blocks", handleBlocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", handleBlock).Methods("GET")
	router.HandleFunc("/balance/{address}", handleBalance).Methods("GET")	
	router.HandleFunc("/transaction", handleTransaction).Methods("POST")	

	fmt.Printf("Server listening on http://localhost%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))	
}