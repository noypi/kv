package proxy

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-zoo/bone"
	"github.com/gorilla/context"
	"github.com/noypi/kv"
	"github.com/noypi/util"
	. "github.com/noypi/webutil"
	"github.com/twinj/uuid"
	"gopkg.in/tylerb/graceful.v1"
)

type _openedReader struct {
	rdr       kv.KVReader
	iterators map[string]kv.KVIterator
}

type Server struct {
	passwordhash string
	passwordsalt []byte
	db           kv.KVStore
	gracesvr     *graceful.Server
	//server       *http.Server
	readers map[string]*_openedReader

	//sec
	bUseTls bool
	privkey *util.PrivKey

	syncDb   sync.Mutex
	syncRdrs sync.Mutex
}

func NewServer(store kv.KVStore, port int, password string, bUseTls bool) (server *Server, err error) {
	bbSecret := make([]byte, 10)
	if _, err := rand.Read(bbSecret); nil == err {
		bbSecret = []byte("some secret")
	}

	h := sha256.New()
	h.Write(bbSecret)
	h.Write([]byte(password))

	server = &Server{
		passwordhash: fmt.Sprintf("%x", h.Sum(nil)),
		db:           store,
		readers:      map[string]*_openedReader{},
		passwordsalt: bbSecret,
		bUseTls:      bUseTls,
	}

	sessionname := "kvproxy-" + uuid.NewV4().String()
	mux := bone.New()
	mux.HandleFunc("/", server.hRoot)
	mux.HandleFunc("/auth/pubkey", server.hAuthPubkey)
	mux.Handle("/auth", MidSeqFunc(
		server.hAuthenticate,
		MidFn(AddCookieSession, sessionname),
	))
	mux.Handle("/logout", MidSeqFunc(
		server.hLogout,
		MidFn(AddCookieSession, sessionname),
		MidFn(server.hValidate),
		MidFn(NoCache),
	))

	fnCommon := func(h http.HandlerFunc) http.Handler {
		return MidSeqFunc(h,
			MidFn(AddCookieSession, sessionname),
			MidFn(server.hValidate),
		)
	}

	mux.Handle("/stat", fnCommon(server.hStat))

	// reader
	mux.Handle("/reader/get", fnCommon(server.hReaderGetHandler))
	mux.Handle("/reader/multiget", fnCommon(server.hReaderMultiGetHandler))
	mux.Handle("/reader/new", fnCommon(server.hReaderNewHandler))
	mux.Handle("/reader/prefix", fnCommon(server.hReaderPrefixHandler))
	mux.Handle("/reader/range", fnCommon(server.hReaderRangeHandler))
	mux.Handle("/reader/close", fnCommon(server.hReaderCloseHandler))

	// iterator
	mux.Handle("/iter/seek", fnCommon(server.hIterSeekHandler))
	mux.Handle("/iter/close", fnCommon(server.hIterCloseHandler))
	mux.Handle("/iter/key", fnCommon(server.hIterKeyHandler))
	mux.Handle("/iter/value", fnCommon(server.hIterValueHandler))
	mux.Handle("/iter/valid", fnCommon(server.hIterValidHandler))
	mux.Handle("/iter/next", fnCommon(server.hIterNextHandler))

	//srv := &http.Server{Addr: ":" + port, Handler: context.ClearHandler(http.DefaultServeMux)}

	server.gracesvr = &graceful.Server{
		Timeout: 5 * time.Second,
		Server: &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: context.ClearHandler(mux),
		},
	}

	if server.bUseTls {
		keypath, certpath := generateTempCert("noypi", "localhost", 2048)
		go server.gracesvr.ListenAndServeTLS(certpath, keypath)
	} else {
		if server.privkey, err = util.GenPrivKey(2048); nil != err {
			log.Fatal("GenerateKey err=", err)
			return
		}
		go server.gracesvr.ListenAndServe()
	}

	return
}

func (this *Server) Close() {
	this.gracesvr.Stop(1 * time.Second)

	this.syncRdrs.Lock()
	for _, rdr := range this.readers {
		for _, iter := range rdr.iterators {
			iter.Close()
		}
		rdr.iterators = nil
		rdr.rdr.Close()
		rdr.rdr = nil
	}
	this.readers = nil
	this.syncRdrs.Unlock()
}

func (this *Server) getIter(rdrid, id string) (iter kv.KVIterator, has bool) {
	this.syncRdrs.Lock()
	defer this.syncRdrs.Unlock()
	openedReader, has := this.readers[rdrid]
	if !has {
		return
	}
	iter, has = openedReader.iterators[id]
	return
}
func (this *Server) newPrefixIter(rdrid string, prefix []byte) (iter kv.KVIterator, id string) {
	this.syncRdrs.Lock()
	defer this.syncRdrs.Unlock()

	id = fmt.Sprintf("%x", uuid.NewV4().Bytes())
	openedReader, _ := this.readers[rdrid]
	iter = openedReader.rdr.PrefixIterator(prefix)
	openedReader.iterators[id] = iter

	return
}

func (this *Server) newRangeIter(rdrid string, start, end []byte) (iter kv.KVIterator, id string) {
	this.syncRdrs.Lock()
	defer this.syncRdrs.Unlock()

	id = fmt.Sprintf("%x", uuid.NewV4().Bytes())
	openedReader, _ := this.readers[rdrid]
	iter = openedReader.rdr.RangeIterator(start, end)
	openedReader.iterators[id] = iter

	return
}

func (this *Server) getRdr(id string) (rdr kv.KVReader, has bool) {
	this.syncRdrs.Lock()
	defer this.syncRdrs.Unlock()
	openedReader, has := this.readers[id]
	if has {
		rdr = openedReader.rdr
	}
	return
}

func (this *Server) closeIter(rdrid, id string) {
	this.syncRdrs.Lock()
	if openedReader, has := this.readers[rdrid]; has {
		if iter, has := openedReader.iterators[id]; has {
			iter.Close()
			delete(openedReader.iterators, id)
		}
	}
	this.syncRdrs.Unlock()
}

func (this *Server) closeRdr(id string) {
	this.syncRdrs.Lock()
	if rdr, has := this.readers[id]; has {
		for _, iter := range rdr.iterators {
			iter.Close()
		}
		rdr.iterators = nil
		rdr.rdr.Close()
		delete(this.readers, id)
	}
	this.syncRdrs.Unlock()
}

func (this *Server) newRdr() (rdr kv.KVReader, id string, err error) {
	this.syncRdrs.Lock()
	defer this.syncRdrs.Unlock()
	id = fmt.Sprintf("%x", uuid.NewV4().Bytes())
	if rdr, err = this.db.Reader(); nil == err {
		this.readers[id] = &_openedReader{
			rdr:       rdr,
			iterators: map[string]kv.KVIterator{},
		}
	}

	return
}

func generateTempCert(org, hosts string, bits int) (keypath, certpath string) {
	var keyTmpF, certTmpF *os.File
	err := func() (err error) {
		certTmpF, err = ioutil.TempFile(".", "tmp-cert.pem")
		if nil != err {
			return
		}
		defer certTmpF.Close()

		keyTmpF, err = ioutil.TempFile(".", "tmp-key.pem")
		if nil != err {
			return
		}
		defer keyTmpF.Close()

		bbKey, bbCert, err := util.GenerateCert(org, hosts, bits)
		if nil != err {
			return
		}

		keyTmpF.Write(bbKey)
		certTmpF.Write(bbCert)
		keypath = keyTmpF.Name()
		certpath = certTmpF.Name()
		return nil
	}()
	if nil != err {
		log.Fatal("failed to create temp key file. err=", err)
	}

	return
}

type ServerStat struct {
	TLS           bool
	ReadersCount  int
	IteratorCount int
}

func (this *Server) Stat() *ServerStat {
	this.syncRdrs.Lock()
	nTotalIters := 0
	nTotalRdrs := 0
	for _, openedRdr := range this.readers {
		nTotalRdrs++
		nTotalIters += len(openedRdr.iterators)
	}
	this.syncRdrs.Unlock()

	return &ServerStat{
		TLS:           this.bUseTls,
		ReadersCount:  nTotalRdrs,
		IteratorCount: nTotalIters,
	}
}

func (this ServerStat) MarshalJSON() (bb []byte, err error) {
	return []byte(fmt.Sprintf(`{"TLS": %v, "ReadersCount": %d, "IteratorCount": %d}`,
		this.TLS, this.ReadersCount, this.IteratorCount)), nil
}
