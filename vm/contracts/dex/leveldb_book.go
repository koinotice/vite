package dex

import (
	"github.com/koinotice/vite/interfaces"
	"github.com/koinotice/vite/vm_db"
)

type levelDbBook struct {
	db       vm_db.VmDb
	marketId int32
	side     bool
	iterator interfaces.StorageIterator
}

func getMakerBook(db vm_db.VmDb, marketId int32, side bool) (book *levelDbBook, err error) {
	book = &levelDbBook{db: db, marketId: marketId, side: side}
	book.iterator, err = db.NewStorageIterator(getBookPrefix(book))
	return
}

func (book *levelDbBook) nextOrder() (order *Order, ok bool) {
	if ok = book.iterator.Next(); ok {
		orderId := book.iterator.Key()
		orderData := book.iterator.Value()
		if len(orderId) != OrderIdBytesLength || len(orderData) == 0 {
			panic(IterateVmDbFailedErr)
		}
		order = &Order{}
		if err := order.DeSerializeCompact(orderData, orderId); err != nil {
			panic(err)
		}
	}
	return
}

func (book *levelDbBook) release() {
	book.iterator.Release()
}

func getBookPrefix(book *levelDbBook) []byte {
	if book.side {
		return append(Uint32ToBytes(uint32(book.marketId))[1:], byte(int8(1)))
	} else {
		return append(Uint32ToBytes(uint32(book.marketId))[1:], byte(int8(0)))
	}
}
