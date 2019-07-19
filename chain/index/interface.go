package chain_index

import (
	"github.com/koinotice/vite/common/types"
	"github.com/koinotice/vite/ledger"
)

type Chain interface {
	IsContractAccount(address types.Address) (bool, error)

	LoadOnRoad(gid types.Gid) (map[types.Address]map[types.Address][]ledger.HashHeight, error)

	IterateContracts(iterateFunc func(addr types.Address, meta *ledger.ContractMeta, err error) bool)
}
