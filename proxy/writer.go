package proxy

/*
import (
	"fmt"

	"github.com/noypi/kv"
	"github.com/syndtr/goleveldb/leveldb"
)

type _writer struct {
	store *_store
}

func (w *_writer) NewBatch() kv.KVBatch {
	rv := _batch{
		store: w.store,
		merge: kv.NewEmulatedMerge(w.store.mo),
		batch: new(leveldb.Batch),
	}
	return &rv
}

func (w *_writer) NewBatchEx(options kv.KVBatchOptions) ([]byte, kv.KVBatch, error) {
	return make([]byte, options.TotalBytes), w.NewBatch(), nil
}

func (w *_writer) ExecuteBatch(b kv.KVBatch) error {
	batch, ok := b.(*_batch)
	if !ok {
		return fmt.Errorf("wrong type of batch")
	}

	// first process merges
	for k, mergeOps := range batch.merge.Merges {
		kb := []byte(k)
		existingVal, err := w.store.db.Get(kb, w.store.defaultReadOptions)
		if err != nil && err != leveldb.ErrNotFound {
			return err
		}
		mergedVal, fullMergeOk := w.store.mo.FullMerge(kb, existingVal, mergeOps)
		if !fullMergeOk {
			return fmt.Errorf("merge operator returned failure")
		}
		// add the final merge to this batch
		batch.batch.Put(kb, mergedVal)
	}

	// now execute the batch
	return w.store.db.Write(batch.batch, w.store.defaultWriteOptions)
}

func (w *_writer) Close() error {
	return nil
}
*/
