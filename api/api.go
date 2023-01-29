package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

type addBlockBody struct {
	Data string
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
			URL: urlText("/blocks"),
			Method: "GET",
			Description: "See All Block",			
		},
		{
			URL: urlText("/blocks"),
			Method: "POST",
			Description: "Add a block",
			Payload: "data:string",
		},
		{
			URL: urlText("/blocks/{height}"),
			Method: "GET",
			Description: "See a block",			
		},
	}	
	json.NewEncoder(rw).Encode(data)
}

func handleBlocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":		
		json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())
	case "POST":
		var body addBlockBody
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&body))
		blockchain.GetBlockchain().AddBlock(body.Data)
		rw.WriteHeader(http.StatusCreated)
	}
}

func handleBlock(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	height, err := strconv.Atoi(vars["height"])
	utils.HandleErr(err)
	block, err := blockchain.GetBlockchain().FindBlock(height)
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

func Start(portNum int){
	router := mux.NewRouter()
	port = portNum	

	router.Use(contentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", handleBlocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{height:[0-9]+}", handleBlock).Methods("GET")

	fmt.Printf("Server listening on http://localhost%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))	
}