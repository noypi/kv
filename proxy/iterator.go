package proxy

/*
import (
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type _iterator struct {
	store    *_store
	iterator iterator.Iterator
}

func (ldi *_iterator) Seek(key []byte) {
	ldi.iterator.Seek(key)
}

func (ldi *_iterator) Next() {
	ldi.iterator.Next()
}

func (ldi *_iterator) Current() ([]byte, []byte, bool) {
	if ldi.Valid() {
		return ldi.Key(), ldi.Value(), true
	}
	return nil, nil, false
}

func (ldi *_iterator) Key() []byte {
	return ldi.iterator.Key()
}

func (ldi *_iterator) Value() []byte {
	return ldi.iterator.Value()
}

func (ldi *_iterator) Valid() bool {
	return ldi.iterator.Valid()
}

func (ldi *_iterator) Close() error {
	ldi.iterator.Release()
	return nil
}

func (ldi *_iterator) Count() int {
	panic("not implemented")
	return -1
}

func (ldi *_iterator) Error() error {
	panic("not implemented")
	return nil
}

func (ldi *_iterator) Reset() {
	panic("not implemented")
}

func (ldi *_iterator) Reset0() {
	panic("not implemented")
}
*/
