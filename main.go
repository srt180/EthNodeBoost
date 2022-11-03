package main

import (
	"context"
	"flag"
)

func main() {
	flag.Parse()

	nodeChan := queryNodes()
	addStaticPeerFrom(context.Background(), nodeChan)
}
