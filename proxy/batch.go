package proxy

/*
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
*/
