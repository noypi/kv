package gtreap

import (
	"github.com/noypi/kv"
)

type _gtreapKey int

const (
	Name _gtreapKey = iota
)

func newWithOpts(opts ...kv.Option) (kv.KVStore, error) {
	var mo kv.MergeOperator = dummymergeop{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case kv.OptMergeOperator:
			mo = v.(kv.MergeOperator)
		}
	}
	return New(mo, nil)
}

func init() {
	kv.Register(Name, newWithOpts)
}
