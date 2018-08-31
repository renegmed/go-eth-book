package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	// Reading the balance of account.

	account := common.HexToAddress("0x71c7656ec7ab88b098defb751b7401b5f6d8976f")

	// Setting nil as the block number will return the latest balance
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance) // 25943679348360745848

	// Reading the balance on a particular block

	blockNumber := big.NewInt(6123635) //5532993) // or 6123635
	balanceAt, err := client.BalanceAt(context.Background(), account, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Balance of block", blockNumber, ":", balanceAt)
}
