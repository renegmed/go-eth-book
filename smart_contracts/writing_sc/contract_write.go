package main

/*

  Writing to a Smart Contract

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

	store "ethereum-go-book/smart_contracts/writing_sc/contracts"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	// Writing to a smart contract requires us to sign the sign transaction with our private key.

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

	// We'll also need to figure the nonce and gas price.

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Next we create a new keyed transactor which takes in the private key.

	auth := bind.NewKeyedTransactor(privateKey)

	// Then we need to set the standard transaction options attached to the keyed transactor.

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	// Now we load an instance of the smart contract. If you recall in the previous sections
	// we create a contract called Store and generated a Go package file using the abigen tool.
	// To initialize it we just invoke the New method of the contract package and give the
	// smart contract address and the ethclient, which returns a contract instance that we can use.

	address := common.HexToAddress("0x35386c483387b87d87eafc0f35504e9539a0b8f2")
	instance, err := store.NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}

	// The smart contract that we created has an external method called SetItem which
	// takes in two arguments (key, value) in the form of solidity bytes32. This means
	// that the Go contract package requires us to pass a byte array of length 32 bytes.
	// Invoking the SetItem method requires us to pass the auth object we created earlier.
	// Behind the scenes this method will encode this function call with it's arguments,
	// set it as the data property of the transaction, and sign it with the private key,
	// bind.NewKeyedTransactor(privateKey).
	// The result will be a signed transaction object.

	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], []byte("foo"))
	copy(value[:], []byte("bar"))

	tx, err := instance.SetItem(auth, key, value)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\ttx sent: %s\n", tx.Hash().Hex())

	result, err := instance.Items(nil, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\tResult: %s\n", string(result[:])) // "bar"

	/*

		tx sent: 0x58043e47acb14787c87c2061be07e1ff726277dad61519d48612a87846d74657
		Result: bar

		NOTE: need to wait for transaction to process, thus initial result will be empty. Need to
		      run the app commenting out the instance.SetItem() part
	*/
}
