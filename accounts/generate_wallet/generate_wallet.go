package main

/*
* Example:
*
* $  go run main.go generate_wallet.go
*
 */
import (
	"crypto/ecdsa"

	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

func main() {

	// To generate a new wallet, call crypto.GenerateKey() to generate
	// a random private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	// Convert privateKey into bytes using crypto/edcsa package and method FromECDSA
	privateKeyBytes := crypto.FromECDSA(privateKey)

	// Print out privateKeyBytes using hexutil package Encode method which takes a
	// byte slice. Then strip the 0x after it's hex encoded
	fmt.Println(hexutil.Encode(privateKeyBytes))
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:])

	// Note that private keys are used for signing transactions and is to be treated
	// like a password

	// Generating Public Key

	// Since public key is derived from the private key, we can use crypto Public() method
	// to generate public key.

	publicKey := privateKey.Public()
	// returns crypto.PublicKey
	//type PublicKey struct {
	//	elliptic.Curve
	//	X, Y *big.Int
	//}

	// Convert to hex. Strip off the 0x and first 2 characters 04
	// which is always the EC prefix and is not required
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes))
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:])

	// Generate public address.
	// crypto.PubkeyToAddress method accepts an ECDSA public key
	// and returns public address

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address)

	// The public address is simply the Keccak-256 hash of the public key.
	// We take the last 40 characters (20 bytes) and prefix it with 0x.

	hash := sha3.NewKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)))
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:]))

}
