package main

/*

In the previous section we learned how to sign a piece of data
with a private key in order to generate a signature. Now we'll
learn how to verify the authenticity of the signature.

We need to have 3 things to verify the signature:
	the signature,
	the hash of the original data,
	and the public key of the signer.

With this information we can determine if the private key holder
of the public key pair did indeed sign the message.

*/
import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {

	// First we'll need the public key in bytes format.

	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	// Next we'll need the original data hashed. In the previous lesson we
	// used Keccak-256 to generate the hash, so we'll do the same in order
	// to verify the signature.

	data := []byte("hello")
	hash := crypto.Keccak256Hash(data)
	fmt.Printf("\tData Keccak256 Hash: %v\n", hash.Hex())

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\tEncoded Signature: %v\n", hexutil.Encode(signature))

	// Now assuming we have the signature in bytes format, we can
	// call Ecrecover (elliptic curve signature recover) from the
	// go-ethereum crypto package to retrieve the public key of the
	// signer. This function takes in the hash and signature in bytes format.

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		log.Fatal(err)
	}

	// To verify we simply now have to compare the signature's public key
	// with the expected public key and if they match then the expected
	// public key holder is indeed the signer of the original message.

	matches := bytes.Equal(sigPublicKey, publicKeyBytes)
	fmt.Printf("\tIs it a match? %v\n", matches)

	// There's also the SigToPub method which does the same thing except
	// it'll return the signature's public key in the ECDSA type.

	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
	verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	fmt.Printf("\tIs verified? %v\n", verified)

	/*

		Data Keccak256 Hash: 0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8
		Encode Signature: 0x789a80053e4927d0a898db8e065e948f5cf086e32f9ccaa54c1908e22ac430c62621578113ddbb62d509bf6049b8fb544ab06d36f916685a2eb8e57ffadde02301
		It is a match? true
		Is verified? true

	*/

}
