package proxy

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
)

func (this Server) validatePassword(pass string) bool {
	h := sha256.New()
	h.Write(this.passwordsalt)
	h.Write([]byte(pass))
	return this.passwordhash == fmt.Sprintf("%x", h.Sum(nil))
}

func (this *Server) AddSession(nexth http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		bValid, _ := session.Values["authenticated"].(bool)
		if !bValid {
			if this.validatePassword(r.FormValue("password")) {
				session.Values["authenticated"] = true
				session.Save(r, w)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("invalid session."))
				return
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "session", session)

		nexth(w, r.WithContext(ctx))
	}
}
