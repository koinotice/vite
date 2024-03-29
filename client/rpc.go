package client

import (
	rpc2 "github.com/koinotice/vite/client/rpc"
	"github.com/koinotice/vite/rpc"
)

type RpcClient interface {
	rpc2.LedgerApi
	rpc2.OnroadApi
	rpc2.TxApi
	rpc2.ContractApi
	rpc2.DexTradeApi

	GetClient() *rpc.Client
}

func NewRpcClient(rawurl string) (RpcClient, error) {
	c, err := rpc.Dial(rawurl)
	if err != nil {
		return nil, err
	}

	r := &rpcClient{
		LedgerApi:   rpc2.NewLedgerApi(c),
		OnroadApi:   rpc2.NewOnroadApi(c),
		TxApi:       rpc2.NewTxApi(c),
		ContractApi: rpc2.NewContractApi(c),
		DexTradeApi: rpc2.NewDexTradeApi(c),
		cc:          c,
	}
	return r, nil
}

type rpcClient struct {
	rpc2.LedgerApi
	rpc2.OnroadApi
	rpc2.TxApi
	rpc2.ContractApi
	rpc2.DexTradeApi

	cc *rpc.Client
}

func (c rpcClient) GetClient() *rpc.Client {
	return c.cc
}
