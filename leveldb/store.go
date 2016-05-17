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
	"fmt"

	"github.com/noypi/kv"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type _store struct {
	path string
	opts *opt.Options
	db   *leveldb.DB
	mo   kv.MergeOperator

	defaultWriteOptions *opt.WriteOptions
	defaultReadOptions  *opt.ReadOptions
}

func GetDefault(path string) (kv.KVStore, error) {
	return New(dummymergeop{}, map[string]interface{}{
		"path": path,
	})
}

func New(mo kv.MergeOperator, config map[string]interface{}) (kv.KVStore, error) {

	path, ok := config["path"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify path")
	}

	opts, err := applyConfig(&opt.Options{}, config)
	if err != nil {
		return nil, err
	}

	db, err := leveldb.OpenFile(path, opts)
	if err != nil {
		return nil, err
	}

	rv := _store{
		path:                path,
		opts:                opts,
		db:                  db,
		mo:                  mo,
		defaultReadOptions:  &opt.ReadOptions{},
		defaultWriteOptions: &opt.WriteOptions{},
	}
	rv.defaultWriteOptions.Sync = true
	return &rv, nil
}

func (ldbs *_store) Close() error {
	return ldbs.db.Close()
}

func (ldbs *_store) Reader() (kv.KVReader, error) {
	snapshot, _ := ldbs.db.GetSnapshot()
	return &_reader{
		store:    ldbs,
		snapshot: snapshot,
	}, nil
}

func (ldbs *_store) Writer() (kv.KVWriter, error) {
	return &_writer{
		store: ldbs,
	}, nil
}

func applyConfig(o *opt.Options, config map[string]interface{}) (*opt.Options, error) {

	ro, ok := config["read_only"].(bool)
	if ok {
		o.ReadOnly = ro
	}

	cim, ok := config["create_if_missing"].(bool)
	if ok {
		o.ErrorIfMissing = !cim
	}

	eie, ok := config["error_if_exists"].(bool)
	if ok {
		o.ErrorIfExist = eie
	}

	wbs, ok := config["write_buffer_size"].(float64)
	if ok {
		o.WriteBuffer = int(wbs)
	}

	bs, ok := config["block_size"].(float64)
	if ok {
		o.BlockSize = int(bs)
	}

	bri, ok := config["block_restart_interval"].(float64)
	if ok {
		o.BlockRestartInterval = int(bri)
	}

	lcc, ok := config["lru_cache_capacity"].(float64)
	if ok {
		o.BlockCacheCapacity = int(lcc)
	}

	bfbpk, ok := config["bloom_filter_bits_per_key"].(float64)
	if ok {
		bf := filter.NewBloomFilter(int(bfbpk))
		o.Filter = bf
	}

	return o, nil
}

type dummymergeop struct{}

func (this dummymergeop) FullMerge(key, existingValue []byte, operands [][]byte) ([]byte, bool) {
	return []byte{}, true
}

func (this dummymergeop) PartialMerge(key, leftOperand, rightOperand []byte) ([]byte, bool) {
	return []byte{}, true
}

func (this dummymergeop) Name() string {
	return "dummy-mergeop"
}
