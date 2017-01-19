package proxy

import (
	"context"
	"net/http"
)

func (this *Server) AddVerifySession(name string, nexth http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if nil != err {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}
		}()

		session, err := this.sessions.Get(r, name)
		if nil != err {
			return
		}
		bValid, _ := session.Values["authenticated"].(bool)
		if !bValid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid session."))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "session", session)

		nexth(w, r.WithContext(ctx))
	}
}
