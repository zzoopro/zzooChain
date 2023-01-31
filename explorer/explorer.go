package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/zzoopro/zzoocoin/blockchain"
)

type homeData struct {
	PageTitle string
	Blocks []*blockchain.Block
}

const (	
	templateDir = "explorer/templates/"
)
var (
	templates *template.Template
)

func handleHome(wr http.ResponseWriter, r *http.Request) {
	data := homeData{"Zzoo Coin", nil}
	templates.ExecuteTemplate(wr, "home", data)
}

func handleAdd(wr http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(wr, "add", "Zzoo Coin")
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.Blockchain().AddBlock(data)
		http.Redirect(wr, r, "/", http.StatusPermanentRedirect)
	}	
}

func Start(port int) {
	mux := http.NewServeMux()
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.html"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.html"))

	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/add", handleAdd)
	fmt.Printf("Server running on Port: %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))	
}