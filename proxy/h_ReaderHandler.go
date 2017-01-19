package proxy

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

func (this *Server) ReaderNewHandler(w http.ResponseWriter, r *http.Request) {
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

func (this *Server) ReaderGetHandler(w http.ResponseWriter, r *http.Request) {
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

func (this *Server) ReaderCloseHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	id := r.FormValue("id")
	this.closeRdr(id)
}

func (this *Server) ReaderPrefixHandler(w http.ResponseWriter, r *http.Request) {
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

	rdr, has := this.getRdr(rdrid)
	if !has {
		err = fmt.Errorf("reader not found.")
		return
	}
	_, id := this.newPrefixIter(rdr, bbPrefix)
	w.Write([]byte(id))
}

func (this *Server) ReaderRangeHandler(w http.ResponseWriter, r *http.Request) {
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

	rdr, has := this.getRdr(rdrid)
	if !has {
		err = fmt.Errorf("reader not found.")
		return
	}
	_, id := this.newRangeIter(rdr, bbStart, bbEnd)
	w.Write([]byte(id))
}
