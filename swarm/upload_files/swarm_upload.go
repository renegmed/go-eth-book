package main

/*

In the previous section we setup a swarm node running as
a daemon on port 8500. Now import the swarm package go-ethereum
swarm/api/client

Commands

geth account new
export BZZKEY=970ef9790b54425bea2c02e25cab01e48cf92573
swarm --bzzaccount $BZZKEY

*/
import (
	"fmt"
	"log"

	bzzclient "github.com/ethereum/go-ethereum/swarm/api/client"
)

func main() {

	// Invoke NewClient function passing it the swarm daemon url.

	client := bzzclient.NewClient("http://127.0.0.1:8500")

	// Create an example text file hello.txt with the content
	// hello world. We'll be uploading this to swarm.

	// n our Go application we'll open the file we just created
	// using Open from the client package. This function will
	// return a File type which represents a file in a swarm
	// manifest and is used for uploading and downloading content
	// to and from swarm.

	file, err := bzzclient.Open("hello.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Now we can invoke the Upload function from our client instance
	// giving it the file object. The second argument is an optional
	// existing manifest string to add the file to, otherwise it'll
	// create on for us. The third argument is if we want our data
	// to be encrypted.

	// The hash returned is the swarm hash of a manifest that contains
	// the hello.txt file as its only entry. So by default both the
	// primary content and the manifest is uploaded. The manifest
	// makes sure you could retrieve the file with the correct mime type.

	manifestHash, err := client.Upload(file, "", false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(manifestHash)

	/*

		56399b64001ed29d65a7823f61f22eda8b31d66a3224052bea489698b884ac1a

		Now we can access our file at
		bzz://56399b64001ed29d65a7823f61f22eda8b31d66a3224052bea489698b884ac1a
	*/
}
