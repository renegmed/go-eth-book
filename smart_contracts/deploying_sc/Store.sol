pragma solidity ^0.4.24;

/*

  Generate the ABI from a solidity source file.
  
  $ solc --abi Store.sol -o contracts/

  In order to deploy a smart contract from Go, we also need to compile
  the solidity smart contract to EVM bytecode. The EVM bytecode is what 
  will be sent in the data field of the transaction. The bin file is 
  required for generating the deploy methods on the Go contract file.

  $ solc --bin Store.sol -o contracts/

  Now we compile the Go contract file which will include the
  deploy methods because we includes the bin file.

  $ abigen --bin=contracts/Store.bin --abi=contracts/Store.abi --pkg=store --out=contracts/Store.go

*/
contract Store {

    event ItemSet(bytes32 key, bytes32 value);

    string public version;
    mapping (bytes32 => bytes32) public items;

    constructor(string _version) public {
        version = _version;
    }

    function setItem(bytes32 key, bytes32 value) external {
        items[key] = value;
        emit ItemSet(key, value);
    }
}