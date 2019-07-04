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
const SimpleMultiSigABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"threshold\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"sigV\",\"type\":\"uint8[]\"},{\"name\":\"sigR\",\"type\":\"bytes32[]\"},{\"name\":\"sigS\",\"type\":\"bytes32[]\"},{\"name\":\"destination\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"executor\",\"type\":\"address\"},{\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"execute\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ownersArr\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOwersLength\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"threshold_\",\"type\":\"uint256\"},{\"name\":\"owners_\",\"type\":\"address[]\"},{\"name\":\"chainId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_confirmAddrs\",\"type\":\"address[]\"},{\"indexed\":false,\"name\":\"_destination\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"Execute\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"}]"

// SimpleMultiSigBin is the compiled bytecode used for deploying new contracts.
const SimpleMultiSigBin = `0x608060405234801561001057600080fd5b50604051610d33380380610d338339818101604052606081101561003357600080fd5b81516020830180519193928301929164010000000081111561005457600080fd5b8201602081018481111561006757600080fd5b815185602082028301116401000000008211171561008457600080fd5b505060209091015181519193509150600a108015906100a4575081518311155b80156100b05750600083115b61011b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f303c7468726573686f6c643c6f776e6572732e6c656e67746800000000000000604482015290519081900360640190fd5b6000805b835181101561022d57816001600160a01b031684828151811061013e57fe5b60200260200101516001600160a01b0316116101bb57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f7265706561746564206f776e6572206f72206e6f7420736f7274656400000000604482015290519081900360640190fd5b6001600260008684815181106101cd57fe5b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060006101000a81548160ff02191690831515021790555083818151811061021857fe5b6020908102919091010151915060010161011f565b50825161024190600390602086019061031a565b505060019290925550604080517fd87cd6ef79d4e2b95e15ce8abf732db51ec771f1ca2edccf22a46c729ac564726020808301919091527fb7a0bfa1b79f2443f4d73ebb9259cddbcd510b18be6fc4da7d1aa7b1786e73e6828401527fc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc6606083015260808201939093523060a08201527f251543af6a222378665a76fe38dbceae4871a070b7fdaf5c6c30cf758dc33cc060c0808301919091528251808303909101815260e090910190915280519101206004556103a6565b82805482825590600052602060002090810192821561036f579160200282015b8281111561036f57825182546001600160a01b0319166001600160a01b0390911617825560209092019160019091019061033a565b5061037b92915061037f565b5090565b6103a391905b8082111561037b5780546001600160a01b0319168155600101610385565b90565b61097e806103b56000396000f3fe6080604052600436106100555760003560e01c80630d8e6e2c1461008d57806342cde4e814610117578063a0ab96531461013e578063aa5df9e21461039b578063affed0e0146103e1578063ca7541ee146103f6575b60408051348152905133917fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c919081900360200190a2005b34801561009957600080fd5b506100a2610421565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100dc5781810151838201526020016100c4565b50505050905090810190601f1680156101095780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561012357600080fd5b5061012c61043f565b60408051918252519081900360200190f35b34801561014a57600080fd5b50610399600480360361010081101561016257600080fd5b810190602081018135600160201b81111561017c57600080fd5b82018360208201111561018e57600080fd5b803590602001918460208302840111600160201b831117156101af57600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b8111156101fe57600080fd5b82018360208201111561021057600080fd5b803590602001918460208302840111600160201b8311171561023157600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b81111561028057600080fd5b82018360208201111561029257600080fd5b803590602001918460208302840111600160201b831117156102b357600080fd5b91908080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525092956001600160a01b0385351695602086013595919450925060608101915060400135600160201b81111561031757600080fd5b82018360208201111561032957600080fd5b803590602001918460018302840111600160201b8311171561034a57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550506001600160a01b038335169350505060200135610445565b005b3480156103a757600080fd5b506103c5600480360360208110156103be57600080fd5b5035610916565b604080516001600160a01b039092168252519081900360200190f35b3480156103ed57600080fd5b5061012c61093d565b34801561040257600080fd5b5061040b610943565b6040805160ff9092168252519081900360200190f35b604080518082019091526004815263322e333360e01b602082015290565b60015481565b60015487511461049c576040805162461bcd60e51b815260206004820152601c60248201527f52206c656e206e6f7420657175616c20746f207468726573686f6c6400000000604482015290519081900360640190fd5b855187511480156104ae575087518751145b6104ff576040805162461bcd60e51b815260206004820152601960248201527f6c656e677468206f6620722f732f76206e6f74206d6174636800000000000000604482015290519081900360640190fd5b6001600160a01b03821633148061051d57506001600160a01b038216155b61055f576040805162461bcd60e51b815260206004820152600e60248201526d3bb937b7339032bc32b1baba37b960911b604482015290519081900360640190fd5b825160208085019190912060008054604080517f3ee892349ae4bbe61dce18f95115b5dc02daf49204cc602458cd4c1f540d56d7818701526001600160a01b038b81168284015260608083018c9052608083019690965260a082019390935291871660c083015260e0808301879052815180840390910181526101008301825280519086012060045461190160f01b6101208501526101228401526101428084018290528251808503909101815261016284018084528151918801919091206001548083529788029094016101820190925294919391801561064b578160200160208202803883390190505b50905060005b6001548110156107ac5760006001858f848151811061066c57fe5b60200260200101518f858151811061068057fe5b60200260200101518f868151811061069457fe5b602002602001015160405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa1580156106f3573d6000803e3d6000fd5b505050602060405103519050836001600160a01b0316816001600160a01b031611801561073857506001600160a01b03811660009081526002602052604090205460ff165b61077d576040805162461bcd60e51b815260206004820152601160248201527015995c9a599e481cda59c819985a5b1959607a1b604482015290519081900360640190fd5b8083838151811061078a57fe5b6001600160a01b03909216602092830291909101909101529250600101610651565b5060008054600101815587518190819060208b018c8e8bf1905080610806576040805162461bcd60e51b815260206004820152600b60248201526a6e6f745f7375636365737360a81b604482015290519081900360640190fd5b7f07f4110a9f6788eae6a0b088d9aca06ec3cd9e2c6eae12a1d393d6d041d18c30828b8b8b6040518080602001856001600160a01b03166001600160a01b0316815260200184815260200180602001838103835287818151815260200191508051906020019060200280838360005b8381101561088d578181015183820152602001610875565b50505050905001838103825284818151815260200191508051906020019080838360005b838110156108c95781810151838201526020016108b1565b50505050905090810190601f1680156108f65780820380516001836020036101000a031916815260200191505b50965050505050505060405180910390a150505050505050505050505050565b6003818154811061092357fe5b6000918252602090912001546001600160a01b0316905081565b60005481565b6003549056fea265627a7a72305820fbcf75b0fe8efcca88f733aa6f967b0d614d263723e263240fa03f39f29b0bca64736f6c63430005090032`

// DeploySimpleMultiSig deploys a new Ethereum contract, binding an instance of SimpleMultiSig to it.
func DeploySimpleMultiSig(auth *bind.TransactOpts, backend bind.ContractBackend, threshold_ *big.Int, owners_ []common.Address, chainId *big.Int) (common.Address, *types.Transaction, *SimpleMultiSig, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleMultiSigABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SimpleMultiSigBin), backend, threshold_, owners_, chainId)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimpleMultiSig{SimpleMultiSigCaller: SimpleMultiSigCaller{contract: contract}, SimpleMultiSigTransactor: SimpleMultiSigTransactor{contract: contract}, SimpleMultiSigFilterer: SimpleMultiSigFilterer{contract: contract}}, nil
}

// SimpleMultiSig is an auto generated Go binding around an Ethereum contract.
type SimpleMultiSig struct {
	SimpleMultiSigCaller     // Read-only binding to the contract
	SimpleMultiSigTransactor // Write-only binding to the contract
	SimpleMultiSigFilterer   // Log filterer for contract events
}

// SimpleMultiSigCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimpleMultiSigCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleMultiSigTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimpleMultiSigTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleMultiSigFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimpleMultiSigFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleMultiSigSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimpleMultiSigSession struct {
	Contract     *SimpleMultiSig   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SimpleMultiSigCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimpleMultiSigCallerSession struct {
	Contract *SimpleMultiSigCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// SimpleMultiSigTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimpleMultiSigTransactorSession struct {
	Contract     *SimpleMultiSigTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// SimpleMultiSigRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimpleMultiSigRaw struct {
	Contract *SimpleMultiSig // Generic contract binding to access the raw methods on
}

// SimpleMultiSigCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimpleMultiSigCallerRaw struct {
	Contract *SimpleMultiSigCaller // Generic read-only contract binding to access the raw methods on
}

// SimpleMultiSigTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimpleMultiSigTransactorRaw struct {
	Contract *SimpleMultiSigTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimpleMultiSig creates a new instance of SimpleMultiSig, bound to a specific deployed contract.
func NewSimpleMultiSig(address common.Address, backend bind.ContractBackend) (*SimpleMultiSig, error) {
	contract, err := bindSimpleMultiSig(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleMultiSig{SimpleMultiSigCaller: SimpleMultiSigCaller{contract: contract}, SimpleMultiSigTransactor: SimpleMultiSigTransactor{contract: contract}, SimpleMultiSigFilterer: SimpleMultiSigFilterer{contract: contract}}, nil
}

// NewSimpleMultiSigCaller creates a new read-only instance of SimpleMultiSig, bound to a specific deployed contract.
func NewSimpleMultiSigCaller(address common.Address, caller bind.ContractCaller) (*SimpleMultiSigCaller, error) {
	contract, err := bindSimpleMultiSig(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleMultiSigCaller{contract: contract}, nil
}

// NewSimpleMultiSigTransactor creates a new write-only instance of SimpleMultiSig, bound to a specific deployed contract.
func NewSimpleMultiSigTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleMultiSigTransactor, error) {
	contract, err := bindSimpleMultiSig(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleMultiSigTransactor{contract: contract}, nil
}

// NewSimpleMultiSigFilterer creates a new log filterer instance of SimpleMultiSig, bound to a specific deployed contract.
func NewSimpleMultiSigFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleMultiSigFilterer, error) {
	contract, err := bindSimpleMultiSig(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleMultiSigFilterer{contract: contract}, nil
}

// bindSimpleMultiSig binds a generic wrapper to an already deployed contract.
func bindSimpleMultiSig(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleMultiSigABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleMultiSig *SimpleMultiSigRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SimpleMultiSig.Contract.SimpleMultiSigCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleMultiSig *SimpleMultiSigRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleMultiSig.Contract.SimpleMultiSigTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleMultiSig *SimpleMultiSigRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleMultiSig.Contract.SimpleMultiSigTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleMultiSig *SimpleMultiSigCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SimpleMultiSig.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleMultiSig *SimpleMultiSigTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleMultiSig.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleMultiSig *SimpleMultiSigTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleMultiSig.Contract.contract.Transact(opts, method, params...)
}

// GetOwersLength is a free data retrieval call binding the contract method 0xca7541ee.
//
// Solidity: function getOwersLength() constant returns(uint8)
func (_SimpleMultiSig *SimpleMultiSigCaller) GetOwersLength(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _SimpleMultiSig.contract.Call(opts, out, "getOwersLength")
	return *ret0, err
}

// GetOwersLength is a free data retrieval call binding the contract method 0xca7541ee.
//
// Solidity: function getOwersLength() constant returns(uint8)
func (_SimpleMultiSig *SimpleMultiSigSession) GetOwersLength() (uint8, error) {
	return _SimpleMultiSig.Contract.GetOwersLength(&_SimpleMultiSig.CallOpts)
}

// GetOwersLength is a free data retrieval call binding the contract method 0xca7541ee.
//
// Solidity: function getOwersLength() constant returns(uint8)
func (_SimpleMultiSig *SimpleMultiSigCallerSession) GetOwersLength() (uint8, error) {
	return _SimpleMultiSig.Contract.GetOwersLength(&_SimpleMultiSig.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() constant returns(string)
func (_SimpleMultiSig *SimpleMultiSigCaller) GetVersion(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SimpleMultiSig.contract.Call(opts, out, "getVersion")
	return *ret0, err
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() constant returns(string)
func (_SimpleMultiSig *SimpleMultiSigSession) GetVersion() (string, error) {
	return _SimpleMultiSig.Contract.GetVersion(&_SimpleMultiSig.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() constant returns(string)
func (_SimpleMultiSig *SimpleMultiSigCallerSession) GetVersion() (string, error) {
	return _SimpleMultiSig.Contract.GetVersion(&_SimpleMultiSig.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() constant returns(uint256)
func (_SimpleMultiSig *SimpleMultiSigCaller) Nonce(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleMultiSig.contract.Call(opts, out, "nonce")
	return *ret0, err
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() constant returns(uint256)
func (_SimpleMultiSig *SimpleMultiSigSession) Nonce() (*big.Int, error) {
	return _SimpleMultiSig.Contract.Nonce(&_SimpleMultiSig.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() constant returns(uint256)
func (_SimpleMultiSig *SimpleMultiSigCallerSession) Nonce() (*big.Int, error) {
	return _SimpleMultiSig.Contract.Nonce(&_SimpleMultiSig.CallOpts)
}

// OwnersArr is a free data retrieval call binding the contract method 0xaa5df9e2.
//
// Solidity: function ownersArr(uint256 ) constant returns(address)
func (_SimpleMultiSig *SimpleMultiSigCaller) OwnersArr(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SimpleMultiSig.contract.Call(opts, out, "ownersArr", arg0)
	return *ret0, err
}

// OwnersArr is a free data retrieval call binding the contract method 0xaa5df9e2.
//
// Solidity: function ownersArr(uint256 ) constant returns(address)
func (_SimpleMultiSig *SimpleMultiSigSession) OwnersArr(arg0 *big.Int) (common.Address, error) {
	return _SimpleMultiSig.Contract.OwnersArr(&_SimpleMultiSig.CallOpts, arg0)
}

// OwnersArr is a free data retrieval call binding the contract method 0xaa5df9e2.
//
// Solidity: function ownersArr(uint256 ) constant returns(address)
func (_SimpleMultiSig *SimpleMultiSigCallerSession) OwnersArr(arg0 *big.Int) (common.Address, error) {
	return _SimpleMultiSig.Contract.OwnersArr(&_SimpleMultiSig.CallOpts, arg0)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() constant returns(uint256)
func (_SimpleMultiSig *SimpleMultiSigCaller) Threshold(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleMultiSig.contract.Call(opts, out, "threshold")
	return *ret0, err
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() constant returns(uint256)
func (_SimpleMultiSig *SimpleMultiSigSession) Threshold() (*big.Int, error) {
	return _SimpleMultiSig.Contract.Threshold(&_SimpleMultiSig.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() constant returns(uint256)
func (_SimpleMultiSig *SimpleMultiSigCallerSession) Threshold() (*big.Int, error) {
	return _SimpleMultiSig.Contract.Threshold(&_SimpleMultiSig.CallOpts)
}

// Execute is a paid mutator transaction binding the contract method 0xa0ab9653.
//
// Solidity: function execute(uint8[] sigV, bytes32[] sigR, bytes32[] sigS, address destination, uint256 value, bytes data, address executor, uint256 gasLimit) returns()
func (_SimpleMultiSig *SimpleMultiSigTransactor) Execute(opts *bind.TransactOpts, sigV []uint8, sigR [][32]byte, sigS [][32]byte, destination common.Address, value *big.Int, data []byte, executor common.Address, gasLimit *big.Int) (*types.Transaction, error) {
	return _SimpleMultiSig.contract.Transact(opts, "execute", sigV, sigR, sigS, destination, value, data, executor, gasLimit)
}

// Execute is a paid mutator transaction binding the contract method 0xa0ab9653.
//
// Solidity: function execute(uint8[] sigV, bytes32[] sigR, bytes32[] sigS, address destination, uint256 value, bytes data, address executor, uint256 gasLimit) returns()
func (_SimpleMultiSig *SimpleMultiSigSession) Execute(sigV []uint8, sigR [][32]byte, sigS [][32]byte, destination common.Address, value *big.Int, data []byte, executor common.Address, gasLimit *big.Int) (*types.Transaction, error) {
	return _SimpleMultiSig.Contract.Execute(&_SimpleMultiSig.TransactOpts, sigV, sigR, sigS, destination, value, data, executor, gasLimit)
}

// Execute is a paid mutator transaction binding the contract method 0xa0ab9653.
//
// Solidity: function execute(uint8[] sigV, bytes32[] sigR, bytes32[] sigS, address destination, uint256 value, bytes data, address executor, uint256 gasLimit) returns()
func (_SimpleMultiSig *SimpleMultiSigTransactorSession) Execute(sigV []uint8, sigR [][32]byte, sigS [][32]byte, destination common.Address, value *big.Int, data []byte, executor common.Address, gasLimit *big.Int) (*types.Transaction, error) {
	return _SimpleMultiSig.Contract.Execute(&_SimpleMultiSig.TransactOpts, sigV, sigR, sigS, destination, value, data, executor, gasLimit)
}

// SimpleMultiSigDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the SimpleMultiSig contract.
type SimpleMultiSigDepositIterator struct {
	Event *SimpleMultiSigDeposit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleMultiSigDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleMultiSigDeposit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SimpleMultiSigDeposit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleMultiSigDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleMultiSigDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleMultiSigDeposit represents a Deposit event raised by the SimpleMultiSig contract.
type SimpleMultiSigDeposit struct {
	From  common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed _from, uint256 _value)
func (_SimpleMultiSig *SimpleMultiSigFilterer) FilterDeposit(opts *bind.FilterOpts, _from []common.Address) (*SimpleMultiSigDepositIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _SimpleMultiSig.contract.FilterLogs(opts, "Deposit", _fromRule)
	if err != nil {
		return nil, err
	}
	return &SimpleMultiSigDepositIterator{contract: _SimpleMultiSig.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed _from, uint256 _value)
func (_SimpleMultiSig *SimpleMultiSigFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *SimpleMultiSigDeposit, _from []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}

	logs, sub, err := _SimpleMultiSig.contract.WatchLogs(opts, "Deposit", _fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleMultiSigDeposit)
				if err := _SimpleMultiSig.contract.UnpackLog(event, "Deposit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SimpleMultiSigExecuteIterator is returned from FilterExecute and is used to iterate over the raw logs and unpacked data for Execute events raised by the SimpleMultiSig contract.
type SimpleMultiSigExecuteIterator struct {
	Event *SimpleMultiSigExecute // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SimpleMultiSigExecuteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleMultiSigExecute)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SimpleMultiSigExecute)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SimpleMultiSigExecuteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleMultiSigExecuteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleMultiSigExecute represents a Execute event raised by the SimpleMultiSig contract.
type SimpleMultiSigExecute struct {
	ConfirmAddrs []common.Address
	Destination  common.Address
	Value        *big.Int
	Data         []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterExecute is a free log retrieval operation binding the contract event 0x07f4110a9f6788eae6a0b088d9aca06ec3cd9e2c6eae12a1d393d6d041d18c30.
//
// Solidity: event Execute(address[] _confirmAddrs, address _destination, uint256 _value, bytes data)
func (_SimpleMultiSig *SimpleMultiSigFilterer) FilterExecute(opts *bind.FilterOpts) (*SimpleMultiSigExecuteIterator, error) {

	logs, sub, err := _SimpleMultiSig.contract.FilterLogs(opts, "Execute")
	if err != nil {
		return nil, err
	}
	return &SimpleMultiSigExecuteIterator{contract: _SimpleMultiSig.contract, event: "Execute", logs: logs, sub: sub}, nil
}

// WatchExecute is a free log subscription operation binding the contract event 0x07f4110a9f6788eae6a0b088d9aca06ec3cd9e2c6eae12a1d393d6d041d18c30.
//
// Solidity: event Execute(address[] _confirmAddrs, address _destination, uint256 _value, bytes data)
func (_SimpleMultiSig *SimpleMultiSigFilterer) WatchExecute(opts *bind.WatchOpts, sink chan<- *SimpleMultiSigExecute) (event.Subscription, error) {

	logs, sub, err := _SimpleMultiSig.contract.WatchLogs(opts, "Execute")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleMultiSigExecute)
				if err := _SimpleMultiSig.contract.UnpackLog(event, "Execute", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
