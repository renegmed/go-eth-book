package main

/**
*  Keystore is a file containing an encrypted wallet private key.
*  Keystores in go-ethereum can only contain one wallet key pair per file.
*
* To generate keystores, first you must invoke NewKeystore giving it the
* directory path to save the keystores.
*
* After, you may generate a new wallet by calling the method NewAccount passing
* it a password for encryption.
*
* Every time you call NewAccount, it will generate a new keystore file on disk
*
 */

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func createKs() {
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	// password := "secret"
	// account, err := ks.NewAccount(password)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(account.Address.Hex())
	fmt.Println(ks)
}

func main() {
	createKs()
}
