package proxy

import (
	"encoding/json"
	"net/http"
)

func (this *Server) hStat(w http.ResponseWriter, r *http.Request) {
	bb, err := json.Marshal(this.Stat())
	if nil != err {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(bb)
}
