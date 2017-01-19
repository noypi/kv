package proxy

import (
	"encoding/base64"
	"fmt"

	"github.com/noypi/kv"
)

type _reader struct {
	store *_store
	ID    string
}

func (r *_reader) Get(key []byte) ([]byte, error) {
	bb, err := r.store.query(fmt.Sprintf("/reader/get?key=%s&id=%s",
		base64.RawURLEncoding.EncodeToString(key),
		r.ID,
	))
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
	bb, err := r.store.query(fmt.Sprintf("/reader/prefix?prefix=%s&id=%s",
		base64.RawURLEncoding.EncodeToString(prefix),
		r.ID,
	))
	if nil != err {
		return nil
	}
	rv := _iterator{
		store:  r.store,
		iterID: string(bb),
	}
	return &rv
}

func (r *_reader) RangeIterator(start, end []byte) kv.KVIterator {
	bb, err := r.store.query(fmt.Sprintf("/reader/range?start=%s&end=%s&id=%s",
		base64.RawURLEncoding.EncodeToString(start),
		base64.RawURLEncoding.EncodeToString(end),
		r.ID,
	))
	if nil != err {
		return nil
	}
	rv := _iterator{
		store:  r.store,
		iterID: string(bb),
	}
	return &rv
}

func (r *_reader) Close() error {
	panic("not implemented")
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
