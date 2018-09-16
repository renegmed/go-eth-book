package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*

	This section will walk you through on how to transfer ERC-20 tokens

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

	value := big.NewInt(0) // in wei (0 eth)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Assuming you've already connected a client, loaded your private key,
	// and configured the gas price, the next step is to set the data field of
	// the transaction

	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")

	// We'll need to figure out the signature of the smart contract function
	// we'll be calling, along with the inputs that the function will be receiving.
	// We then take the keccak-256 hash of the function signature to retreive the
	// method ID which is the first 8 characters (4 bytes).
	//
	// Afterwards we append the address we're sending to, as well append the amount
	// of tokens we're transferring.
	//
	// These inputs will need to be 256 bits long (32 bytes) and left padded. The method ID
	// is not padded.
	tokenAddress := common.HexToAddress("0x28b149020d2152179873ec60bed6bf7cd705775d")

	// The function signature will be the name of the transfer function, which is transfer
	// in the ERC-20 specification, and the argument types. The first argument type is
	// address (receiver of the tokens) and the second type is uint256 (amount of tokens to send).
	// There should be no spaces or argument names. We'll also need it as a byte slice.

	transferFnSignature := []byte("transfer(address,uint256")

	// We'll now import the crypto/sha3 package from go-ethereum to generate the Keccak256
	// hash of the function signature. We then take only the first 4 bytes to have the method ID.

	hash := sha3.NewKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Printf("\thexuitil.Encode(methodID): %v\n", hexutil.Encode(methodID))

	// Next we'll need to left pad 32 bytes the address we're sending tokens to.

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Printf("\tSend to address - hexutil.Encode(paddedAmount): %v\n", hexutil.Encode(paddedAddress)) // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d

	// Next we determine how many tokens we want to send, in this case it'll be 1,000 tokens
	// which will need to be formatted to wei in a big.Int

	amount := new(big.Int)
	amount.SetString("1000000000000000000000", 10) // 1000 tokens

	// Left padding to 32 bytes will also be required for the amount.

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Printf("\thexutil.Encode(paddedAmount)): %v\n", paddedAmount)

	// Now we simply concanate the method ID, padded address, and padded amount to a byte slice
	// that will be our data field.

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// The gas limit will depend on the size of the transaction data and computational steps
	// that the smart contract has to perform. Fortunately the client provides the method
	// EstimateGas which is able to estimate the gas for us. This function takes a CallMsg struct
	// from the ethereum package where we specify the data and to address. It'll return the
	// estimated gas limit units we'll be needing for generating the complete transaction.

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\tEstimated Gas limit: %v\n", gasLimit) // 23256

	// Next thing we need to do is generate the transaction type, similar to what you've seen
	// in the transfer ETH section, EXCEPT the 'to' field will be the token smart contract address.
	// This is a gotcha that confuses people. We must also include the value field which will
	// be 0 ETH, and the data bytes that we just generated.

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	// The next step is to sign the transaction with the private key of the sender.

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// The SignTx method requires the EIP155 signer, which we derive the chain ID from the client.

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// And finally broadcast the transaction.

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\ttx send: %s\n", signedTx.Hash().Hex()) // tx sent: 0xa56316b637a94c4cc0331c73ef26389d6c097506d581073f927275e7a6ece0bc

}
