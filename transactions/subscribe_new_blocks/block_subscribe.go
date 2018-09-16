package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*

  In this section we'll go over how to set up a subscription
  to get events when there is a new block mined.

*/

func main() {
	// First thing is we need an Ethereum provider that
	// supports RPC over websockets. In this example
	// we'll use the infura websocket endpoint.

	client, err := ethclient.Dial("wss://ropsten.infura.io/ws")
	if err != nil {
		log.Fatal(err)
	}

	// Next we'll create a new channel that will be receiving
	// the latest block headers.

	headers := make(chan *types.Header)

	// Now we call the client's SubscribeNewHead method which
	// takes in the headers channel we just created, which will
	// return a subscription object.

	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	// The subscription will push new block headers to our channel
	// so we'll use a select statement to listen for new messages.
	// The subscription object also contains an error channel that
	// will send a message in case of a failure with the subscription.

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Printf("\theader.Hash().Hex(): %v\n", header.Hash().Hex())

			// To get the full contents of the block, we can pass the
			// block header hash to the client's BlockByHash function.

			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("\tblock.Hash().Hex(): %v\n", block.Hash().Hex())
			fmt.Printf("\tblock.Number().Uint64(): %v\n", block.Number().Uint64())
			fmt.Printf("\tblock.Time().Uint64(): %v\n", block.Time().Uint64())
			fmt.Printf("\tblock.Nonce(): %v\n", block.Nonce())
			fmt.Printf("\tBlock Transactions Count: %d\n\n", len(block.Transactions()))
		}
	}
}

/*
 Sample result:
 	header.Hash().Hex(): 0xe5924cc552484e1b7a0094c6497ead04a0f6d076d77047a89421482f89b3952e
	block.Hash().Hex(): 0xe5924cc552484e1b7a0094c6497ead04a0f6d076d77047a89421482f89b3952e
	block.Number().Uint64(): 4050854
	block.Time().Uint64(): 1537138784
	block.Nonce(): 6488584796222382448
	Block Transactions Count: 12

	header.Hash().Hex(): 0x55349543ebd10fa7ddb449caa50cc827344396eb188fa7b6b35f6cdb5f49e293
	block.Hash().Hex(): 0x55349543ebd10fa7ddb449caa50cc827344396eb188fa7b6b35f6cdb5f49e293
	block.Number().Uint64(): 4050855
	block.Time().Uint64(): 1537138814
	block.Nonce(): 13850456764717066103
	Block Transactions Count: 40

	header.Hash().Hex(): 0xf919b3923fe591d05a4cd60273032d016ca661c845cb059e4d73b2b7cd48d62c
	block.Hash().Hex(): 0xf919b3923fe591d05a4cd60273032d016ca661c845cb059e4d73b2b7cd48d62c
	block.Number().Uint64(): 4050856
	block.Time().Uint64(): 1537138830
	block.Nonce(): 2769471145547757776
	Block Transactions Count: 8

	header.Hash().Hex(): 0x0c77989dadf52142b780cbe69568edb24b23fe42168ba723559ee60b40076888
	block.Hash().Hex(): 0x0c77989dadf52142b780cbe69568edb24b23fe42168ba723559ee60b40076888
	block.Number().Uint64(): 4050857
	block.Time().Uint64(): 1537138848
	block.Nonce(): 7940714359153374125
	Block Transactions Count: 9

*/
