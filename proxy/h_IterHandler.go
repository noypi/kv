package proxy

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

func (this *Server) hIterSeekHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	bbKey, err := base64.RawURLEncoding.DecodeString(r.FormValue("key"))
	if nil != err {
		return
	}

	id := r.FormValue("id")
	rdrid := r.FormValue("rdrid")

	iter, has := this.getIter(rdrid, id)
	if !has {
		err = fmt.Errorf("iterator not found.")
		return
	}

	iter.Seek(bbKey)
}

func (this *Server) hIterKeyHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	id := r.FormValue("id")
	rdrid := r.FormValue("rdrid")
	iter, has := this.getIter(rdrid, id)
	if !has {
		err = fmt.Errorf("iterator not found.")
		return
	}

	w.Write(iter.Key())
}

func (this *Server) hIterValidHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()
	id := r.FormValue("id")
	rdrid := r.FormValue("rdrid")

	iter, has := this.getIter(rdrid, id)
	if !has {
		err = fmt.Errorf("iterator not found.")
		return
	}

	if iter.Valid() {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
}

func (this *Server) hIterValueHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()
	id := r.FormValue("id")
	rdrid := r.FormValue("rdrid")

	iter, has := this.getIter(rdrid, id)
	if !has {
		err = fmt.Errorf("iterator not found.")
		return
	}

	w.Write(iter.Value())
}

func (this *Server) hIterNextHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()
	id := r.FormValue("id")
	rdrid := r.FormValue("rdrid")

	iter, has := this.getIter(rdrid, id)
	if !has {
		err = fmt.Errorf("iterator not found.")
		return
	}

	iter.Next()
}

func (this *Server) hIterCloseHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()
	id := r.FormValue("id")
	rdrid := r.FormValue("rdrid")
	this.closeIter(rdrid, id)
}
