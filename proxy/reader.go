package proxy

/*
import (
	"fmt"

	"github.com/noypi/kv"
)

type _reader struct {
	store *_store
}

func (r *_reader) Get(key []byte) ([]byte, error) {
	bb, err := r.store.query(fmt.Sprintf("/get?key=%s", string(key)))
	return bb, err
}

func (r *_reader) MultiGet(keys [][]byte) ([][]byte, error) {
	panic("need to implement")
	vals := make([][]byte, len(keys))

	for i, key := range keys {
		val, err := r.Get(key)
		if err != nil {
			return nil, err
		}

		vals[i] = val
	}

	return vals, nil
}

func (r *_reader) PrefixIterator(prefix []byte) kv.KVIterator {
	panic("need to implement")
	rv := _iterator{
		store:    r.store,
		iterator: nil,
	}
	return &rv
}

func (r *_reader) RangeIterator(start, end []byte) kv.KVIterator {
	panic("need to implement")
	rv := _iterator{
		store:    r.store,
		iterator: nil,
	}
	return &rv
}

func (r *_reader) Close() error {
	panic("need to implement")
	return nil
}

func (r *_reader) PrefixIterator0(prefix []byte) kv.KVIterator {
	panic("not implemented")
	return nil
}

func (r *_reader) RangeIterator0(start, end []byte) kv.KVIterator {
	panic("not implemented")
	return nil
}

func (r *_reader) ReversePrefixIterator(prefix []byte) kv.KVIterator {
	panic("not implemented")
	return nil
}

func (r *_reader) ReverseRangeIterator(start, end []byte) kv.KVIterator {
	panic("not implemented")
	return nil
}
*/
