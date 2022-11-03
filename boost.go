package main

import (
	"context"
	"flag"
	"log"
)

var ipc string

func init() {
	flag.StringVar(&ipc, "ipc", "/opt/geth/geth.ipc", "geth ipc path")

}

func addStaticPeerFrom(ctx context.Context, peerChan chan string) {
	//ipcCli, err := rpc.DialIPC(context.Background(), ipc)
	//if err != nil {
	//	log.Printf("rpc dial err: %v", err)
	//	return
	//}

	for {
		select {
		case <-ctx.Done():
			return
		case nodeURL, ok := <-peerChan:
			if !ok {
				return
			}
			var resultBool bool
			//if err := ipcCli.Call(
			//	&resultBool,
			//	"admin_addPeer",
			//	nodeURL,
			//); err != nil {
			//	panic(err)
			//}

			log.Printf("add node %s result: %v", nodeURL, resultBool)
		}
	}

}
