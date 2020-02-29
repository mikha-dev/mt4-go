package main

import (
	"flag"
	"mtdealerapi"
)

func main() {
	flag.Parse()
	mtdealerapi.Start()
}
