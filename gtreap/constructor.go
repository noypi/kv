package gtreap

import (
	"github.com/noypi/kv"
)

type _gtreapKey int

const (
	Name _gtreapKey = iota
	DefaultName
)

type MergeOperatorOpt struct {
	kv.Option
}

func newDefault(opts ...kv.Option) (kv.KVStore, error) {
	return GetDefault(), nil
}

func newWithOpts(opts ...kv.Option) (kv.KVStore, error) {
	var mo kv.MergeOperator = dummymergeop{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case MergeOperatorOpt:
			mo = v.Option.(kv.MergeOperator)
		}
	}
	return New(mo, nil)
}

func init() {
	kv.Register(DefaultName, newDefault)
	kv.Register(Name, newWithOpts)
}
