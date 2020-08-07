package internalized

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dabankio/wallet-core/core/eth/internalized/contracts"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// DeploySimpleMultiSigContract 创建部署一个多签合约，返回合约地址
// address: 签名的参与者， mRequired: m-n 至少多少人签名生效
// return: address_hex,txid
func DeploySimpleMultiSigContract(nonceBucketNum uint16, chainID *big.Int, backend bind.ContractBackend, privkHex string, hexAddress []string, mRequired uint8) (string, error) {
	//TODO 参数验证chainID枚举,privkHex格式,hexAddress格式和排序和长度
	var (
		err    error
		privk  *ecdsa.PrivateKey
		owners []common.Address
	)

	// 预处理: owners需要保持incr序，合约内的简单查重依赖于排序性
	for i := 0; i < len(hexAddress); i++ {
		hexAddress[i] = strings.ToLower(hexAddress[i])
	}
	sort.Slice(hexAddress, func(i, j int) bool {
		return hexAddress[i] < hexAddress[j]
	})

	// init vars
	privk, err = crypto.HexToECDSA(privkHex)
	if err != nil {
		return "", fmt.Errorf("转换私钥时发生错误,%v", err)
	}
	for _, ha := range hexAddress {
		owners = append(owners, common.HexToAddress(ha))
	}

	// 部署多签合约
	auth := bind.NewKeyedTransactor(privk)
	//TODO set gas,gasPrice
	add, tx, _, err := contracts.DeploySimpleMultiSig(auth, backend, nonceBucketNum, big.NewInt(int64(mRequired)), owners, chainID)
	if err != nil {
		return "", fmt.Errorf("_部署多签合约失败, %v", err)
	}
	// 下面的代码在测试链/主链无效，因为不是即时完成的
	// ver, err := w.GetVersion(&bind.CallOpts{Pending: true})
	// if err != nil {
	// 	return "", fmt.Errorf("无法调用合约函数获取version, %v", err)
	// }
	// return ver + "," + add.Hex() + "," + tx.Hash().Hex(), nil
	// info := add.Hex() + "," + tx.Hash().Hex() + ", ver:" + ver
	// fmt.Println("deploy info", info)
	_ = tx
	return add.Hex(), nil

}

// TxParams 交易参数
type TxParams struct {
	Backend                              bind.ContractBackend
	BucketIdx                            uint16
	ExpireTime                           time.Time
	SigV                                 []uint8    //签名
	SigR, SigS                           [][32]byte //签名
	PrivkHex                             string
	MultisigContractAddress, FromAddress string //多签合约地址，发起地址
	Destination, Executor                string //toAddress
	Value, GasLimit                      *big.Int
	Data                                 []byte
}

// ExecuteTX .
func ExecuteTX(txp *TxParams) (string, error) {
	var (
		err              error
		privk            *ecdsa.PrivateKey
		multisigContract *contracts.SimpleMultiSig
	)

	{ // init vars
		privk, err = crypto.HexToECDSA(txp.PrivkHex)
		if err != nil {
			return "", err
		}

		multisigContract, err = contracts.NewSimpleMultiSig(common.HexToAddress(txp.MultisigContractAddress), txp.Backend)
		if err != nil {
			return "", fmt.Errorf("构建多签合约调用时异常,检查合约地址和rpc server,%v", err)
		}
	}

	//TODO 参数校验
	{ // 调用合约方法
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		nonce, err := txp.Backend.PendingNonceAt(ctx, common.HexToAddress(txp.FromAddress))
		if err != nil {
			return "", fmt.Errorf("获取多签地址nonce时发生错误, %v", err)
		}

		tx, err := multisigContract.Execute(&bind.TransactOpts{
			From:     common.HexToAddress(txp.FromAddress),
			Nonce:    big.NewInt(int64(nonce)),
			GasLimit: uint64(txp.GasLimit.Int64()),
			Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
				signature, err := crypto.Sign(signer.Hash(tx).Bytes(), privk)
				if err != nil {
					return nil, err
				}
				// return tx.WithSignature(signer, append(signature, byte(12))) 这里瞎写的，测试失败的情况，可以删掉
				return tx.WithSignature(signer, signature)
			},
		},
			txp.BucketIdx,
			big.NewInt(txp.ExpireTime.Unix()),
			txp.SigV,
			txp.SigR,
			txp.SigS,
			common.HexToAddress(txp.Destination),
			txp.Value,
			txp.Data,
			common.HexToAddress(txp.Executor),
			txp.GasLimit)
		if err != nil {
			return "", fmt.Errorf("调用合约交易方法时发生错误, %v", err)
		}
		return tx.Hash().Hex(), nil
	}

}

func ecrecoverAddress(hashBytes, sig []byte) string {
	rePub, err := crypto.SigToPub(hashBytes, sig)
	// rePub, err := crypto.Ecrecover([]byte(hash), sig)
	if err != nil {
		panic(fmt.Errorf("ecrecover err: %v", err))
	}
	reAddr := crypto.PubkeyToAddress(*rePub)
	// addrFromPriv := crypto.PubkeyToAddress(privk.PublicKey)
	fmt.Println("addrFromPrivKey vs recoverdAddr")
	// fmt.Println(addrFromPriv.Hex())
	fmt.Println(reAddr.Hex())
	return reAddr.Hex()
}

//合约内部固定常量, 参考SimpleMultisig.sol
const (
	_TxTypeHash           = "0x3ee892349ae4bbe61dce18f95115b5dc02daf49204cc602458cd4c1f540d56d7"
	_NameHash             = "0xb7a0bfa1b79f2443f4d73ebb9259cddbcd510b18be6fc4da7d1aa7b1786e73e6"
	_VersionHash          = "0xc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc6"
	_Eip712DomaintypeHash = "0xd87cd6ef79d4e2b95e15ce8abf732db51ec771f1ca2edccf22a46c729ac56472"
	_Salt                 = "0x251543af6a222378665a76fe38dbceae4871a070b7fdaf5c6c30cf758dc33cc0"
	_AllZero              = "000000000000000000000000000000000000000000000000000000000000000000" //做padding用
)

// SimpleMultiSigExecuteSign return v,r,s
func SimpleMultiSigExecuteSign(expireTime time.Time, chainID int64, signerPrivkHex string, multisigContractAddr, destinationAddr, executor string, nonce *big.Int, value, gasLimit *big.Int, data []byte) (uint8, [32]byte, [32]byte, error) {
	zeroAndErr := func(e error) (uint8, [32]byte, [32]byte, error) { return 0, [32]byte{}, [32]byte{}, e }
	leftPad64 := func(str string) string { return _AllZero[:64-len(str)] + str } // 将小于64位的字符串(hex编码的)填充至64位（64位转为byte即32位，对应32*8=256 bit）
	hexKeccak256Hash := func(byts []byte) common.Hash {
		decodedData, err := hex.DecodeString(string(byts))
		if err != nil { //should not happen
			panic(errors.Wrap(err, "hex decode err"))
		}
		return crypto.Keccak256Hash([]byte(decodedData))
	}

	domainSeparatorHashHex := hexKeccak256Hash([]byte(strings.Join([]string{
		_Eip712DomaintypeHash[2:],
		_NameHash[2:],
		_VersionHash[2:],
		leftPad64(strconv.FormatInt(chainID, 16)),
		leftPad64(multisigContractAddr[2:]),
		_Salt[2:],
	}, "")))
	executor = "0x" + strings.TrimLeft(executor, "0x")
	txInputHashHex := hexKeccak256Hash([]byte(strings.Join([]string{
		_TxTypeHash[2:],
		leftPad64(big.NewInt(expireTime.Unix()).Text(16)),
		leftPad64(destinationAddr[2:]),
		leftPad64(value.Text(16)),
		crypto.Keccak256Hash(data).Hex()[2:],
		leftPad64(nonce.Text(16)),
		leftPad64(executor[2:]),
		leftPad64(gasLimit.Text(16)),
	}, "")))
	hashBytes := crypto.Keccak256(bytes.Join([][]byte{
		{25, 01}, //hex.DecodeString("1901") 的值, 0x19 01 为固定常量，参考合约函数execute
		domainSeparatorHashHex[:],
		txInputHashHex[:],
	}, nil))

	privk, err := crypto.HexToECDSA(signerPrivkHex)
	if err != nil {
		return zeroAndErr(err)
	}
	sig, err := crypto.Sign(hashBytes, privk)
	if err != nil {
		return zeroAndErr(errors.Wrap(err, "crypto sign failed"))
	}
	// if ecrecoverAddress(hashBytes, sig) != crypto.PubkeyToAddress(privk.PublicKey).Hex() {//【调试用】做内部的ecrecover验证,可移除
	// 	panic("ecrecover err")
	// }
	r, s, v := sig[:32], sig[32:64], sig[64]+27
	toBytes32 := func(b []byte) [32]byte {
		if len(b) == 32 {
			var b32 [32]byte
			copy(b32[:], b)
			return b32
		}
		panic(fmt.Sprintf("not [32]byte, actual len: %d", len(b))) //should not happen
	}
	return v, toBytes32(r), toBytes32(s), nil
}

// GetContractInfo 获取多签合约地址内的地址和签名人数
func GetContractInfo(caller bind.ContractCaller, contractAddress string) (addrs []string, mRequired int64, err error) {
	var multisigContract *contracts.SimpleMultiSigCaller

	{ // init vars
		multisigContract, err = contracts.NewSimpleMultiSigCaller(common.HexToAddress(contractAddress), caller)
		if err != nil {
			return nil, 0, fmt.Errorf("构建多签合约调用时异常,检查合约地址和rpc server,%v", err)
		}
	}
	callOpts := &bind.CallOpts{}
	{ //获取地址
		length, err := multisigContract.GetOwersLength(callOpts)
		if err != nil {
			return nil, 0, fmt.Errorf("获取持有人数量时发生异常, %v", err)
		}

		for i := uint16(0); i < length; i++ {
			addr, err := multisigContract.OwnersArr(callOpts, big.NewInt(int64(i)))
			if err != nil {
				return addrs, 0, fmt.Errorf("获取地址[%d/%d]时发生错误, %v", i+1, length, err)
			}
			addrs = append(addrs, addr.Hex())
		}
	}
	{ // 获取mRequired
		mInt, err := multisigContract.Threshold(callOpts)
		if err != nil {
			return nil, 0, fmt.Errorf("获取签名数时发生错误，%v", err)
		}
		mRequired = mInt.Int64()
	}
	return
}
