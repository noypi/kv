package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/net/context"
)

func (this *Server) PutHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()
	if "PUT" != r.Method {
		err = fmt.Errorf("invalid method. use PUT.")
		return ctx
	}

	session := ctx.Value("session").(*sessions.Session)
	path := session.Values["path"].(string)

	key := r.FormValue("key")
	store, has := this.opendb[path]
	if !has {
		fmt.Errorf("db not opened.")
		return ctx
	}

	bb, err := ioutil.ReadAll(r.Body)
	if nil != err {
		return ctx
	}

	wrtr, _ := store.Writer()
	batch := wrtr.NewBatch()
	batch.Set([]byte(key), bb)
	err = wrtr.ExecuteBatch(batch)

	return ctx
}
