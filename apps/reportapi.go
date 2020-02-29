package main

import (
	"flag"
	"mtreportapi"
)

func main() {
	flag.Parse()
	mtreportapi.Start()
}
