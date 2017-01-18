package proxy

import (
	"github.com/gorilla/sessions"
	"github.com/noypi/kv"
)

type _server struct {
	path         string
	passwordhash string
	sessions     *sessions.CookieStore
	opendb       map[string]kv.KVStore
}

/*

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	"crypto/sha256"

	"bitbucket.org/noypi/handlers"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/noypi/kv"
)

type _store struct {
	path    string
	baseurl string
	port    string
	client  *http.Client
	mo      kv.MergeOperator

	passwordhash string
	sessions     *sessions.CookieStore
	opendb       map[string]kv.KVStore
}

func GetCreateServer(port string) (kv.KVStore, error) {
	return NewServer(dummymergeop{}, map[string]interface{}{
		"port": port,
	})
}

func Get(path, pass, port string) (kv.KVStore, error) {
	return NewServer(dummymergeop{}, map[string]interface{}{
		"path":     path,
		"port":     port,
		"password": pass,
	})
}

func NewClient(mo kv.MergeOperator, config map[string]interface{}) (kv.KVStore, error) {
	path, ok := config["path"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify path")
	}

	port, ok := config["port"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify port")
	}

	password, ok := config["password"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify password")
	}

	h := sha256.New()
	h.Write([]byte(password))

	jar, _ := cookiejar.New(nil)
	rv := _store{
		path: path,
		port: port,
		client: &http.Client{
			Jar: jar,
		},
		mo:           mo,
		passwordhash: fmt.Sprintf("%x", h.Sum(nil)),
	}

	rv.baseurl = fmt.Sprintf("http://locahost:%s", port)
	_, err := rv.query(fmt.Sprintf("/open?path=%s", path))
	if nil != err {
		return nil, err
	}

	bbSecret := make([]byte, 10)
	if _, err := rand.Read(bbSecret); nil == err {
		bbSecret = []byte("some secret")
	}
	rv.sessions = sessions.NewCookieStore(bbSecret)

	return &rv, nil

}

func (this _store) query(q string) (bb []byte, err error) {
	resp, err := this.client.Get(fmt.Sprintf("%s%s", this.baseurl, q))
	if nil != err {
		return
	}
	defer resp.Body.Close()

	if 200 != resp.StatusCode {
		err = fmt.Errorf("%s", resp.Status)
		return
	}

	bb, err = ioutil.ReadAll(resp.Body)
	return
}

func NewServer(mo kv.MergeOperator, config map[string]interface{}) (kv.KVStore, error) {

	port, ok := config["port"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify port")
	}

	password, ok := config["password"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify password")
	}

	h := sha256.New()
	h.Write([]byte(password))

	rv := _server{
		passwordhash: fmt.Sprintf("%x", h.Sum(nil)),
		opendb:       map[string]kv.KVStore{},
	}

	http.Handle("/open", handlers.HttpSeq(
		rv.GetSessionHandler,
		rv.OpenHandler,
	))
	http.Handle("/get", handlers.HttpSeq(
		rv.GetSessionHandler,
		rv.GetHandler,
	))
	http.Handle("/put", handlers.HttpSeq(
		rv.GetSessionHandler,
		rv.PutHandler,
	))
	http.Handle("/quit", handlers.HttpSeq(
		rv.GetSessionHandler,
		rv.QuitHandler,
	))

	go http.ListenAndServe(":"+port, context.ClearHandler(http.DefaultServeMux))

	return &rv, nil
}

func (this *_store) Close() error {
	return nil
}

func (this *_store) Reader() (kv.KVReader, error) {
	return &_reader{
		store: this,
	}, nil
}

func (this *_store) Writer() (kv.KVWriter, error) {
	return &_writer{
		store: this,
	}, nil
}

type dummymergeop struct{}

func (this dummymergeop) FullMerge(key, existingValue []byte, operands [][]byte) ([]byte, bool) {
	return []byte{}, true
}

func (this dummymergeop) PartialMerge(key, leftOperand, rightOperand []byte) ([]byte, bool) {
	return []byte{}, true
}

func (this dummymergeop) Name() string {
	return "dummy-mergeop"
}
*/
