package proxy

import (
	"net/http"

	"golang.org/x/net/context"
)

func (this *_server) QuitHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	var err error
	defer func() {
		if nil != err {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}()

	for _, store := range this.opendb {
		store.Close()
	}

	return ctx
}
