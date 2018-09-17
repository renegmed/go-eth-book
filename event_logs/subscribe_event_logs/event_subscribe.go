package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	// First thing we need to do in order to subscribe to event logs
	// is dial to a websocket enabled Ethereum client. Fortunately
	// for us, Infura supports websockets.

	client, err := ethclient.Dial("wss://rinkeby.infura.io/ws")
	if err != nil {
		log.Fatal(err)
	}

	// The next step is to create a filter query. In this example
	// we'll be reading all events coming from the example contract
	// that we've created in the previous lessons.

	contractAddress := common.HexToAddress("0x147B8eb97fD247D06C4006D269c90C1908Fb5D54")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	// The way we'll be receiving events is through a Go channel.
	// Let's create one with type of Log from the go-ethereum core/types package.

	logs := make(chan types.Log)

	// Now all we have to do is subscribe by calling SubscribeFilterLogs
	// from the client, which takes in the query options and the output channel.
	// This will return a subscription struct containing unsubscribe and error methods.

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	// Finally all we have to do is setup an continuous loop with a
	// select statement to read in either new log events or the subscription error.

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Printf("\tLog: %v\n", vLog) // pointer to event log
		}
	}
}
