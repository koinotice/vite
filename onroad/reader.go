package onroad

import (
	"github.com/koinotice/vite/common/types"
	"github.com/koinotice/vite/net"
	"github.com/koinotice/vite/producer/producerevent"
	"github.com/koinotice/vite/vm_db"
)

type pool interface {
	AddDirectAccountBlock(address types.Address, vmAccountBlock *vm_db.VmAccountBlock) error
}

type producer interface {
	SetAccountEventFunc(func(producerevent.AccountEvent))
}

type netReader interface {
	SubscribeSyncStatus(fn func(net.SyncState)) (subId int)
	UnsubscribeSyncStatus(subId int)
	SyncState() net.SyncState
}
