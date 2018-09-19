package main

/*
	The components for generating a signature are:
		the signers private key,
		and the hash of the data that will be signed.

	Any hashing algorithm may be used as long as the output is
	32 bytes. We'll be using Keccak-256 as the hashing algorithm
	which is what Ethereum prefers to use.

*/
import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {

	// First we'll load private key.

	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}

	// Next we'll take the Keccak-256 of the data that we wish to sign,
	// in this case it'll be the word hello. The go-ethereum crypto package
	// provides a handy Keccak256Hash method for doing this.

	data := []byte("hello")
	hash := crypto.Keccak256Hash(data)
	fmt.Printf("\tData Keccak256Hash: %v\n", hash.Hex())

	// Finally we sign the hash with our private key,
	// which gives us the signature.

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\tHash Signature: %v\n", hexutil.Encode(signature))

	/*

		Data Keccak256Hash: 0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8

		Hash Signature:
		  0x789a80053e4927d0a898db8e065e948f5cf086e32f9ccaa54c1908e22ac430c62621578113ddbb62d509bf6049b8fb544ab06d36f916685a2eb8e57ffadde02301


	*/

}
