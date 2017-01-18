package proxy

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/net/context"
)

func (this *_server) GetHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	session := ctx.Value("session").(*sessions.Session)
	path := session.Values["path"].(string)

	key := r.FormValue("key")
	store, has := this.opendb[path]
	if !has {
		fmt.Errorf("db not opened.")
		return ctx
	}

	rdr, _ := store.Reader()
	bb, err := rdr.Get([]byte(key))
	if nil != err {
		return ctx
	}
	w.Write(bb)
	return ctx
}
