package proxy

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"net/http"
)

func (this *Server) hReaderNewHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	_, id, err := this.newRdr()
	if nil != err {
		return
	}

	w.Write([]byte(id))
}

func (this *Server) hReaderGetHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	id := r.FormValue("id")

	bbKey, err := base64.RawURLEncoding.DecodeString(r.FormValue("key"))
	if nil != err {
		return
	}

	rdr, has := this.getRdr(id)
	if !has {
		fmt.Errorf("reader not found.")
		return
	}

	bb, err := rdr.Get(bbKey)
	if nil != err {
		return
	}

	w.Write(bb)
}

func (this *Server) hReaderMultiGetHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	id := r.FormValue("id")
	rdr, has := this.getRdr(id)
	if !has {
		fmt.Errorf("reader not found.")
		return
	}

	var bbKeys = [][]byte{}
	if err = gob.NewDecoder(r.Body).Decode(&bbKeys); nil != err {
		return
	}

	bbVals, err := rdr.MultiGet(bbKeys)
	if nil != err {
		return
	}

	buf := bytes.NewBufferString("")
	if err = gob.NewEncoder(buf).Encode(bbVals); nil != err {
		return
	}

	w.Write(buf.Bytes())
}

func (this *Server) hReaderCloseHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	this.closeRdr(id)
}

func (this *Server) hReaderPrefixHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	rdrid := r.FormValue("id")
	bbPrefix, err := base64.RawURLEncoding.DecodeString(r.FormValue("prefix"))
	if nil != err {
		return
	}

	if _, has := this.getRdr(rdrid); !has {
		err = fmt.Errorf("reader not found.")
		return
	}
	_, id := this.newPrefixIter(rdrid, bbPrefix)
	w.Write([]byte(id))
}

func (this *Server) hReaderRangeHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	rdrid := r.FormValue("id")
	bbStart, err := base64.RawURLEncoding.DecodeString(r.FormValue("start"))
	if nil != err {
		return
	}
	bbEnd, err := base64.RawURLEncoding.DecodeString(r.FormValue("end"))
	if nil != err {
		return
	}

	if _, has := this.getRdr(rdrid); !has {
		err = fmt.Errorf("reader not found.")
		return
	}
	_, id := this.newRangeIter(rdrid, bbStart, bbEnd)
	w.Write([]byte(id))
}
