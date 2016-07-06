// modified and copied from "github.com/blevesearch/bleve/index/store"

//  Copyright (c) 2016 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package kv

// MultiGet is a helper function to retrieve mutiple keys from a
// KVReader, and might be used by KVStore implementations that don't
// have a native multi-get facility.
func MultiGet(kvreader KVReader, keys [][]byte) ([][]byte, error) {
	vals := make([][]byte, 0, len(keys))

	for i, key := range keys {
		val, err := kvreader.Get(key)
		if err != nil {
			return nil, err
		}

		vals[i] = val
	}

	return vals, nil
}