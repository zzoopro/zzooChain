package main

import (
	"github.com/zzoopro/zzoocoin/cli"
	"github.com/zzoopro/zzoocoin/db"
)


func main() {	
	defer db.Close()
	cli.Start()
} 
