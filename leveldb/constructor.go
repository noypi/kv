package leveldb

import (
	"github.com/noypi/kv"
)

type _leveldbKey int

const (
	Name _leveldbKey = iota
)

type MergeOperatorOpt struct {
	kv.Option
}

type FilePathOpt struct {
	kv.Option
}

type ConfigOpt struct {
	kv.Option
}

func newWithOpts(opts ...kv.Option) (kv.KVStore, error) {
	var mo kv.MergeOperator = dummymergeop{}
	var config0 map[string]interface{}
	var config = map[string]interface{}{}

	for _, opt := range opts {
		switch v := opt.(type) {
		case MergeOperatorOpt:
			mo = v.Option.(kv.MergeOperator)
		case ConfigOpt:
			config0 = v.Option.(map[string]interface{})
			if 0 < len(config) {
				for k, v := range config0 {
					config[k] = v
				}
			} else {
				config = config0
			}
		case FilePathOpt:
			config["path"] = filepath
		}
	}
	return New(mo, config)
}

func init() {
	kv.Register(Name, newWithOpts)
}
