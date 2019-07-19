package vm_db

import (
	"github.com/koinotice/vite/common/db"
	"github.com/koinotice/vite/interfaces"
)

// Cannot be concurrent with write
func (vdb *vmDb) NewStorageIterator(prefix []byte) (interfaces.StorageIterator, error) {
	iter, err := vdb.chain.GetStorageIterator(*vdb.address, prefix)
	if err != nil {
		return nil, err
	}

	unsavedIter := vdb.unsaved().NewStorageIterator(prefix)

	return db.NewMergedIterator([]interfaces.StorageIterator{
		unsavedIter,
		iter,
	}, vdb.unsaved().IsDelete), nil
}
