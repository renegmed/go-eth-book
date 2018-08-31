package main

import (
	"context"
	"fmt"
	"log"
	"math"
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
	accountStr := "0x71c7656ec7ab88b098defb751b7401b5f6d8976f"
	account := common.HexToAddress(accountStr)

	// Setting nil as the block number will return the latest balance
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance) // 25943679348360745848

	// Reading the balance on a particular block

	blockNumber := big.NewInt(5532993) // or 6123635
	balanceAt, err := client.BalanceAt(context.Background(), account, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Balance of block", blockNumber, ":", balanceAt)

	// With conversion to Ether units since 1 ether = 10^18 weis
	// we need to use Go's math and  math/big packages

	fbalance := new(big.Float)
	fbalance.SetString(balanceAt.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println(ethValue) // 25.729324269165216041

	// Pending Balance
	//   Sometimes there is pending account balance. For example after submitting or
	//   waiting for a transaction to be confirmed, pending account balance will be created
	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	fmt.Println("Pending Balance of account", accountStr, ":", pendingBalance) // 25943679348360745848
}
