package eth

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/lomocoin/wallet-core/core/eth/contracts"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/lomocoin/wallet-core/core/eth/testtool"
)

const (
	ropstenRPCHost = "https://ropsten.infura.io/v3/69f29be376784f37b36a146ce7581efc"

	addrHex        = "0x6d2f5E9DDCa612ec835D943A2c117B870e3e9Adb"
	addrPrvk       = "74E1BB30C0C0BEC64986B49640ECB7EF453BEE0FAFEEA2B2EBB7C5F00F70E5F3"
	ropstenChainID = 4

	addr2Hex  = "0xBd573B656B36AB9c9781dede2f9cae7658cE3C08"
	addr2Prvk = "45c299c5415348a75678d40dfadaab55d5201e6201da6b2b8dbe3d2ca7773442"
)

//在ropsten上的测试跳过
func TestRopstenDeploySimpleMultisig(t *testing.T) {
	t.SkipNow() //在ropsten上的测试跳过

	mRequired := uint8(1)

	client, err := ethclient.Dial(ropstenRPCHost)
	testtool.FailOnErr(t, err, "dial failed")
	got, err := DeploySimpleMultiSigContract(*big.NewInt(ropstenChainID), client, addrPrvk, []string{addrHex, addr2Hex}, mRequired)
	if err != nil {
		t.Errorf("DeployMultiSigWalletContract() error = %v", err)
		t.FailNow()
	}
	fmt.Println("deployMultisigWalletContract got:", got)

	fmt.Println("============部署完成，请人工确认合约状态，拷贝地址进入下一阶段测试===========")
}

func TestTransferETH2multisigAddress(t *testing.T) {
	t.SkipNow() //在ropsten上的测试跳过

	contractAddress := "0x9344ffc32e48e3c0cadc9ac6444fdbbbc3c27d21"
	prvk, err := crypto.HexToECDSA(addrPrvk)
	testtool.FailOnErr(t, err, "prvk failed")
	addr := crypto.PubkeyToAddress(prvk.PublicKey)

	client, err := ethclient.Dial(ropstenRPCHost)
	testtool.FailOnErr(t, err, "dial failed")
	value := big.NewInt(testtool.E18 / 5)
	ctx := context.Background()
	nonce, err := client.NonceAt(ctx, addr, nil)
	testtool.FailOnErr(t, err, "Failed to get nonce")
	tx := types.NewTransaction(nonce, common.HexToAddress(contractAddress), value, uint64(6721975), big.NewInt(20000000000), nil)
	signer := types.MakeSigner(params.TestChainConfig, nil)
	signature, err := crypto.Sign(signer.Hash(tx).Bytes(), prvk)
	testtool.FailOnErr(t, err, "签名交易失败")
	tx, err = tx.WithSignature(signer, signature)
	testtool.FailOnErr(t, err, "为交易附加签名数据错误")
	err = client.SendTransaction(ctx, tx)
	testtool.FailOnErr(t, err, "充值到合约地址异常")
	fmt.Printf("==========转账完成，请人工确认余额，然后进入下一阶段========\nhttps://ropsten.etherscan.io/address/%s \n", contractAddress)
}

//在ropsten上的测试跳过
func TestRopstenSimpleMultisig(t *testing.T) {
	t.SkipNow() //在ropsten上的测试跳过

	contractAddress := "0xce87809b742789af9a5acac51d9d826cb9c721cf"
	contractAddressHex := "0xce87809b742789af9a5acac51d9d826cb9c721cf"
	// prvk, err := crypto.HexToECDSA(addrPrvk)
	// testtool.FailOnErr(t, err, "prvk failed")
	// addr := crypto.PubkeyToAddress(prvk.PublicKey)

	var (
		rpcClient *ethclient.Client
		err       error
	)
	rpcClient, err = ethclient.Dial(ropstenRPCHost)
	testtool.FailOnErr(t, err, "dial failed")

	// multisigContract, err := contracts.NewSimpleMultiSigCaller(contractAddress, rpcClient)
	// testtool.FailOnErr(t, err, fmt.Sprintf("构建多签合约调用时异常,检查合约地址和rpc server,%v", err))

	// addr:0xBd573B656B36AB9c9781dede2f9cae7658cE3C08
	// metamask://0x6d2f5E9DDCa612ec835D943A2c117B870e3e9Adb
	// metamask://private//74E1BB30C0C0BEC64986B49640ECB7EF453BEE0FAFEEA2B2EBB7C5F00F70E5F3

	{ // 交易测试
		transferValue := big.NewInt(testtool.E18 * 1.2)

		multisigContract, err := contracts.NewSimpleMultiSig(common.HexToAddress(contractAddress), rpcClient)
		testtool.FailOnErr(t, err, "构建多签合约调用时异常,检查合约地址和rpc server")
		nonce, err := multisigContract.Nonce(&bind.CallOpts{Pending: true})
		testtool.FailOnErr(t, err, "无法获取合约内部nonce")
		var (
			sigV                                 []uint8    //签名
			sigR, sigS                           [][32]byte //签名
			privkHex                             string
			multisigContractAddress, fromAddress string //多签合约地址，发起地址
			destination, executor                string //toAddress
			value, gasLimit                      *big.Int
			data                                 []byte
		)
		// 012由0发起，0和2签名, 把钱赚到1的地址上,executor 为0
		// 由metamask发起，自己签名，赚到另一个账户上
		privkHex = addrPrvk
		multisigContractAddress = contractAddressHex
		fromAddress = addrHex
		executor = addrHex
		destination = addr2Hex
		value = transferValue
		gasLimit = big.NewInt(239963)
		data = []byte("")

		v, r, s, err := SimpleMultiSigExecuteSign(ropstenChainID, addrPrvk, multisigContractAddress, destination, executor, uint64(nonce.Int64()), value, gasLimit, data)
		testtool.FailOnErr(t, err, "create sig failed")
		sigV = append(sigV, v)
		sigR = append(sigR, r)
		sigS = append(sigS, s)

		txid, err := ExecuteTX(&TxParams{
			backend:                 rpcClient,
			sigV:                    sigV,
			sigR:                    sigR,
			sigS:                    sigS,
			privkHex:                privkHex,
			multisigContractAddress: multisigContractAddress,
			fromAddress:             fromAddress,
			destination:             destination,
			executor:                executor,
			value:                   value,
			gasLimit:                gasLimit,
			data:                    data,
		})
		testtool.FailOnErr(t, err, "Execute Failed")
		fmt.Println("execute txid", txid)
	}

}
