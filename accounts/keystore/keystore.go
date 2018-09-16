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
	"io/ioutil"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func createKs() {
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	password := "secret"
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(account.Address.Hex())
	fmt.Println(account.Address.Hex()) // 0x9fd3F47d11d9F7454C4D7A1D222A1839F9A6b21b
}

func importKs() {
	file := "./tmp/UTC--2018-07-04T09-58-30.122808598Z--20f8d42fb0f667f2e53930fed426f225752453b3"
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	password := "secret"
	account, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex()) // 0x20F8D42FB0F667F2E53930fed426f225752453b3

	if err := os.Remove(file); err != nil {
		log.Fatal(err)
	}
}

func main() {
	createKs()
	// importKs()
}
