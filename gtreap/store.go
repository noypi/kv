// modified and copied from "github.com/blevesearch/bleve/index/store/gtreap"

//  Copyright (c) 2015 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the
//  License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on an "AS
//  IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
//  express or implied. See the License for the specific language
//  governing permissions and limitations under the License.

// Package gtreap provides an in-memory implementation of the
// KVStore interfaces using the gtreap balanced-binary treap,
// copy-on-write data structure.

package gtreap

import (
	"bytes"
	"sync"

	"github.com/noypi/kv"
	"github.com/steveyen/gtreap"
)

const Name = "gtreap"

type Store struct {
	m  sync.Mutex
	t  *gtreap.Treap
	mo kv.MergeOperator
}

type Item struct {
	k []byte
	v []byte
}

func itemCompare(a, b interface{}) int {
	return bytes.Compare(a.(*Item).k, b.(*Item).k)
}

func GetDefault() kv.KVStore {
	store, _ := New(dummymergeop{}, nil)
	return store
}

func New(mo kv.MergeOperator, config map[string]interface{}) (kv.KVStore, error) {
	rv := Store{
		t:  gtreap.NewTreap(itemCompare),
		mo: mo,
	}
	return &rv, nil
}

func (s *Store) Close() error {
	return nil
}

func (s *Store) Reader() (kv.KVReader, error) {
	s.m.Lock()
	t := s.t
	s.m.Unlock()
	return &Reader{t: t}, nil
}

func (s *Store) Writer() (kv.KVWriter, error) {
	return &Writer{s: s}, nil
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
