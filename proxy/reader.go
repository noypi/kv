package proxy

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
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
	buf := bytes.NewBufferString("")
	if err := gob.NewEncoder(buf).Encode(keys); nil != err {
		return nil, err
	}

	bb, err := r.store.postData(fmt.Sprintf("/reader/multiget?id=%s",
		r.ID,
	), buf.Bytes())
	if nil != err {
		return nil, err
	}

	var result = [][]byte{}
	bufDec := bytes.NewBuffer(bb)
	if err = gob.NewDecoder(bufDec).Decode(&result); nil != err {
		return nil, err
	}

	return result, nil
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
		rdrID:  r.ID,
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
		rdrID:  r.ID,
	}
	return &rv
}

func (r *_reader) Close() error {
	_, err := r.store.query(fmt.Sprintf("/reader/close?id=%s",
		r.ID,
	))
	return err
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
