package kv

type Setter interface {
	Close()
	Flush()
	Set(k, v []byte)
}

func NewSetter(kv KVStore) (setter Setter, err error) {
	o := new(_setter)
	wrtr, err := kv.Writer()
	if nil != err {
		return
	}

	batch := wrtr.NewBatch()
	o.set = func(k, v []byte) {
		batch.Set(k, v)
	}

	o.flush = func() {
		wrtr.ExecuteBatch(batch)
	}

	o.close = func() {
		batch.Close()
		wrtr.Close()
	}
	setter = o
	return
}

type _setter struct {
	close func()
	flush func()
	set   func(k, v []byte)
}

func (this *_setter) Close() {
	this.close()
}
func (this *_setter) Flush() {
	this.Flush()
}
func (this *_setter) Set(k, v []byte) {
	this.Set(k, v)
}
