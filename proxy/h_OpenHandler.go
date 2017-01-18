package proxy

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/noypi/kv/leveldb"
	"golang.org/x/net/context"
)

func (this *Server) OpenHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	session := ctx.Value("session").(*sessions.Session)

	path := r.FormValue("path")
	store, has := this.opendb[path]
	if !has {
		if store, err = leveldb.GetDefault(path); nil != err {
			return ctx
		} else {
			this.opendb[path] = store
		}
	}

	ctx = context.WithValue(ctx, "kvstore", store)
	session.Values["path"] = path
	session.Save(r, w)
	return ctx
}
