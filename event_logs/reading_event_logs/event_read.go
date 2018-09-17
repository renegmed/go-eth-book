package main

/*

  Reading Event Logs

  A smart contract may optionally emit "events" which get stored a logs as part of the
  transaction receipt. Reading these events are pretty simple.

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
	"fmt"
	"log"
	"math/big"
	"strings"

	store "ethereum-go-book/event_logs/reading_event_logs/contracts"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	client, err := ethclient.Dial("wss://rinkeby.infura.io/ws")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress("0x147B8eb97fD247D06C4006D269c90C1908Fb5D54") // contract address
	//contractAddress := common.HexToAddress("0x35386c483387B87d87EaFC0f35504E9539a0B8F2")

	// First we need to construct a filter query. We import the FilterQuery struct from
	// the go-ethereum package and initialize it with filter options. We tell it the
	// range of blocks that we want to filter through and specify the contract address
	// to read this logs from. In this example we'll be reading all the logs from a
	// particular block, from the smart contract we created in the smart contract sections.

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(2394201),
		ToBlock:   big.NewInt(2394201),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	// The next is step is to call FilterLogs from the ethclient that
	// takes in our query and will return all the matching event logs.

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	// All the logs returned will be ABI encoded so by themselves they
	// won't be very readable. In order to decode the logs we'll need to import
	// our smart contract ABI. To do that, we import our compiled smart contract
	// Go package which will contain an external property in the name
	// format <ContractName>ABI containing our ABI. Afterwards we use
	// the abi.JSON function from the go-ethereum accounts/abi go-ethereum package
	// to return a parsed ABI interface that we can use in our Go application.

	contractAbi, err := abi.JSON(strings.NewReader(string(store.StoreABI)))
	if err != nil {
		log.Fatal(err)
	}

	// Now we can interate through the logs and decode them into a type we can use.
	// If you recall the logs that our sample contract emitted were of
	// type bytes32 in Solidity, so the equivalent in Go would be [32]byte.
	// We can create an anonymous struct with these types and pass a pointer as
	// the first argument to the Unpack function of the parsed ABI interface
	// to decode the raw log data. The second argument is the name of the event
	// we're trying to decode and the last argument is the encoded log data.

	for _, vLog := range logs {

		// The log struct contains additional information such as the block hash,
		// block number, and transaction hash.

		fmt.Printf("\tLog.BlockHash.Hex() %v\n", vLog.BlockHash.Hex())
		fmt.Printf("\tLog.BlockNumber: %d\n", vLog.BlockNumber)
		fmt.Printf("\tLog.TxHash.Hex(): %v\n", vLog.TxHash.Hex())

		event := struct {
			Key   [32]byte
			Value [32]byte
		}{}

		err := contractAbi.Unpack(&event, "ItemSet", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\tEvent Key: %s\n", string(event.Key[:]))
		fmt.Printf("\tEvent Value: %s\n", string(event.Value[:]))

		// If your solidity event contains indexed event types, then they
		// become a topic rather than part of the data property of the log.
		// In solidity you may only have up to 4 topics but only 3 indexed
		// event types. The first topic is always the signature of the event.
		// Our example contract didn't contain indexed events, but if it did
		// this is how to read the event topics.

		var topics [4]string
		for i := range vLog.Topics {
			topics[i] = vLog.Topics[i].Hex()
		}

		fmt.Printf("\tTopics: %v\n\n", topics[0])
	}

	eventSignature := []byte("ItemSet(bytes32,bytes32)")
	hash := crypto.Keccak256Hash(eventSignature)
	fmt.Printf("\tEvent Signature Hash: %v\n", hash.Hex())
}
