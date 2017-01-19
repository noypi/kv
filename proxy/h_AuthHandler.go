package proxy

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (this Server) validatePassword(pass string) bool {
	h := sha256.New()
	h.Write(this.passwordsalt)
	h.Write([]byte(pass))
	return this.passwordhash == fmt.Sprintf("%x", h.Sum(nil))
}

func (this *Server) Authenticate(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	session, err := this.sessions.Get(r, "kvproxy")
	if nil != err {
		return
	}
	bbPass, _ := ioutil.ReadAll(r.Body)
	if this.validatePassword(string(bbPass)) {
		session.Values["authenticated"] = true
		session.Save(r, w)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("not authorized"))
	}

}

func (this *Server) Logout(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	session, err := this.sessions.Get(r, "kvproxy")
	if nil != err {
		return
	}
	session.Values["authenticated"] = false
	for k, _ := range session.Values {
		delete(session.Values, k)
	}
	session.Save(r, w)

	this.gracesvr.Stop(1 * time.Second)

}
