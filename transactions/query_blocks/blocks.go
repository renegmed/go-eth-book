package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	// You can call the client's HeaderByNumber to return header information
	// about a block. It'll return the latest block header if you pass nil
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(header.Number.String()) // 6339747

	// Call the client's BlockByNumber method to get the full block. You can
	// read all the contents and metadata of the block such as block number,
	// block timestamp, block hash, block difficulty, as well as the list of
	// transactions and much much more
	blockNumber := big.NewInt(6339747)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Block number: %d\n", block.Number().Uint64())
	fmt.Printf("Time: %d\n", block.Time().Uint64())
	fmt.Printf("Difficulty: %d\n", block.Difficulty().Uint64())
	fmt.Printf("Block Hash: %v\n", block.Hash().Hex())
	fmt.Printf("No. of transactions: %d\n", len(block.Transactions()))

	// Call TransactionCount to return just the count of transactions in a block.
	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Count: %d\n\n", count)
}
