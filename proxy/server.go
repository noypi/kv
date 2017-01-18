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
	"time"

	"bitbucket.org/noypi/handlers"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/noypi/kv"
)

type Server struct {
	path         string
	passwordhash string
	passwordsalt []byte
	server       *http.Server
	sessions     *sessions.CookieStore
	opendb       map[string]kv.KVStore
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
		opendb:       map[string]kv.KVStore{},
		passwordsalt: bbSecret,
	}

	http.Handle("/get", handlers.HttpSeq(
		server.GetSessionHandler,
		server.GetHandler,
	))
	http.Handle("/put", handlers.HttpSeq(
		server.GetSessionHandler,
		server.PutHandler,
	))

	srv := &http.Server{Addr: ":" + port, Handler: context.ClearHandler(http.DefaultServeMux)}
	server.sessions = sessions.NewCookieStore(bbSecret)

	keypath, certpath := generateTempCert("noypi", "localhost", 2048)
	go log.Fatal(srv.ListenAndServeTLS(certpath, keypath))

	return server, nil
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
