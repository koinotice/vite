package verifier

import (
	"github.com/koinotice/vite/common/types"
	"github.com/koinotice/vite/consensus/core"
	"github.com/koinotice/vite/ledger"
	"github.com/koinotice/vite/vm_db"
)

type consensus interface {
	VerifyAccountProducer(block *ledger.AccountBlock) (bool, error)
	SBPReader() core.SBPStatReader
}

type onRoadPool interface {
	IsFrontOnRoadOfCaller(gid types.Gid, orAddr, caller types.Address, hash types.Hash) (bool, error)
}

type accountChain interface {
	vm_db.Chain

	IsReceived(sendBlockHash types.Hash) (bool, error)
	GetReceiveAbBySendAb(sendBlockHash types.Hash) (*ledger.AccountBlock, error)
	IsGenesisAccountBlock(block types.Hash) bool
	IsSeedConfirmedNTimes(blockHash types.Hash, n uint64) (bool, error)
}
