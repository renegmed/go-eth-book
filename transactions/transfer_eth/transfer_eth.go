package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*
	In this lesson you'll learn how to transfer ETH from one account to another account.
	If you're already familar with Ethereum then you know that a transaction consists of
	the amount of ether you're transferring, the gas limit, the gas price, a nonce, the
	receiving address, and optionally data. The transaction must be signed with the private
	key of the sender before it's broadcasted to the network.
*/
func main() {
	client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// The next step is to set the amount of ETH that we'll be transferring. However
	// we must convert ether to wei since that's what the Ethereum blockchain uses.
	// Ether supports up to 18 decimal places so 1 ETH is 1 plus 18 zeros. Here's a
	// little tool to help you convert between ETH and wei: https://etherconverter.online
	value := big.NewInt(1000000000000000000) // in wei (1 eth)

	// The gas limit for a standard ETH transfer is 21000 units.
	gasLimit := uint64(21000) // in units

	// The gas price must be set in wei. At the time of this writing, a gas price that
	// will get your transaction included pretty fast in a block is 30 gwei.

	// However, gas prices are always fluctuating based on market demand and what users
	// are willing to pay, so hardcoding a gas price is sometimes not ideal. The go-ethereum
	// client provides the SuggestGasPrice function for getting the average gas price based
	// on x number of previous blocks.
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// We figure out who we're sending the ETH to
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	var data []byte

	// Now we can finally generate our unsigned ethereum transaction by importing
	// the go-ethereum core/types package and invoking NewTransaction which takes
	// in the nonce, to address, value, gas limit, gas price, and optional data.
	// The data field is nil for just sending ETH. We'll be using the data field when
	// it comes to interacting with smart contracts.
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Now we are finally ready to broadcast the transaction to the entire network by
	// calling SendTransaction on the client which takes in the signed transaction.
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}
