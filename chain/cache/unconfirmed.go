package chain_cache

import (
	"github.com/koinotice/vite/common/types"
	"github.com/koinotice/vite/ledger"
)

func (cache *Cache) GetUnconfirmedBlocks() []*ledger.AccountBlock {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	return cache.unconfirmedPool.GetBlocks()
}

func (cache *Cache) GetUnconfirmedBlocksByAddress(address *types.Address) []*ledger.AccountBlock {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	return cache.unconfirmedPool.GetBlocksByAddress(address)
}
