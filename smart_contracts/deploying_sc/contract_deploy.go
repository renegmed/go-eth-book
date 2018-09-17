package main

/*

  Deploying a Smart Contract

  Generate the ABI from a solidity source file.

  $ solc --abi Store.sol -o contracts/

  In order to deploy a smart contract from Go, we also need to compile
  the solidity smart contract to EVM bytecode. The EVM bytecode is what
  will be sent in the data field of the transaction. The bin file is
  required for generating the deploy methods on the Go contract file.

  $ solc --bin Store.sol -o contracts/

  Now we compile the Go contract file which will include the
  deploy methods because we includes the bin file.

  $ abigen --bin=contracts/Store.bin --abi=contracts/Store.abi --pkg=store --out=contracts/Store.go

*/
import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	store "ethereum-go-book/smart_contracts/deploying_sc/contracts" // for demo
)

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

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Assuming you've imported the newly created Go package file generated
	// from abigen, and set the ethclient, loaded your private key, the next step
	// is to create a keyed transactor.

	// First import the accounts/abi/bind package from go-ethereum and then invoke
	// NewKeyedTransactor passing in the private key. Afterwards set the usual
	// properties such as the nonce, gas price, gas limit, and ETH value.

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	// The deploy function takes in the keyed transactor, the ethclient, and any
	// input arguments that the smart contract constructor might takes in. We've
	// set our smart contract to take in a string argument for the version. This
	// function will return the Ethereum address of the newly deployed contract,
	// the transaction object, the contract instance so that we can start
	// interacting with, and the error if any.

	input := "1.0"
	address, tx, instance, err := store.DeployStore(auth, client, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\tAddress: %v\n", address.Hex())
	fmt.Printf("\tTx Hash: %v\n", tx.Hash().Hex())

	_ = instance

	/*
		RESULT:

		Address: 0x35386c483387B87d87EaFC0f35504E9539a0B8F2
		Tx Hash: 0x93c889c233f28ff9af72ab614a90dfb12a5b3758b212c95d88a997543920ae4f

		https://rinkeby.etherscan.io/address/0x35386c483387b87d87eafc0f35504e9539a0b8f2

		https://rinkeby.etherscan.io/tx/0x93c889c233f28ff9af72ab614a90dfb12a5b3758b212c95d88a997543920ae4f
	*/
}
