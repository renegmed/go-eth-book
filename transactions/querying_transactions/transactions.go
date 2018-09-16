package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	blockNumber := big.NewInt(5671744)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	for idx, tx := range block.Transactions() {
		fmt.Printf("\n%d.\tTx Hash: %v\n", idx+1, tx.Hash().Hex())
		fmt.Printf("\tTx Value: %s\n", tx.Value().String())
		fmt.Printf("\tTx Gas: %d\n", tx.Gas())
		fmt.Printf("\tTx Gas Price: %d\n", tx.GasPrice().Uint64())
		fmt.Printf("\tTx Nonce: %d\n", tx.Nonce())
		fmt.Printf("\tTx Data: %v\n", tx.Data())
		fmt.Printf("\tTx To: %v\n", tx.To().Hex())

		if msg, err := tx.AsMessage(types.HomesteadSigner{}); err == nil {
			fmt.Printf("\tTx Message From: %v\n", msg.From().Hex())
		}

		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\tTx Receipt Status: %d\n", receipt.Status)
	}

	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9")
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	blockHash := common.HexToHash("0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9")
	count, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		log.Fatal(err)
	}

	for idx := uint(0); idx < count; idx++ {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, idx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\tTx In Block Hash: %v\n", tx.Hash().Hex())
	}

	txHash := common.HexToHash("0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2")
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\tTx Hash: %v\n", tx.Hash().Hex())
	fmt.Printf("\tIs Pending: %v\n", isPending)
}
