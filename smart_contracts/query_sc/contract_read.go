package main

/*

  Querying a Smart Contract

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
	"fmt"
	"log"

	store "ethereum-go-book/smart_contracts/query_sc/contracts"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0x35386c483387b87d87eafc0f35504e9539a0b8f2")
	instance, err := store.NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}

	// In the previous section (Deploying a Smart Contract) we learned how to
	// initialize a contract instance in our Go application. Now we're going
	// to read the smart contract using the provided methods by the new contract
	// instance. If you recall we had a global variable named version in our
	// contract that was set during deployment. Because it's public that means
	// that they'll be a getter function automatically created for us. Constant
	// and view functions also accept bind.CallOpts as the first argument.
	// To learn about what options you can pass checkout the type's documentation
	// but usually this is set to nil.
	version, err := instance.Version(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(version) // "1.0"
}
