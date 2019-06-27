// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SimpleMultiSigABI is the input ABI used to generate the binding from.
// const SimpleMultiSigABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"threshold\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sigV\",\"type\":\"uint8[]\"},{\"name\":\"sigR\",\"type\":\"bytes32[]\"},{\"name\":\"sigS\",\"type\":\"bytes32[]\"},{\"name\":\"destination\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"executor\",\"type\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"execute\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ownersArr\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOwersLength\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"threshold_\",\"type\":\"uint256\"},{\"name\":\"owners_\",\"type\":\"address[]\"},{\"name\":\"chainId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_confirmAddrs\",\"type\":\"address[]\"},{\"indexed\":false,\"name\":\"_destination\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"Execute\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"}]"

// SimpleMultiSigBin is the compiled bytecode used for deploying new contracts.
// const SimpleMultiSigBin = `0x608060405234801561001057600080fd5b50604051610d33380380610d338339818101604052606081101561003357600080fd5b81516020830180519193928301929164010000000081111561005457600080fd5b8201602081018481111561006757600080fd5b815185602082028301116401000000008211171561008457600080fd5b505060209091015181519193509150600a108015906100a4575081518311155b80156100b05750600083115b61011b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f303c7468726573686f6c643c6f776e6572732e6c656e67746800000000000000604482015290519081900360640190fd5b6000805b835181101561022d57816001600160a01b031684828151811061013e57fe5b60200260200101516001600160a01b0316116101bb57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564206f776e6572206f72206e6f7420736f7274656400000000604482015290519081900360640190fd5b6001600260008684815181106101cd57fe5b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060006101000a81548160ff02191690831515021790555083818151811061021857fe5b6020908102919091010151915060010161011f565b50825161024190600390602086019061031a565b505060019290925550604080517fd87cd6ef79d4e2b95e15ce8abf732db51ec771f1ca2edccf22a46c729ac564726020808301919091527fb7a0bfa1b79f2443f4d73ebb9259cddbcd510b18be6fc4da7d1aa7b1786e73e6828401527fc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc6606083015260808201939093523060a08201527f251543af6a222378665a76fe38dbceae4871a070b7fdaf5c6c30cf758dc33cc060c0808301919091528251808303909101815260e090910190915280519101206004556103a6565b82805482825590600052602060002090810192821561036f579160200282015b8281111561036f57825182546001600160a01b0319166001600160a01b0390911617825560209092019160019091019061033a565b5061037b92915061037f565b5090565b6103a391905b8082111561037b5780546001600160a01b0319168155600101610385565b90565b61097e806103b56000396000f3fe6080604052600436106100555760003560e01c80630d8e6e2c1461008d57806342cde4e814610117578063a0ab96531461013e578063aa5df9e21461039b578063affed0e0146103e1578063ca7541ee146103f6575b60408051348152905133917fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c919081900360200190a2005b34801561009957600080fd5b506100a2610421565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100dc5781810151838201526020016100c4565b50505050905090810190601f1680156101095780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561012357600080fd5b5061012c61043f565b60408051918252519081900360200190f35b34801561014a57600080fd5b50610399600480360361010081101561016257600080fd5b810190602081018135600160201b81111561017c57600080fd5b82018360208201111561018e57600080fd5b803590602001918460208302840111600160201b831117156101af57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b8111156101fe57600080fd5b82018360208201111561021057600080fd5b803590602001918460208302840111600160201b8311171561023157600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b81111561028057600080fd5b82018360208201111561029257600080fd5b803590602001918460208302840111600160201b831117156102b357600080fd5b91908080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525092956001600160a01b0385351695602086013595919450925060608101915060400135600160201b81111561031757600080fd5b82018360208201111561032957600080fd5b803590602001918460018302840111600160201b8311171561034a57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550506001600160a01b038335169350505060200135610445565b005b3480156103a757600080fd5b506103c5600480360360208110156103be57600080fd5b5035610916565b604080516001600160a01b039092168252519081900360200190f35b3480156103ed57600080fd5b5061012c61093d565b34801561040257600080fd5b5061040b610943565b6040805160ff9092168252519081900360200190f35b604080518082019091526004815263322e333360e01b602082015290565b60015481565b60015487511461049c576040805162461bcd60e51b815260206004820152601c60248201527f52206c656e206e6f7420657175616c20746f207468726573686f6c6400000000604482015290519081900360640190fd5b855187511480156104ae575087518751145b6104ff576040805162461bcd60e51b815260206004820152601960248201527f6c656e677468206f6620722f732f76206e6f74206d6174636800000000000000604482015290519081900360640190fd5b6001600160a01b03821633148061051d57506001600160a01b038216155b61055f576040805162461bcd60e51b815260206004820152600e60248201526d3bb937b7339032bc32b1baba37b960911b604482015290519081900360640190fd5b825160208085019190912060008054604080517f3ee892349ae4bbe61dce18f95115b5dc02daf49204cc602458cd4c1f540d56d7818701526001600160a01b038b81168284015260608083018c9052608083019690965260a082019390935291871660c083015260e0808301879052815180840390910181526101008301825280519086012060045461190160f01b6101208501526101228401526101428084018290528251808503909101815261016284018084528151918801919091206001548083529788029094016101820190925294919391801561064b578160200160208202803883390190505b50905060005b6001548110156107ac5760006001858f848151811061066c57fe5b60200260200101518f858151811061068057fe5b60200260200101518f868151811061069457fe5b602002602001015160405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa1580156106f3573d6000803e3d6000fd5b505050602060405103519050836001600160a01b0316816001600160a01b031611801561073857506001600160a01b03811660009081526002602052604090205460ff165b61077d576040805162461bcd60e51b815260206004820152601160248201527015995c9a599e481cda59c819985a5b1959607a1b604482015290519081900360640190fd5b8083838151811061078a57fe5b6001600160a01b03909216602092830291909101909101529250600101610651565b5060008054600101815587518190819060208b018c8e8bf1905080610806576040805162461bcd60e51b815260206004820152600b60248201526a6e6f745f7375636365737360a81b604482015290519081900360640190fd5b7f07f4110a9f6788eae6a0b088d9aca06ec3cd9e2c6eae12a1d393d6d041d18c30828b8b8b6040518080602001856001600160a01b03166001600160a01b0316815260200184815260200180602001838103835287818151815260200191508051906020019060200280838360005b8381101561088d578181015183820152602001610875565b50505050905001838103825284818151815260200191508051906020019080838360005b838110156108c95781810151838201526020016108b1565b50505050905090810190601f1680156108f65780820380516001836020036101000a031916815260200191505b50965050505050505060405180910390a150505050505050505050505050565b6003818154811061092357fe5b6000918252602090912001546001600160a01b0316905081565b60005481565b6003549056fea265627a7a72305820fbcf75b0fe8efcca88f733aa6f967b0d614d263723e263240fa03f39f29b0bca64736f6c63430005090032`

// PackedDeploySimpleMultiSig deploys a new Ethereum contract, binding an instance of SimpleMultiSig to it.
func PackedDeploySimpleMultiSig(threshold_ *big.Int, owners_ []common.Address, chainId *big.Int) ([]byte, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleMultiSigABI))
	if err != nil {
		return nil, err
	}
	arguments, err := parsed.Constructor.Inputs.Pack(threshold_, owners_, chainId)
	if err != nil {
		return nil, err
	}
	bytecode := common.FromHex(SimpleMultiSigBin)
	return append(bytecode, arguments...), nil
}

type SimpleMultiSigABIHelper struct {
	abi abi.ABI
}

func NewSimpleMultiSigABIHelper() *SimpleMultiSigABIHelper {
	parsed, _ := abi.JSON(strings.NewReader(SimpleMultiSigABI))
	return &SimpleMultiSigABIHelper{parsed}
}

// GetOwersLength is a free data retrieval call binding the contract method 0xca7541ee.
// Solidity: function getOwersLength() constant returns(uint8)
func (_SimpleMultiSig *SimpleMultiSigABIHelper) PackedGetOwersLength() ([]byte, error) {
	return _SimpleMultiSig.abi.Pack("getOwersLength")
}

// UnpackGetOwersLength is a free data retrieval call binding the contract method 0xca7541ee.
// Solidity: function getOwersLength() constant returns(uint8)
func (_SimpleMultiSig *SimpleMultiSigABIHelper) UnpackGetOwersLength(output []byte) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _SimpleMultiSig.abi.Unpack(out, "getOwersLength", output)
	return *ret0, err
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
// Solidity: function getVersion() constant returns(string)
func (_SimpleMultiSig *SimpleMultiSigABIHelper) PackedGetVersion() ([]byte, error) {
	return _SimpleMultiSig.abi.Pack("getVersion")
}

// UnpackGetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
// Solidity: function getVersion() constant returns(string)
func (_SimpleMultiSig *SimpleMultiSigABIHelper) UnpackGetVersion(output []byte) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SimpleMultiSig.abi.Unpack(out, "getVersion", output)
	return *ret0, err
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
// Solidity: function nonce() constant returns(uint256)
func (_SimpleMultiSig *SimpleMultiSigABIHelper) PackedNonce() ([]byte, error) {
	return _SimpleMultiSig.abi.Pack("nonce")
}

// UnpackNonce is a free data retrieval call binding the contract method 0xaffed0e0.
// Solidity: function nonce() constant returns(uint256)
func (_SimpleMultiSig *SimpleMultiSigABIHelper) UnpackNonce(output []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleMultiSig.abi.Unpack(out, "nonce", output)
	return *ret0, err
}

// OwnersArr is a free data retrieval call binding the contract method 0xaa5df9e2.
// Solidity: function ownersArr(uint256 ) constant returns(address)
func (_SimpleMultiSig *SimpleMultiSigABIHelper) PackedOwnersArr(arg0 *big.Int) ([]byte, error) {
	return _SimpleMultiSig.abi.Pack("ownersArr", arg0)
}

// UnpackOwnersArr is a free data retrieval call binding the contract method 0xaa5df9e2.
// Solidity: function ownersArr(uint256 ) constant returns(address)
func (_SimpleMultiSig *SimpleMultiSigABIHelper) UnpackOwnersArr(output []byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SimpleMultiSig.abi.Unpack(out, "ownersArr", output)
	return *ret0, err
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
// Solidity: function threshold() constant returns(uint256)
func (_SimpleMultiSig *SimpleMultiSigABIHelper) PackedThreshold() ([]byte, error) {
	return _SimpleMultiSig.abi.Pack("threshold")
}

// UnpackThreshold is a free data retrieval call binding the contract method 0x42cde4e8.
// Solidity: function threshold() constant returns(uint256)
func (_SimpleMultiSig *SimpleMultiSigABIHelper) UnpackThreshold(output []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleMultiSig.abi.Unpack(out, "threshold", output)
	return *ret0, err
}

// PackedExecute is a paid mutator transaction binding the contract method 0xa0ab9653.
// Solidity: function execute(uint8[] sigV, bytes32[] sigR, bytes32[] sigS, address destination, uint256 value, bytes data, address executor, uint256 gasLimit) returns()
func (_SimpleMultiSig *SimpleMultiSigABIHelper) PackedExecute(sigV []uint8, sigR [][32]byte, sigS [][32]byte, destination common.Address, value *big.Int, data []byte, executor common.Address, gasLimit *big.Int) ([]byte, error) {
	return _SimpleMultiSig.abi.Pack("execute", sigV, sigR, sigS, destination, value, data, executor, gasLimit)
}
