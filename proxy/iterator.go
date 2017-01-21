package proxy

import (
	"encoding/base64"
	"fmt"
)

type _iterator struct {
	store  *_store
	iterID string
	rdrID  string
}

func (ldi *_iterator) Seek(key []byte) {
	ldi.store.query(fmt.Sprintf("/iter/seek?key=%s&rdrid=%s&id=%s",
		base64.RawURLEncoding.EncodeToString(key),
		ldi.rdrID,
		ldi.iterID,
	))
}

func (ldi *_iterator) Next() {
	ldi.store.query(fmt.Sprintf("/iter/next?rdrid=%s&id=%s", ldi.rdrID, ldi.iterID))
}

func (ldi *_iterator) Current() ([]byte, []byte, bool) {
	if ldi.Valid() {
		return ldi.Key(), ldi.Value(), true
	}
	return nil, nil, false
}

func (ldi *_iterator) Key() []byte {
	bbKey, err := ldi.store.query(fmt.Sprintf("/iter/key?rdrid=%s&id=%s", ldi.rdrID, ldi.iterID))
	if nil != err {
		return nil
	}
	return bbKey
}

func (ldi *_iterator) Value() []byte {
	bbVal, err := ldi.store.query(fmt.Sprintf("/iter/value?rdrid=%s&id=%s", ldi.rdrID, ldi.iterID))
	if nil != err {
		return nil
	}
	return bbVal
}

func (ldi *_iterator) Valid() bool {
	bbValid, err := ldi.store.query(fmt.Sprintf("/iter/valid?rdrid=%s&id=%s", ldi.rdrID, ldi.iterID))
	if nil != err || "true" != string(bbValid) {
		return false
	}
	return true
}

func (ldi *_iterator) Close() error {
	_, err := ldi.store.query(fmt.Sprintf("/iter/close?rdrid=%s&id=%s", ldi.rdrID, ldi.iterID))
	return err
}

func (ldi *_iterator) Count() int {
	panic("not implemented")
	return -1
}

func (ldi *_iterator) Error() error {
	panic("not implemented")
	return nil
}

func (ldi *_iterator) Reset() {
	panic("not implemented")
}

func (ldi *_iterator) Reset0() {
	panic("not implemented")
}
