// modified and copied from "github.com/blevesearch/bleve/index/store"

//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package kv

// blevesearch compatible interfaces

type KVStore interface {
	Writer() (KVWriter, error)
	Reader() (KVReader, error)
	Close() error
}

type KVReader interface {
	Get(key []byte) ([]byte, error)
	MultiGet(keys [][]byte) ([][]byte, error)
	PrefixIterator(prefix []byte) KVIterator
	PrefixIterator0(prefix []byte) KVIterator
	RangeIterator(start, end []byte) KVIterator
	RangeIterator0(start, end []byte) KVIterator

	ReversePrefixIterator(prefix []byte) KVIterator
	ReverseRangeIterator(start, end []byte) KVIterator
	Close() error
}

type KVIterator interface {
	Seek(key []byte)
	Next()
	Key() []byte
	Value() []byte
	Valid() bool
	Current() ([]byte, []byte, bool)
	Close() error
	Reset()
	Reset0()
	Error() error
	Count() int
}

type KVWriter interface {
	NewBatch() KVBatch
	NewBatchEx(KVBatchOptions) ([]byte, KVBatch, error)
	ExecuteBatch(batch KVBatch) error
	Close() error
}

type KVBatchOptions struct {
	TotalBytes int
	NumSets    int
	NumDeletes int
	NumMerges  int
}

type KVBatch interface {
	Set(key, val []byte)
	Delete(key []byte)
	Merge(key, val []byte)
	Reset()
	Close() error
}
