package main

import (
	"flag"
	"tcapi"
)

func main() {
	flag.Parse()
	tcapi.Start()
}
