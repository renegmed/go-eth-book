package main

/*

In the previous section we uploaded a hello.txt file to swarm and in return we got a manifest hash.

manifestHash = "56399b64001ed29d65a7823f61f22eda8b31d66a3224052bea489698b884ac1a"

*/

import (
	"fmt"
	"io/ioutil"
	"log"

	bzzclient "github.com/ethereum/go-ethereum/swarm/api/client"
)

func main() {

	client := bzzclient.NewClient("http://127.0.0.1:8500")
	manifestHash := "56399b64001ed29d65a7823f61f22eda8b31d66a3224052bea489698b884ac1a"

	// inspect the manifest by downloading it first by calling DownloadManfest.

	manifest, isEncrypted, err := client.DownloadManifest(manifestHash)
	if err != nil {
		log.Fatal(err)
	}

	// We can iterate over the manifest entries and see what the
	// content-type, size, and content hash are.

	fmt.Printf("\tIs encrypted: %t\n", isEncrypted)

	for _, entry := range manifest.Entries {
		fmt.Printf("\tEntry Hash: %v\n", entry.Hash)
		fmt.Printf("\tContent Type: %s\n", entry.ContentType)
		fmt.Printf("\tSize: %d\n", entry.Size)
		fmt.Printf("\tPath: %s\n", entry.Path)
	}

	// If you're familiar with swarm urls, they're in the
	// format bzz:/<hash>/<path>, so in order to download
	// the file we specify the manifest hash and path. The path
	// in this case is an empty string. We pass this data to
	// the Download function and get back a file object.

	file, err := client.Download(manifestHash, "")
	if err != nil {
		log.Fatal(err)
	}

	// If you're familiar with swarm urls, they're in the
	// format bzz:/<hash>/<path>, so in order to download
	// the file we specify the manifest hash and path. The
	// path in this case is an empty string. We pass this
	// data to the Download function and get back a file object.

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\tFile Content:\n\t\t%v\n", string(content))
}
