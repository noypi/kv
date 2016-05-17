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
)

type _batch struct {
	store *_store
	merge *kv.EmulatedMerge
	batch *leveldb.Batch
}

func (b *_batch) Set(key, val []byte) {
	b.batch.Put(key, val)
}

func (b *_batch) Delete(key []byte) {
	b.batch.Delete(key)
}

func (b *_batch) Merge(key, val []byte) {
	b.merge.Merge(key, val)
}

func (b *_batch) Reset() {
	b.batch.Reset()
	b.merge = kv.NewEmulatedMerge(b.store.mo)
}

func (b *_batch) Close() error {
	b.batch.Reset()
	b.batch = nil
	b.merge = nil
	return nil
}
