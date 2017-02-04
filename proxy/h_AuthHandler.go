package proxy

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/noypi/util"
	"github.com/noypi/webutil"
)

func (this Server) validatePassword(pass string) bool {
	h := sha256.New()
	h.Write(this.passwordsalt)
	h.Write([]byte(pass))
	return this.passwordhash == fmt.Sprintf("%x", h.Sum(nil))
}

func (this *Server) hRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (this *Server) hAuthPubkey(w http.ResponseWriter, r *http.Request) {
	if this.bUseTls {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bbPubk, err := this.privkey.PubKey().Marshal()
	if nil != err {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bbPubk)
}

func (this *Server) hAuthenticate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := ctx.Value(webutil.SessionName).(*sessions.Session)

	bbPass, _ := ioutil.ReadAll(r.Body)
	bbPass, _ = this.privkey.Decrypt(bbPass)
	if this.validatePassword(string(bbPass)) {
		session.Values["$authenticated"] = true
		session.Save(r, w)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

}

func (this *Server) hValidate(nexth http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session := ctx.Value(webutil.SessionName).(*sessions.Session)

		bValid, _ := session.Values["$authenticated"].(bool)
		if !bValid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid session."))
			util.LogErr(ctx, "Server.Validate() unauthenticated user.")
			return
		}

		nexth.ServeHTTP(w, r)
	})
}

func (this *Server) hLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := ctx.Value(webutil.SessionName).(*sessions.Session)
	for k, _ := range session.Values {
		delete(session.Values, k)
	}
	session.Save(r, w)

}
