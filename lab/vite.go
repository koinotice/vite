package main

import (
	"fmt"
	client2 "github.com/koinotice/vite/client"
	"github.com/koinotice/vite/common/types"
	//"github.com/koinotice/vite/crypto/ed25519"
	"github.com/koinotice/vite/ledger"
	"github.com/koinotice/vite/wallet"
	"math/big"

	//"github.com/koinotice/vite/ledger"

	//"github.com/koinotice/vite/common/types"

)

var RawUrl = "wss://premainnet.vitewallet.com/test/ws"

// var RawUrl = "http://118.25.182.202:48132"

type TokenBalance struct {
	TotalAmount string `json:"totalAmount"`
	TInfo       `json:"tokenInfo"`
}
type TInfo struct {
	TokenId     types.TokenTypeId `json:"tokenId"`
	TokenName   string            `json:"tokenName"`
	TokenSymbol string            `json:"tokenSymbol"`
}

func balance() {
	client, err := client2.NewRpcClient(RawUrl)
	if err != nil {
		fmt.Print(err)
		return
	}

	addr, err := types.HexToAddress("vite_d3896e6f21a8c0c90ce66584ecd8c3e6fa75f557befd116e6d")
	if err != nil {
		fmt.Print(addr)
		return
	}



	bs, err := client.GetOnroad(client2.OnroadQuery{
		Address: addr,
		Index:   1,
		Cnt:     10,
	})
	if err != nil {
		fmt.Print(err)
		return
	}
	if len(bs) > 0 {
		for _, v := range bs {
			fmt.Print(v)
		}
	}

	//bs,e:=client.Balance(addr)
	//client2.RpcClient().
	b := TokenBalance{}
	query := client2.BalanceQuery{
		Addr:    addr,
		TokenId: ledger.ViteTokenId,
	}
	err1 := client.GetClient().Call(&b, "ledger_getBalanceByAccAddrToken", query.Addr, query.TokenId)
	if err1 != nil {

	}

	fmt.Printf("%s \n",b)
}

func mkwallet() {

	//manager1 := wallet.New(&wallet.Config{
	//	DataDir: "./wallet/",
	//})
	//mnemonic, em, err := manager1.NewMnemonicAndEntropyStore("123456")
	//if err != nil {
	//	fmt.Print(err)
	//}
	//fmt.Println(mnemonic, ",", em.GetEntropyStoreFile())

	manager := wallet.New(&wallet.Config{
		DataDir: fmt.Sprintf("/vedex/vite/lab/wallet"),
	})
	manager.Start()
	files := manager.ListAllEntropyFiles()
	for _, v := range files {
		storeManager, err := manager.GetEntropyStoreManager(v)
		if err != nil {
			fmt.Print(err)
		}
		storeManager.Unlock("123456")
		_, key, err := storeManager.DeriveForIndexPath(0)
		if err != nil {
			fmt.Print(err)
		}
		keys, err := key.PrivateKey()
		if err != nil {
			fmt.Print(err)
		}
		addr, err := key.Address()
		if err != nil {
			fmt.Print(err)
		}
		fmt.Printf("%s,0:%s,%s\n", addr, addr, keys.Hex())
	}
}

func send() {
	rpc, err := client2.NewRpcClient(RawUrl)
	if err != nil {
		fmt.Print(err)
		return
	}

	client, err := client2.NewClient(rpc)

	if err != nil {
		fmt.Print(err)
		return
	}
	to, err := types.HexToAddress("vite_f4c39830a39c87dcb63867cd7d6ab6cb9ef82824ef620d287f")
	self, err := types.HexToAddress("vite_159df2046e8ab0348fd3d25070bbeaec58abd2c03bb93dce0a")
	if err != nil {
		fmt.Print(err)
		return
	}

	//manager := wallet.New(&wallet.Config{
	//	DataDir: "./wallet",
	//})
	//mnemonic, em, err := manager.NewMnemonicAndEntropyStore("123456")
	//if err != nil {
	//	fmt.Print(err)
	//}
	//fmt.Println(mnemonic)
	//fmt.Println(em.GetPrimaryAddr())
	//fmt.Println(em.GetEntropyStoreFile())

	manager := wallet.New(&wallet.Config{
		DataDir: "/vedex/vite/lab/wallet",
	})

	manager.Start()

	 mnemonic, em, err := manager.NewMnemonicAndEntropyStore("123456")

	fmt.Println(mnemonic)
	 fmt.Println(em.GetPrimaryAddr())
	 fmt.Println(em.GetEntropyStoreFile())

	storeManager, err := manager.GetEntropyStoreManager("vite_159df2046e8ab0348fd3d25070bbeaec58abd2c03bb93dce0a")
	if err != nil {
		fmt.Print(err)
	}
	storeManager.Unlock("123456")

	_, key, err := storeManager.DeriveForIndexPath(0)
	if err != nil {
		panic(err)
	}
	keys, err := key.PrivateKey()
	if err != nil {
		panic(err)
	}
	addr, err := key.Address()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s,0:%s,%s %d,\n", addr, keys, keys.Hex(), len(keys))

	//
	//for i := 0; i < 25; i++ {
	//	_, key, _ := storeManager.DeriveForIndexPath(uint32(i))
	//	address, _ := key.Address()
	//	fmt.Println(strconv.Itoa(i) + ":" + address.String())
	//}
	//mnemonic, em, err := manager.NewMnemonicAndEntropyStore("123456")
	//if err != nil {
	//	fmt.Print(err)
	//}
	//fmt.Println(mnemonic, ",", em.GetEntropyStoreFile())
	//despair candy brain nerve cart stereo undo arrange spirit wing valve leopard lawsuit fetch label beyond zone orange zebra photo canal fruit guide panic
	//vite_51fe5388251f0474502c7400a30fb7ea9adaf87275b5c1dad0
	//wallet/vite_51fe5388251f0474502c7400a30fb7ea9adaf87275b5c1dad0
	//

	DU,err:=types.HexToTokenTypeId("tti_d7d6d5fe81d5f8c69d9c6e17")
	if err!=nil{
		fmt.Print(err)
	}

	//privateKey:=[]byte("7aaeaddd27a2ccb468976c5c42702c1fa990c28de987b3010f378b859e2de705fa65ed3b349471f58388308fda141b4638baf1ed2509f78f647a0c5c1378a1fc")
	hashH, err := client.SubmitRequestTx(client2.RequestTxParams{
		ToAddr:   to,
		SelfAddr: self,
		Amount:   big.NewInt(1999999999999999999),
		TokenId:  DU,
		Data:     []byte("hello pow"),
	}, nil, func(addr types.Address, data []byte) (signedData, pubkey []byte, err error) {

		//ed25519.Sign(priv, message), priv.PubByte()
		//return ed25519.Sign(keys,data),keys.PubByte(),err
		return key.SignData(data)
	})
	fmt.Print(hashH)
	if err != nil {
		fmt.Print(err)
		fmt.Print("err")
		return
	}
}

func main() {
	//mkwallet()
	//go balance()
	go send()
	select {}
	//client, err := client2.NewRpcClient(RawUrl)
	//if err != nil {
	//	 fmt.Print(err)
	//	return
	//}
	////fmt.Print(client)
	//hash, err := types.HexToHash("bfff83c40823c60ff8b28430f988334e60f49a9adacfc4b94b2fce224aa97d14")
	//if err != nil {
	//	fmt.Print(err)
	//	return
	//}
	//block, e := client.GetAccBlock(hash)
	//if e != nil {
	//	fmt.Print(e)
	//	return
	//}
	//fmt.Print(block)
	//fmt.Print(block.TokenId)
	//

	//client, err := vite.NewRpcClient(RawUrl)
	//if err != nil {
	//	fmt.Print(err)
	//	return
	//}
	//
	//accountAddress, _ := types.HexToAddress("vite_00000000000000000000000000000000000000056ad6d26692")
	//toAddress, _ := types.HexToAddress("vite_098dfae02679a4ca05a4c8bf5dd00a8757f0c622bfccce7d68")
	//tokenId, _ := types.HexToTokenTypeId("tti_3cd880a76b7524fc2694d607")
	//snapshotHash, _ := types.HexToHash("68d458d52a13d5594c069a365345d2067ccbceb63680ec384697dda88de2ada8")
	//publicKey, _ := hex.DecodeString("4sYVHCR0fnpUZy3Acj8Wy0JOU81vH/khAW1KLYb19Hk=")
	//
	//amount := big.NewInt(1000000000).String()
	//block := vite.RawBlock{
	//	BlockType: 3,
	//	//PrevHash:       prevHash,
	//	AccountAddress: accountAddress,
	//	PublicKey:      publicKey,
	//	ToAddress:      toAddress,
	//	TokenId:        tokenId,
	//	SnapshotHash:   snapshotHash,
	//	Height:         "6",
	//	Amount:         &amount,
	//	Timestamp:      time.Now().Unix(),
	//}
	//err = client.SubmitRaw(block)
	//if err != nil {
	//	fmt.Print(err)
	//	return
	//}
}
