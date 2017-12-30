package kv

import (
	"fmt"
)

var ErrNotRegistered = fmt.Errorf("not registered.")

type Constructor func(...Option) (KVStore, error)

type Option interface{}
type ConfigType map[string]interface{}

type OptMergeOperator struct{ MergeOperator }
type OptConfig struct{ ConfigType }
type OptFilePath struct{ FilePath string }

type KVRegistry struct {
	m map[interface{}]Constructor
}

var defaultRegistry = new(KVRegistry)

func Register(k interface{}, c Constructor) {
	defaultRegistry.Register(k, c)
}

func New(k interface{}, opts ...Option) (KVStore, error) {
	return defaultRegistry.New(k, opts...)
}

func (this *KVRegistry) Register(k interface{}, c Constructor) {
	if nil == this.m {
		this.m = map[interface{}]Constructor{}
	}
	this.m[k] = c
}

func (this KVRegistry) New(k interface{}, opts ...Option) (KVStore, error) {
	if nil == this.m {
		return nil, ErrNotRegistered
	}

	constructor, has := this.m[k]
	if !has {
		return nil, ErrNotRegistered
	}
	return constructor(opts...)
}
