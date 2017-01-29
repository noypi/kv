// modified and copied from "github.com/blevesearch/bleve/index/store/goleveldb"

//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package leveldb

import (
	"github.com/noypi/kv"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type _reader struct {
	store    *_store
	snapshot *leveldb.Snapshot
}

func (r *_reader) Get(key []byte) ([]byte, error) {
	b, err := r.snapshot.Get(key, r.store.defaultReadOptions)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	return b, err
}

func (r *_reader) MultiGet(keys [][]byte) ([][]byte, error) {
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
	byteRange := util.BytesPrefix(prefix)
	iter := r.snapshot.NewIterator(byteRange, r.store.defaultReadOptions)
	iter.First()
	rv := _iterator{
		store:    r.store,
		iterator: iter,
	}
	return &rv
}

func (r *_reader) RangeIterator(start, end []byte) kv.KVIterator {
	byteRange := &util.Range{
		Start: start,
		Limit: end,
	}
	iter := r.snapshot.NewIterator(byteRange, r.store.defaultReadOptions)
	iter.First()
	rv := _iterator{
		store:    r.store,
		iterator: iter,
	}
	return &rv
}

func (r *_reader) Close() error {
	r.snapshot.Release()
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
	iter := r.PrefixIterator(prefix).(*_iterator)
	iter.reverse = true
	iter.iterator.Last()
	return iter
}

func (r *_reader) ReverseRangeIterator(start, end []byte) kv.KVIterator {
	iter := r.ReverseRangeIterator(start, end).(*_iterator)
	iter.reverse = true
	iter.iterator.Last()
	return iter
}
