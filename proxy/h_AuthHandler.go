package proxy

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/noypi/webutil"
)

func (this Server) validatePassword(pass string) bool {
	h := sha256.New()
	h.Write(this.passwordsalt)
	h.Write([]byte(pass))
	return this.passwordhash == fmt.Sprintf("%x", h.Sum(nil))
}

func (this *Server) Authenticate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := ctx.Value(webutil.SessionName).(*sessions.Session)

	bbPass, _ := ioutil.ReadAll(r.Body)
	if this.validatePassword(string(bbPass)) {
		session.Values["authenticated"] = true
		session.Save(r, w)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("not authorized"))
	}

}

func (this *Server) Validate(nexth http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session := ctx.Value(webutil.SessionName).(*sessions.Session)

		bValid, _ := session.Values["authenticated"].(bool)
		if !bValid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid session."))
			webutil.LogErr(ctx, "Server.Validate() unauthenticated user.")
			return
		}

		nexth(w, r)
	}
}

func (this *Server) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := ctx.Value(webutil.SessionName).(*sessions.Session)
	for k, _ := range session.Values {
		delete(session.Values, k)
	}
	session.Save(r, w)

}
