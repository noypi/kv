package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strconv"

	"github.com/noypi/kv"
	"github.com/noypi/util"
)

type _store struct {
	store   kv.KVStore
	port    string
	client  *http.Client
	mo      kv.MergeOperator
	baseurl string
}

func NewClient(port int, basename, password string, bUseTls bool) (prv kv.KVStore, err error) {
	return New(dummymergeop{}, map[string]interface{}{
		"password": password,
		"port":     port,
		"usetls":   bUseTls,
		"basename": basename,
	})
}

func New(mo kv.MergeOperator, config map[string]interface{}) (prv kv.KVStore, err error) {
	port, ok := config["port"].(int)
	if !ok {
		return nil, fmt.Errorf("must specify port")
	}
	basename, ok := config["basename"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify basename")
	}
	password, ok := config["password"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify password")
	}

	bUseTls, _ := config["usetls"].(bool)

	jar, _ := cookiejar.New(nil)
	rv := _store{
		port: strconv.Itoa(port),
		client: &http.Client{
			Jar: jar,
		},
		mo:      mo,
		baseurl: fmt.Sprintf("http://localhost:%d/%s", port, basename),
	}
	prv = &rv

	// if not using tls
	// get public key and encrypt password
	if !bUseTls {
		bbPubKey, err := rv.query("/auth/pubkey")
		if nil != err {
			return nil, err
		}
		pubk, err := util.ParsePublickey(bbPubKey)
		if nil != err {
			return nil, err
		}
		bbPassCipher, err := pubk.Encrypt([]byte(password))
		if nil != err {
			return nil, err
		}

		_, err = rv.postData("/auth", bbPassCipher)
		if nil != err {
			return nil, err
		}

	} else {
		if _, err = rv.postData("/auth", []byte(password)); nil != err {
			return
		}
	}

	return
}

func Stat(kvs kv.KVStore) (stat *ServerStat, err error) {
	store, ok := kvs.(*_store)
	if !ok {
		err = fmt.Errorf("not a proxy kv")
		return
	}

	stat = new(ServerStat)
	bb, err := store.query("/stat")
	if nil != err {
		return
	}

	if err = json.Unmarshal(bb, stat); nil != err {
		return
	}

	return
}

func (this *_store) Close() error {
	return nil
}

func (this *_store) Reader() (kv.KVReader, error) {
	bb, err := this.query("/reader/new")
	if nil != err {
		return nil, err
	}
	rv := _reader{
		store: this,
		ID:    string(bb),
	}
	return &rv, nil
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

func respAsError(resp *http.Response, def error) (err error) {
	if nil == resp {
		err = def
		return
	}

	bValid := 100 < resp.StatusCode && resp.StatusCode < 300
	if !bValid {
		var content string
		if nil != resp {
			bb, _ := ioutil.ReadAll(resp.Body)
			content = string(bb)
		}
		err = fmt.Errorf("%s.%s", resp.Status, content)
	}

	return
}

func (this _store) query(q string) (bb []byte, err error) {
	resp, err := this.client.Get(fmt.Sprintf("%s%s", this.baseurl, q))
	if nil != err {
		return
	}
	defer resp.Body.Close()

	if err = respAsError(resp, nil); nil != err {
		return
	}

	bb, err = ioutil.ReadAll(resp.Body)
	return
}

func (this _store) postData(url string, data []byte) (bb []byte, err error) {
	buf := bytes.NewBuffer(data)
	resp, err := this.client.Post(this.baseurl+url, "application/octet-stream", buf)
	if nil != err {
		return
	}
	defer resp.Body.Close()

	if err = respAsError(resp, nil); nil != err {
		return
	}

	bb, err = ioutil.ReadAll(resp.Body)
	return
}
