package main

import (
	"fmt"
	client2 "github.com/koinotice/vite/client"
	"github.com/koinotice/vite/common/types"
	"github.com/koinotice/vite/ledger"

	//"github.com/koinotice/vite/ledger"

	//"github.com/koinotice/vite/common/types"

)

 var RawUrl = "http://127.0.0.1:48132"

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
func balance(){
	client, err := client2.NewRpcClient(RawUrl)
	if err != nil {
		fmt.Print(err)
		return
	}

	addr, err := types.HexToAddress("vite_512a50000f53c2aeb913f206cf9ee107a850e9f3f72caed6d4")
	if err != nil {
		fmt.Print(addr)
		return
	}
	//bs,e:=client.Balance(addr)
	//client2.RpcClient().
	b := TokenBalance{}
	query:=client2.BalanceQuery{
		Addr:    addr,
		TokenId: ledger.ViteTokenId,
	}
	err1:=client.GetClient().Call(&b, "ledger_getBalanceByAccAddrToken", query.Addr, query.TokenId)
	if err1 != nil {

	}


	fmt.Print(b)
} 
 
func main() {

	balance()
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
