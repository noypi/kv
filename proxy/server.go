package proxy

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/context"
	"github.com/noypi/kv"
	. "github.com/noypi/webutil"
	"github.com/twinj/uuid"
	"gopkg.in/tylerb/graceful.v1"
)

type Server struct {
	passwordhash string
	passwordsalt []byte
	db           kv.KVStore
	gracesvr     *graceful.Server
	//server       *http.Server
	iterators map[string]kv.KVIterator
	readers   map[string]kv.KVReader

	syncDb    sync.Mutex
	syncIters sync.Mutex
	syncRdrs  sync.Mutex
}

func NewServer(store kv.KVStore, port, password string) (*Server, error) {

	bbSecret := make([]byte, 10)
	if _, err := rand.Read(bbSecret); nil == err {
		bbSecret = []byte("some secret")
	}

	h := sha256.New()
	h.Write(bbSecret)
	h.Write([]byte(password))

	server := &Server{
		passwordhash: fmt.Sprintf("%x", h.Sum(nil)),
		db:           store,
		iterators:    map[string]kv.KVIterator{},
		readers:      map[string]kv.KVReader{},
		passwordsalt: bbSecret,
	}

	sessionname := "kvproxy-" + uuid.NewV4().String()
	http.HandleFunc("/auth", server.Authenticate)
	http.Handle("/logout", MidSeqFunc(
		server.Logout,
		MidFn(AddCookieSession, sessionname),
		MidFn(server.Validate),
		MidFn(NoCache),
	))

	fnCommon := func(h http.HandlerFunc) http.Handler {
		return MidSeqFunc(h,
			MidFn(AddCookieSession, sessionname),
			MidFn(server.Validate),
		)
	}

	// reader
	http.Handle("/reader/get", fnCommon(server.ReaderGetHandler))
	http.Handle("/reader/new", fnCommon(server.ReaderNewHandler))
	http.Handle("/reader/prefix", fnCommon(server.ReaderPrefixHandler))
	http.Handle("/reader/range", fnCommon(server.ReaderRangeHandler))

	// iterator
	http.Handle("/iter/seek", fnCommon(server.IterSeekHandler))
	http.Handle("/iter/close", fnCommon(server.IterCloseHandler))
	http.Handle("/iter/key", fnCommon(server.IterKeyHandler))
	http.Handle("/iter/value", fnCommon(server.IterValueHandler))
	http.Handle("/iter/valid", fnCommon(server.IterValidHandler))
	http.Handle("/iter/next", fnCommon(server.IterNextHandler))

	//srv := &http.Server{Addr: ":" + port, Handler: context.ClearHandler(http.DefaultServeMux)}

	keypath, certpath := generateTempCert("noypi", "localhost", 2048)

	server.gracesvr = &graceful.Server{
		Timeout: 5 * time.Second,
		Server: &http.Server{
			Addr:    ":" + port,
			Handler: context.ClearHandler(http.DefaultServeMux),
		},
	}

	go log.Fatal(server.gracesvr.ListenAndServeTLS(certpath, keypath))

	return server, nil
}

func (this *Server) Close() {
	this.syncIters.Lock()
	for _, iter := range this.iterators {
		iter.Close()
	}
	this.syncIters.Unlock()

	this.syncRdrs.Lock()
	for _, rdr := range this.readers {
		rdr.Close()
	}
	this.syncRdrs.Unlock()
}

func (this *Server) getIter(id string) (iter kv.KVIterator, has bool) {
	this.syncIters.Lock()
	defer this.syncIters.Unlock()
	iter, has = this.iterators[id]
	return
}
func (this *Server) newPrefixIter(rdr kv.KVReader, prefix []byte) (iter kv.KVIterator, id string) {
	this.syncIters.Lock()
	defer this.syncIters.Unlock()
	id = fmt.Sprintf("%x", uuid.NewV4().Bytes())

	iter = rdr.PrefixIterator(prefix)
	this.iterators[id] = iter

	return
}

func (this *Server) newRangeIter(rdr kv.KVReader, start, end []byte) (iter kv.KVIterator, id string) {
	this.syncIters.Lock()
	defer this.syncIters.Unlock()
	id = fmt.Sprintf("%x", uuid.NewV4().Bytes())

	iter = rdr.RangeIterator(start, end)
	this.iterators[id] = iter

	return
}

func (this *Server) getRdr(id string) (rdr kv.KVReader, has bool) {
	this.syncRdrs.Lock()
	defer this.syncRdrs.Unlock()
	rdr, has = this.readers[id]
	return
}

func (this *Server) closeIter(id string) {
	this.syncIters.Lock()
	if iter, has := this.iterators[id]; has {
		iter.Close()
		delete(this.iterators, id)
	}
	this.syncIters.Unlock()
}

func (this *Server) closeRdr(id string) {
	this.syncRdrs.Lock()
	if rdr, has := this.readers[id]; has {
		rdr.Close()
		delete(this.readers, id)
	}
	this.syncRdrs.Unlock()
}

func (this *Server) newRdr() (rdr kv.KVReader, id string, err error) {
	this.syncRdrs.Lock()
	defer this.syncRdrs.Unlock()
	id = fmt.Sprintf("%x", uuid.NewV4().Bytes())
	if rdr, err = this.db.Reader(); nil == err {
		this.readers[id] = rdr
	}

	return
}

func generateTempCert(org, hosts string, bits int) (keypath, certpath string) {
	priv, err := rsa.GenerateKey(rand.Reader, bits)

	tValidTo := time.Now().AddDate(1, 0, 0)
	tValidFrom := time.Now().AddDate(-1, 0, 0)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatal("failed to generate serial number err=", err)
	}
	tmpl := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{org},
		},
		NotBefore:             tValidFrom,
		NotAfter:              tValidTo,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA: true,
	}
	for _, h := range strings.Split(hosts, ",") {
		if ip := net.ParseIP(h); nil != ip {
			tmpl.IPAddresses = append(tmpl.IPAddresses, ip)
		} else {
			tmpl.DNSNames = append(tmpl.DNSNames, h)
		}
	}

	bb, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
	if nil != err {
		log.Fatal("failed to create certificate err=", err)
	}

	certTmpF, err := ioutil.TempFile(".", "tmp-cert.pem")
	if nil != err {
		log.Fatal("failed to create temp cert file. err=", err)
	}
	defer certTmpF.Close()
	pem.Encode(certTmpF, &pem.Block{Type: "CERTIFICATE", Bytes: bb})

	keyTmpF, err := ioutil.TempFile(".", "tmp-key.pem")
	if nil != err {
		log.Fatal("failed to create temp key file. err=", err)
	}
	defer keyTmpF.Close()
	pem.Encode(keyTmpF, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return keyTmpF.Name(), certTmpF.Name()
}
