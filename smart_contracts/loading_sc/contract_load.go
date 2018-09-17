package main

/*

  Loading a Smart Contract

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

	store "ethereum-go-book/smart_contracts/loading_sc/contracts"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0x35386c483387b87d87eafc0f35504e9539a0b8f2")

	// Once you've compiled your smart contract's ABI to a Go package using the
	// abigen tool, the next step is to call the "New" method, which is in the
	// format New<ContractName>, so in our example if you recall it's going to
	// be NewStore. This initializer method takes in the address of the smart contract
	// and returns a contract instance that you can start interact with it.

	instance, err := store.NewStore(address, client)

	fmt.Println("Contract is loaded")

	_ = instance // not used in this app

	// see installed contract https://rinkeby.etherscan.io/address/0x35386c483387b87d87eafc0f35504e9539a0b8f2
}
