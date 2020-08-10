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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SimpleMultiSigABI is the input ABI used to generate the binding from.
const SimpleMultiSigABI = "[{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"nonceNum_\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"threshold_\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"owners_\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"_confirmAddrs\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_destination\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"Execute\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"bucketIdx\",\"type\":\"uint16\"},{\"internalType\":\"uint256\",\"name\":\"expireTime\",\"type\":\"uint256\"},{\"internalType\":\"uint8[]\",\"name\":\"sigV\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"sigR\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"sigS\",\"type\":\"bytes32[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBucketLength\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwersLength\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"nonceBucket\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ownersArr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"threshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// SimpleMultiSigFuncSigs maps the 4-byte function signature to its string representation.
var SimpleMultiSigFuncSigs = map[string]string{
	"163947e5": "execute(uint16,uint256,uint8[],bytes32[],bytes32[],address,uint256,bytes,address,uint256)",
	"c683a2db": "getBucketLength()",
	"ca7541ee": "getOwersLength()",
	"0d8e6e2c": "getVersion()",
	"c3142f4a": "nonceBucket(uint256)",
	"aa5df9e2": "ownersArr(uint256)",
	"42cde4e8": "threshold()",
}

// SimpleMultiSigBin is the compiled bytecode used for deploying new contracts.
var SimpleMultiSigBin = "0x60806040523480156200001157600080fd5b5060405162000e7638038062000e76833981810160405260808110156200003757600080fd5b815160208301516040808501805191519395929483019291846401000000008211156200006357600080fd5b9083019060208201858111156200007957600080fd5b82518660208202830111640100000000821117156200009757600080fd5b82525081516020918201928201910280838360005b83811015620000c6578181015183820152602001620000ac565b505050509190910160405250602001518351909250600a1080159150620000ee575081518311155b8015620000fb5750600083115b6200013d576040805162461bcd60e51b815260206004820152600d60248201526c18985917dd1a1c995cda1bdb19609a1b604482015290519081900360640190fd5b60008461ffff161180156200015757506101008461ffff16105b62000199576040805162461bcd60e51b815260206004820152600d60248201526c6261645f6e6f6e63655f6e756d60981b604482015290519081900360640190fd5b6000805b83518110156200029757816001600160a01b0316848281518110620001be57fe5b60200260200101516001600160a01b03161162000222576040805162461bcd60e51b815260206004820152601760248201527f72657065617465645f6f776e65722f756e736f72746564000000000000000000604482015290519081900360640190fd5b6001600260008684815181106200023557fe5b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060006101000a81548160ff0219169083151502179055508381815181106200028157fe5b602090810291909101015191506001016200019d565b508251620002ad906003906020860190620003d2565b50836001819055508461ffff16604051908082528060200260200182016040528015620002e4578160200160208202803883390190505b508051620002fb916000916020909101906200043c565b5050604080517fd87cd6ef79d4e2b95e15ce8abf732db51ec771f1ca2edccf22a46c729ac564726020808301919091527fb7a0bfa1b79f2443f4d73ebb9259cddbcd510b18be6fc4da7d1aa7b1786e73e6828401527fc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc6606083015260808201939093523060a08201527f251543af6a222378665a76fe38dbceae4871a070b7fdaf5c6c30cf758dc33cc060c0808301919091528251808303909101815260e0909101909152805191012060045550620004cf915050565b8280548282559060005260206000209081019282156200042a579160200282015b828111156200042a57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190620003f3565b506200043892915062000488565b5090565b8280548282559060005260206000209081019282156200047a579160200282015b828111156200047a5782518255916020019190600101906200045d565b5062000438929150620004b2565b620004af91905b80821115620004385780546001600160a01b03191681556001016200048f565b90565b620004af91905b80821115620004385760008155600101620004b9565b61099780620004df6000396000f3fe6080604052600436106100705760003560e01c8063aa5df9e21161004e578063aa5df9e2146103c6578063c3142f4a1461040c578063c683a2db14610436578063ca7541ee1461046257610070565b80630d8e6e2c146100a8578063163947e51461013257806342cde4e81461039f575b60408051348152905133917fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c919081900360200190a2005b3480156100b457600080fd5b506100bd610477565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100f75781810151838201526020016100df565b50505050905090810190601f1680156101245780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561013e57600080fd5b5061039d600480360361014081101561015657600080fd5b61ffff82351691602081013591810190606081016040820135600160201b81111561018057600080fd5b82018360208201111561019257600080fd5b803590602001918460208302840111600160201b831117156101b357600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b81111561020257600080fd5b82018360208201111561021457600080fd5b803590602001918460208302840111600160201b8311171561023557600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295949360208101935035915050600160201b81111561028457600080fd5b82018360208201111561029657600080fd5b803590602001918460208302840111600160201b831117156102b757600080fd5b91908080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525092956001600160a01b0385351695602086013595919450925060608101915060400135600160201b81111561031b57600080fd5b82018360208201111561032d57600080fd5b803590602001918460018302840111600160201b8311171561034e57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550506001600160a01b038335169350505060200135610496565b005b3480156103ab57600080fd5b506103b461090a565b60408051918252519081900360200190f35b3480156103d257600080fd5b506103f0600480360360208110156103e957600080fd5b5035610910565b604080516001600160a01b039092168252519081900360200190f35b34801561041857600080fd5b506103b46004803603602081101561042f57600080fd5b5035610937565b34801561044257600080fd5b5061044b610955565b6040805161ffff9092168252519081900360200190f35b34801561046e57600080fd5b5061044b61095b565b604080518082019091526005815264312e302e3160d81b602082015290565b6001548751146104d9576040805162461bcd60e51b81526020600482015260096024820152683130b22fb92fb632b760b91b604482015290519081900360640190fd5b855187511480156104eb575087518751145b61052c576040805162461bcd60e51b815260206004820152600d60248201526c3130b22fb632b72fb917b997bb60991b604482015290519081900360640190fd5b6001600160a01b03821633148061054a57506001600160a01b038216155b61058a576040805162461bcd60e51b815260206004820152600c60248201526b3130b22fb2bc32b1baba37b960a11b604482015290519081900360640190fd5b8842106105cd576040805162461bcd60e51b815260206004820152600c60248201526b1d1a5b5957d95e1c1a5c995960a21b604482015290519081900360640190fd5b60005461ffff8b161061061d576040805162461bcd60e51b81526020600482015260136024820152726275636b65745f6f75745f6f665f72616e676560681b604482015290519081900360640190fd5b6000808b61ffff168154811061062f57fe5b9060005260206000200154905060007f3ee892349ae4bbe61dce18f95115b5dc02daf49204cc602458cd4c1f540d56d760001b8b8888888051906020012086898960405160200180898152602001888152602001876001600160a01b03166001600160a01b03168152602001868152602001858152602001848152602001836001600160a01b03166001600160a01b0316815260200182815260200198505050505050505050604051602081830303815290604052805190602001209050600060045482604051602001808061190160f01b81525060020183815260200182815260200192505050604051602081830303815290604052805190602001209050600080905060008090505b6001548110156108665760006001848f848151811061075557fe5b60200260200101518f858151811061076957fe5b60200260200101518f868151811061077d57fe5b602002602001015160405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa1580156107dc573d6000803e3d6000fd5b505050602060405103519050826001600160a01b0316816001600160a01b031611801561082157506001600160a01b03811660009081526002602052604090205460ff165b61085c576040805162461bcd60e51b81526020600482015260076024820152666261645f73696760c81b604482015290519081900360640190fd5b915060010161073a565b5060008e61ffff168154811061087857fe5b906000526020600020015460010160008f61ffff168154811061089757fe5b90600052602060002001819055506000809050600080895160208b018c8e8bf19050806108f9576040805162461bcd60e51b815260206004820152600b60248201526a18d85b1b17d9985a5b195960aa1b604482015290519081900360640190fd5b505050505050505050505050505050565b60015481565b6003818154811061091d57fe5b6000918252602090912001546001600160a01b0316905081565b6000818154811061094457fe5b600091825260209091200154905081565b60005490565b6003549056fea2646970667358221220f8f4a55eaec2221782282b7a9a8049127f2626bc1ccc51074ad08293d9e9099464736f6c63430006000033"

// DeploySimpleMultiSig deploys a new Ethereum contract, binding an instance of SimpleMultiSig to it.
func DeploySimpleMultiSig(auth *bind.TransactOpts, backend bind.ContractBackend, nonceNum_ uint16, threshold_ *big.Int, owners_ []common.Address, chainId *big.Int) (common.Address, *types.Transaction, *SimpleMultiSig, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleMultiSigABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SimpleMultiSigBin), backend, nonceNum_, threshold_, owners_, chainId)
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

// GetBucketLength is a free data retrieval call binding the contract method 0xc683a2db.
//
// Solidity: function getBucketLength() view returns(uint16)
func (_SimpleMultiSig *SimpleMultiSigCaller) GetBucketLength(opts *bind.CallOpts) (uint16, error) {
	var (
		ret0 = new(uint16)
	)
	out := ret0
	err := _SimpleMultiSig.contract.Call(opts, out, "getBucketLength")
	return *ret0, err
}

// GetBucketLength is a free data retrieval call binding the contract method 0xc683a2db.
//
// Solidity: function getBucketLength() view returns(uint16)
func (_SimpleMultiSig *SimpleMultiSigSession) GetBucketLength() (uint16, error) {
	return _SimpleMultiSig.Contract.GetBucketLength(&_SimpleMultiSig.CallOpts)
}

// GetBucketLength is a free data retrieval call binding the contract method 0xc683a2db.
//
// Solidity: function getBucketLength() view returns(uint16)
func (_SimpleMultiSig *SimpleMultiSigCallerSession) GetBucketLength() (uint16, error) {
	return _SimpleMultiSig.Contract.GetBucketLength(&_SimpleMultiSig.CallOpts)
}

// GetOwersLength is a free data retrieval call binding the contract method 0xca7541ee.
//
// Solidity: function getOwersLength() view returns(uint16)
func (_SimpleMultiSig *SimpleMultiSigCaller) GetOwersLength(opts *bind.CallOpts) (uint16, error) {
	var (
		ret0 = new(uint16)
	)
	out := ret0
	err := _SimpleMultiSig.contract.Call(opts, out, "getOwersLength")
	return *ret0, err
}

// GetOwersLength is a free data retrieval call binding the contract method 0xca7541ee.
//
// Solidity: function getOwersLength() view returns(uint16)
func (_SimpleMultiSig *SimpleMultiSigSession) GetOwersLength() (uint16, error) {
	return _SimpleMultiSig.Contract.GetOwersLength(&_SimpleMultiSig.CallOpts)
}

// GetOwersLength is a free data retrieval call binding the contract method 0xca7541ee.
//
// Solidity: function getOwersLength() view returns(uint16)
func (_SimpleMultiSig *SimpleMultiSigCallerSession) GetOwersLength() (uint16, error) {
	return _SimpleMultiSig.Contract.GetOwersLength(&_SimpleMultiSig.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(string version)
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
// Solidity: function getVersion() view returns(string version)
func (_SimpleMultiSig *SimpleMultiSigSession) GetVersion() (string, error) {
	return _SimpleMultiSig.Contract.GetVersion(&_SimpleMultiSig.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(string version)
func (_SimpleMultiSig *SimpleMultiSigCallerSession) GetVersion() (string, error) {
	return _SimpleMultiSig.Contract.GetVersion(&_SimpleMultiSig.CallOpts)
}

// NonceBucket is a free data retrieval call binding the contract method 0xc3142f4a.
//
// Solidity: function nonceBucket(uint256 ) view returns(uint256)
func (_SimpleMultiSig *SimpleMultiSigCaller) NonceBucket(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SimpleMultiSig.contract.Call(opts, out, "nonceBucket", arg0)
	return *ret0, err
}

// NonceBucket is a free data retrieval call binding the contract method 0xc3142f4a.
//
// Solidity: function nonceBucket(uint256 ) view returns(uint256)
func (_SimpleMultiSig *SimpleMultiSigSession) NonceBucket(arg0 *big.Int) (*big.Int, error) {
	return _SimpleMultiSig.Contract.NonceBucket(&_SimpleMultiSig.CallOpts, arg0)
}

// NonceBucket is a free data retrieval call binding the contract method 0xc3142f4a.
//
// Solidity: function nonceBucket(uint256 ) view returns(uint256)
func (_SimpleMultiSig *SimpleMultiSigCallerSession) NonceBucket(arg0 *big.Int) (*big.Int, error) {
	return _SimpleMultiSig.Contract.NonceBucket(&_SimpleMultiSig.CallOpts, arg0)
}

// OwnersArr is a free data retrieval call binding the contract method 0xaa5df9e2.
//
// Solidity: function ownersArr(uint256 ) view returns(address)
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
// Solidity: function ownersArr(uint256 ) view returns(address)
func (_SimpleMultiSig *SimpleMultiSigSession) OwnersArr(arg0 *big.Int) (common.Address, error) {
	return _SimpleMultiSig.Contract.OwnersArr(&_SimpleMultiSig.CallOpts, arg0)
}

// OwnersArr is a free data retrieval call binding the contract method 0xaa5df9e2.
//
// Solidity: function ownersArr(uint256 ) view returns(address)
func (_SimpleMultiSig *SimpleMultiSigCallerSession) OwnersArr(arg0 *big.Int) (common.Address, error) {
	return _SimpleMultiSig.Contract.OwnersArr(&_SimpleMultiSig.CallOpts, arg0)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint256)
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
// Solidity: function threshold() view returns(uint256)
func (_SimpleMultiSig *SimpleMultiSigSession) Threshold() (*big.Int, error) {
	return _SimpleMultiSig.Contract.Threshold(&_SimpleMultiSig.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint256)
func (_SimpleMultiSig *SimpleMultiSigCallerSession) Threshold() (*big.Int, error) {
	return _SimpleMultiSig.Contract.Threshold(&_SimpleMultiSig.CallOpts)
}

// Execute is a paid mutator transaction binding the contract method 0x163947e5.
//
// Solidity: function execute(uint16 bucketIdx, uint256 expireTime, uint8[] sigV, bytes32[] sigR, bytes32[] sigS, address destination, uint256 value, bytes data, address executor, uint256 gasLimit) returns()
func (_SimpleMultiSig *SimpleMultiSigTransactor) Execute(opts *bind.TransactOpts, bucketIdx uint16, expireTime *big.Int, sigV []uint8, sigR [][32]byte, sigS [][32]byte, destination common.Address, value *big.Int, data []byte, executor common.Address, gasLimit *big.Int) (*types.Transaction, error) {
	return _SimpleMultiSig.contract.Transact(opts, "execute", bucketIdx, expireTime, sigV, sigR, sigS, destination, value, data, executor, gasLimit)
}

// Execute is a paid mutator transaction binding the contract method 0x163947e5.
//
// Solidity: function execute(uint16 bucketIdx, uint256 expireTime, uint8[] sigV, bytes32[] sigR, bytes32[] sigS, address destination, uint256 value, bytes data, address executor, uint256 gasLimit) returns()
func (_SimpleMultiSig *SimpleMultiSigSession) Execute(bucketIdx uint16, expireTime *big.Int, sigV []uint8, sigR [][32]byte, sigS [][32]byte, destination common.Address, value *big.Int, data []byte, executor common.Address, gasLimit *big.Int) (*types.Transaction, error) {
	return _SimpleMultiSig.Contract.Execute(&_SimpleMultiSig.TransactOpts, bucketIdx, expireTime, sigV, sigR, sigS, destination, value, data, executor, gasLimit)
}

// Execute is a paid mutator transaction binding the contract method 0x163947e5.
//
// Solidity: function execute(uint16 bucketIdx, uint256 expireTime, uint8[] sigV, bytes32[] sigR, bytes32[] sigS, address destination, uint256 value, bytes data, address executor, uint256 gasLimit) returns()
func (_SimpleMultiSig *SimpleMultiSigTransactorSession) Execute(bucketIdx uint16, expireTime *big.Int, sigV []uint8, sigR [][32]byte, sigS [][32]byte, destination common.Address, value *big.Int, data []byte, executor common.Address, gasLimit *big.Int) (*types.Transaction, error) {
	return _SimpleMultiSig.Contract.Execute(&_SimpleMultiSig.TransactOpts, bucketIdx, expireTime, sigV, sigR, sigS, destination, value, data, executor, gasLimit)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_SimpleMultiSig *SimpleMultiSigTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _SimpleMultiSig.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_SimpleMultiSig *SimpleMultiSigSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _SimpleMultiSig.Contract.Fallback(&_SimpleMultiSig.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_SimpleMultiSig *SimpleMultiSigTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _SimpleMultiSig.Contract.Fallback(&_SimpleMultiSig.TransactOpts, calldata)
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

// ParseDeposit is a log parse operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed _from, uint256 _value)
func (_SimpleMultiSig *SimpleMultiSigFilterer) ParseDeposit(log types.Log) (*SimpleMultiSigDeposit, error) {
	event := new(SimpleMultiSigDeposit)
	if err := _SimpleMultiSig.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	return event, nil
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

// ParseExecute is a log parse operation binding the contract event 0x07f4110a9f6788eae6a0b088d9aca06ec3cd9e2c6eae12a1d393d6d041d18c30.
//
// Solidity: event Execute(address[] _confirmAddrs, address _destination, uint256 _value, bytes data)
func (_SimpleMultiSig *SimpleMultiSigFilterer) ParseExecute(log types.Log) (*SimpleMultiSigExecute, error) {
	event := new(SimpleMultiSigExecute)
	if err := _SimpleMultiSig.contract.UnpackLog(event, "Execute", log); err != nil {
		return nil, err
	}
	return event, nil
}
