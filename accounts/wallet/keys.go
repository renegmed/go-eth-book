package main

import (
	"fmt"
	"log"

	"github.com/miguelmota/go-ethereum-hdwallet"
)

func main() {
	mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"

	// Returns a new wallet from a BIP-39 mnemonic
	wallet, err := hdwallet.NewFromMnemonic(mnemonic) // returns hdwallet *Wallet
	if err != nil {
		log.Fatal(err)
	}

	// MustParseDerivationPath parses the derivation path in string format into
	// []uint32 but will panic if it can't parse it.
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")

	// Derive implements accounts.Wallet, deriving a new account at the specific
	// derivation path. If pin is set to true, the account will be added to the list
	// of tracked accounts.
	account, err := wallet.Derive(path, false) // returns github.com/ethereum/go-ethereum/accounts Account
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Account address: %s\n", account.Address.Hex())

	// PrivateKeyHex return the ECDSA private key in hex string format of the account.
	privateKey, err := wallet.PrivateKeyHex(account) // returns string
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Private key in hex: %s\n", privateKey)

	// PublicKeyHex return the ECDSA public key in hex string format of the account.
	publicKey, _ := wallet.PublicKeyHex(account) // returns string
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Public key in hex: %s\n", publicKey)

}
