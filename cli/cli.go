package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/zzoopro/zzoocoin/api"
	"github.com/zzoopro/zzoocoin/explorer"
)


func usage() {
	fmt.Printf("Welcome to zzoocoin.\n")
	fmt.Printf("Please use the following flags:\n")
	fmt.Printf("-mode=html :  What do you want?\n")
	fmt.Printf("-port=4000 :  Set port of the server.\n")	
	runtime.Goexit()
}

func Start() {
	if len(os.Args) < 2 {
		usage()
	}
	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "api", "Choose between 'html' and 'api'")
	
	flag.Parse()
	switch *mode {
		case "html":
			explorer.Start(*port)
		case "api":
			api.Start(*port)
		default:
			usage()
	}
}