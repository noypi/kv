package proxy

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

func (this *Server) IterSeekHandler(w http.ResponseWriter, r *http.Request) {
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

	iter, has := this.getIter(id)
	if !has {
		err = fmt.Errorf("iterator not found.")
		return
	}

	iter.Seek(bbKey)
}

func (this *Server) IterKeyHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	id := r.FormValue("id")
	iter, has := this.getIter(id)
	if !has {
		err = fmt.Errorf("iterator not found.")
		return
	}

	w.Write(iter.Key())
}

func (this *Server) IterValidHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()
	id := r.FormValue("id")

	iter, has := this.getIter(id)
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

func (this *Server) IterValueHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()
	id := r.FormValue("id")

	iter, has := this.getIter(id)
	if !has {
		err = fmt.Errorf("iterator not found.")
		return
	}

	w.Write(iter.Value())
}

func (this *Server) IterNextHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()
	id := r.FormValue("id")

	iter, has := this.getIter(id)
	if !has {
		err = fmt.Errorf("iterator not found.")
		return
	}

	iter.Next()
}

func (this *Server) IterCloseHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()
	id := r.FormValue("id")
	this.closeIter(id)
}
