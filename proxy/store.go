package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	"github.com/gorilla/sessions"
	"github.com/noypi/kv"
)

type _store struct {
	store   kv.KVStore
	port    string
	client  *http.Client
	mo      kv.MergeOperator
	baseurl string

	sessions *sessions.CookieStore
	opendb   map[string]kv.KVStore
}

func NewClient(mo kv.MergeOperator, config map[string]interface{}) (kv.KVStore, error) {
	port, ok := config["port"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify port")
	}

	password, ok := config["password"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify password")
	}

	jar, _ := cookiejar.New(nil)
	rv := _store{
		port: port,
		client: &http.Client{
			Jar: jar,
		},
		mo: mo,
	}

	rv.baseurl = fmt.Sprintf("http://locahost:%s", port)
	_, err := rv.query(fmt.Sprintf("/open?password=%s", password))
	if nil != err {
		return nil, err
	}

	return &rv, nil
}

func (this *_store) Close() error {
	return nil
}

func (this *_store) Reader() (kv.KVReader, error) {
	panic("implement")
	return nil, nil
	/*return &_reader{
		store: this,
	}, nil*/
}

func (this *_store) Writer() (kv.KVWriter, error) {
	return nil, fmt.Errorf("Writer is not supported.")
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
