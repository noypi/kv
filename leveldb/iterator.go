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
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type _iterator struct {
	store    *_store
	iterator iterator.Iterator
	reverse  bool
}

func (ldi *_iterator) Seek(key []byte) {
	ldi.iterator.Seek(key)
}

func (ldi *_iterator) Next() {
	if ldi.reverse {
		ldi.iterator.Prev()
	} else {
		ldi.iterator.Next()
	}
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
	ldi.Reset()
	n := 0
	for ; ldi.Valid(); ldi.Next() {
		n++
	}
	ldi.Reset()
	return n
}

func (ldi *_iterator) Error() error {
	panic("not implemented")
	return nil
}

func (ldi *_iterator) Reset() {
	if ldi.reverse {
		ldi.iterator.Last()
	} else {
		ldi.iterator.First()
	}
}

func (ldi *_iterator) Reset0() {
	panic("not implemented")
}
