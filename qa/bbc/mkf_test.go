package bbc

import (
	"fmt"
	"testing"

	"github.com/dabankio/bbrpc"
	"github.com/dabankio/devtools4chains"
	"github.com/dabankio/wallet-core/bip39"
	"github.com/dabankio/wallet-core/bip44"
	"github.com/dabankio/wallet-core/core/bbc"
	"github.com/stretchr/testify/require"
)

const mkfImage = "dabankio/mkfdev:1"

func TestDecodeTX(t *testing.T) {
	t.Skip("bad test")
	raw := "0200000051b64d5f00000000017a26dfcba4ed6c3e46a7b0f0ec6c8214a23c211613559b8768a843aa4cb64d5f0002030001e21d6d49931304681ac8ed683d8e90dc8eb6793a875d5361b0bb72bdf9601823000000000030750000000000000000"
	tx, err := bbc.DecodeTX(raw)
	fmt.Println(tx, err)
}

// 演示BBC sdk一般性用法
// 警告: 不要在生产环境中直接使用注释中的助记词
func TestExampleMKF(t *testing.T) {
	r := require.New(t)
	const pass = "123"
	// 临时启动一个测试节点

	nodeInfo := devtools4chains.MustRunDockerMKFDev(t, mkfImage, true, true)

	jsonRPC := nodeInfo.Client
	minerAddress := nodeInfo.MinerAddress

	var err error
	var key *bbc.KeyInfo
	t.Run("私钥、地址推导", func(t *testing.T) {
		entropy, err := bip39.NewEntropy(128) // <<=== sdk 生成熵, 128-256 32倍数
		require.NoError(t, err)
		err = bip39.SetWordListLang(bip39.LangChineseSimplified) // <<=== sdk 设定助记词语言,参考语言常量
		require.NoError(t, err)
		mnemonic, err := bip39.NewMnemonic(entropy) // <<=== sdk 生成助记词
		require.NoError(t, err)
		fmt.Println("mnemonic:", mnemonic)  //mnemonic: 旦 件 言 毫 树 名 当 氧 旨 弧 落 功
		seed := bip39.NewSeed(mnemonic, "") // <<=== sdk 获取种子，第二个参数相当于salt,生产后请始终保持一致

		d, err := bbc.NewSymbolCoin("MKF", seed, bip44.PathFormat, "")
		require.NoError(t, err)
		r.NoError(err)
		add, err := d.DeriveAddress()
		require.NoError(t, err)
		privk, err := d.DerivePrivateKey()
		require.NoError(t, err)
		pubk, err := d.DerivePublicKey()
		require.NoError(t, err)

		key = &bbc.KeyInfo{
			PrivateKey: privk,
			Address:    add,
			PublicKey:  pubk,
		}

		fmt.Println("key", key) //addr
	})

	registeredAssets := 12.34
	{ // 导入公钥
		_, err = jsonRPC.Importpubkey(key.PublicKey) // <<=== RPC 导入公钥
		r.NoError(err)
		r.NoError(bbrpc.Wait4balanceReach(minerAddress, 10, jsonRPC))
		jsonRPC.Unlockkey(nodeInfo.MinerOwnerPubk, nodeInfo.UnlockPass, nil)
		_, err = jsonRPC.Sendfrom(bbrpc.CmdSendfrom{
			From: minerAddress, To: key.Address, Amount: registeredAssets,
		})
		r.NoError(err)
		r.NoError(bbrpc.Wait4balanceReach(key.Address, registeredAssets, jsonRPC))
	}

	outAmount := 2.3

	t.Run("简单交易签名", func(t *testing.T) {
		//创建交易、签名、广播、检查余额
		rawTX, err := jsonRPC.Createtransaction(bbrpc.CmdCreatetransaction{ // <<=== RPC 创建交易
			From: key.Address, To: minerAddress, Amount: outAmount,
		})
		r.NoError(err)

		rawTX = replaceTXVersion(*rawTX)

		fmt.Println("rawTX:", *rawTX)
		deTx, err := bbc.DecodeSymbolTX("MKF", *rawTX) // <<=== sdk 反序列化交易
		r.NoError(err)
		fmt.Println("decoded tx", deTx) //decoded tx {"Version":1,"Typ":0,"Timestamp":1584952846,"LockUntil":0,"SizeIn":1,"Prefix":2,"Amount":1340000,"TxFee":100,"SizeOut":0,"SizeSign":0,"HashAnchor":"00000000c335f935650a427bf548242eac4e4a444e25691b47351e7945f4a8d4","Address":"10g06z2bmwb71n9xg9zsv4vzay86ab7avt6n97hm6ra2z3rsbrtc2ncer","Sign":""}

		signedTX, err := bbc.SymbolSignWithPrivateKey(bbc.SymbolMKF, *rawTX, "", key.PrivateKey) // <<=== sdk 使用私钥对交易进行签名
		r.NoError(err)

		sendTxid, err := jsonRPC.Sendtransaction(signedTX) // <<=== RPC 发送交易
		r.NoError(err)

		sdkTxid, err := bbc.CalcTxid("MKF", signedTX)
		r.NoError(err)
		rpcDe, err := jsonRPC.Decodetransaction(signedTX)
		r.NoError(err)
		fmt.Printf("txidS\n sdk: %s\n snd: %s \n rpc: %s", sdkTxid, *sendTxid, rpcDe.Txid)
		r.Equal(sdkTxid, *sendTxid)

		r.NoError(bbrpc.Wait4nBlocks(1, jsonRPC))

		bal, err := jsonRPC.Getbalance(nil, &key.Address) // <<=== RPC 查询余额
		r.NoError(err)
		r.Len(bal, 1)
		r.True(bal[0].Avail < registeredAssets-outAmount)
		fmt.Println("balance after send", bal[0]) //balance after send {1dmyvkbkbk5zaqvx46zqpy2vzywjz02sv5kdd0gq2c56mwb48925hfhpd 0.9899 0 0}
	})

	t.Run("直接使用私钥", func(t *testing.T) {
		key, err = bbc.ParsePrivateKey(key.PrivateKey) // <<=== sdk 解析私钥为公钥/地址
		require.NoError(t, err)
	})

	t.Run("多重签名示例", func(t *testing.T) { //多签示例,警告：不要使用示例中的私钥
		r := require.New(t)
		member0 := bbrpc.AddrKeypair{Keypair: bbrpc.Keypair{
			Privkey: "195cd69eff4580ad2430f92d2c86865c596e72edb33f40df5d41c97883241c7c",
			Pubkey:  "a7386f6cbe769fda91462637393970850ae7528d2cee5214c26cc4b27c014a65",
		}, Address: "1cn502z5jrhpc452jxrp8tmq71a2q0e9s6wk4d4etkxvbwv3f72ksbkdn"}
		member1 := bbrpc.AddrKeypair{Keypair: bbrpc.Keypair{
			Privkey: "3de774bfb200a46f6d969f5e080572859bc5d7b297fdb34471f55be3326b2153",
			Pubkey:  "1fb8c0c79a506fd8fcca12065331110ae4aedceb2eac38f75379174c6a5b1bff",
		}, Address: "1zwdnptjc2xwn7xsrngqeqq5ewg512cak0r9cnz6rdx89nhy0q0fstv2y"}
		member2 := bbrpc.AddrKeypair{Keypair: bbrpc.Keypair{
			Privkey: "8c49b0f3788e07025303ef763e55d14781c09d43cb749628d26280f8f6912336",
			Pubkey:  "5910534ab7629ccb73659df42afc3c382597223a9caa4040a687dbebbe1aa88a",
		}, Address: "1ham1nfqbve3tcg20nae3m8mq4mw3sz1ayjepawybkhhbejjk21cvjnx3"}

		//创建多签地址
		multisigAddress, err := jsonRPC.Addnewtemplate(bbrpc.AddnewtemplateParamMultisig{
			Required: 2,
			Pubkeys: []string{
				member0.Pubkey,
				member1.Pubkey,
				member2.Pubkey,
			},
		})
		r.NoError(err)

		//为多签地址准备资金
		_, err = jsonRPC.Sendfrom(bbrpc.CmdSendfrom{
			From: minerAddress, To: *multisigAddress, Amount: 100,
		})
		r.NoError(err)
		r.NoError(bbrpc.Wait4balanceReach(*multisigAddress, 100, jsonRPC)) //等待多签地址资金到账

		//从多签地址转出
		// 创建交易
		rawTX, err := jsonRPC.Createtransaction(bbrpc.CmdCreatetransaction{
			From:   *multisigAddress,
			To:     member0.Address,
			Amount: 23,
		})
		r.NoError(err)

		var templateData string //多签模版地址数据
		{
			vret, err := jsonRPC.Validateaddress(*multisigAddress)
			r.NoError(err)
			templateData = vret.Addressdata.Templatedata.Hex
		}

		//member0签名
		member0Sign, err := bbc.SymbolSignWithPrivateKey(bbc.SymbolMKF, *rawTX, templateData, member0.Privkey)
		r.NoError(err)

		//member1签名
		member1Sign, err := bbc.SymbolSignWithPrivateKey(bbc.SymbolMKF, member0Sign, templateData, member1.Privkey)
		r.NoError(err)

		txid, err := jsonRPC.Sendtransaction(member1Sign)
		r.NoError(err)
		fmt.Println("多签广播txid", *txid)

		r.NoError(bbrpc.Wait4nBlocks(1, jsonRPC))
		bal, err := jsonRPC.Getbalance(nil, nil)
		r.NoError(err)
		fmt.Printf("%#v\n", bal)
		// r.NoError(bbrpc.Wait4balanceReach(member0.Address, 20, jsonRPC)) //等待member0到账
	})
}
