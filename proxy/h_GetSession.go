package proxy

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"golang.org/x/net/context"
)

func (this _server) validatePassword(pass string) bool {
	h := sha256.New()
	h.Write([]byte(pass))
	return this.passwordhash == fmt.Sprintf("%x", h.Sum(nil))
}

func (this *_server) GetSessionHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	session, err := this.sessions.Get(r, "kvproxy")
	if nil != err {
		return ctx
	}
	bValid, _ := session.Values["authenticated"].(bool)
	if !bValid {
		if this.validatePassword(r.FormValue("password")) {
			session.Values["authenticated"] = true
			session.Save(r, w)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid session."))
			return ctx
		}
	}

	return ctx
}
