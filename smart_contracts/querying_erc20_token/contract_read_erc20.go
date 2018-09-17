package main

/*

  Querying an ERC20 Token Smart Contract

  Generate the ABI from a solidity source file.

  $ solc --abi erc20.sol -o contracts/

  Now we compile the Go contract file which will include the
  deploy methods because we includes the bin file.

  $ abigen --abi=contracts/erc20.abi --pkg=token --out=contracts/erc20.go

*/
import (
	token "ethereum-go-book/smart_contracts/querying_erc20_token/contracts"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	// Assuming we already have Ethereum client set up as usual,
	// we can now import the new token package into our application
	// and instantiate it. In this example we'll be using the Golem token.
	// Golem (GNT) Address

	tokenAddress := common.HexToAddress("0xa74476443119A942dE498590Fe1f2454d7D4aC0d")
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// We may now call any ERC20 method that we like. For example,
	// we can query the token balance of a user.

	address := common.HexToAddress("0x0536806df512d6cdde913cf95c9886f65b1d3462")
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}

	// We can also read the public variables of the ERC20 smart contract.

	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\tName: %s\n", name)         // "name: Golem Network"
	fmt.Printf("\tSymbol: %s\n", symbol)     // "symbol: GNT"
	fmt.Printf("\tDecimals: %v\n", decimals) // "decimals: 18"

	fmt.Printf("\tWei: %s\n", bal)

	// We can do some simple math to convert the balance into a
	// human readable decimal format.

	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))

	fmt.Printf("\tBalance: %g\n", value)

	/*

		Name: Golem Network Token
		Symbol: GNT
		Decimals: 18
		Wei: 74219874220317902384781652
		Balance: 7.4219874220317902385e+07

		See the same information on etherscan:
		https://etherscan.io/token/0xa74476443119a942de498590fe1f2454d7d4ac0d?a=0x0536806df512d6cdde913cf95c9886f65b1d3462

	*/
}
