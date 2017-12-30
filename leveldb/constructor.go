package leveldb

import (
	"github.com/noypi/kv"
)

type _leveldbKey int

const (
	Name _leveldbKey = iota
)

func newWithOpts(opts ...kv.Option) (kv.KVStore, error) {
	var mo kv.MergeOperator = dummymergeop{}
	var config0 map[string]interface{}
	var config = map[string]interface{}{}

	for _, opt := range opts {
		switch v := opt.(type) {
		case kv.OptMergeOperator:
			mo = v.(kv.MergeOperator)
		case kv.OptConfig:
			config0 = map[string]interface{}(v)
			if 0 < len(config) {
				for k, v := range config0 {
					config[k] = v
				}
			} else {
				config = config0
			}
		case kv.OptFilePath:
			config["path"] = string(v)
		}
	}
	return New(mo, config)
}

func init() {
	kv.Register(Name, newWithOpts)
}
