package proxy

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	"github.com/noypi/kv"
)

type _store struct {
	store   kv.KVStore
	port    string
	client  *http.Client
	mo      kv.MergeOperator
	baseurl string
}

func NewClient(mo kv.MergeOperator, config map[string]interface{}) (prv kv.KVStore, err error) {
	port, ok := config["port"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify port")
	}

	password, ok := config["password"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify password")
	}

	bUseTls, _ := config["usetls"].(bool)

	jar, _ := cookiejar.New(nil)
	rv := _store{
		port: port,
		client: &http.Client{
			Jar: jar,
		},
		mo:      mo,
		baseurl: fmt.Sprintf("http://locahost:%s", port),
	}
	prv = &rv

	// if not using tls
	// get public key and encrypt password
	if !bUseTls {
		bbPubKey, err := rv.query("/auth/pubkey")
		if nil != err {
			return nil, err
		}
		hash := sha256.New()

		block, _ := pem.Decode(bbPubKey)
		pubKif, err := x509.ParsePKIXPublicKey(block.Bytes)
		if nil != err {
			return nil, err
		}
		bbPassEnc, err := rsa.EncryptOAEP(hash, rand.Reader, pubKif.(*rsa.PublicKey), []byte(password), []byte(""))
		if nil != err {
			return nil, err
		}

		_, err = rv.postData("/auth", bbPassEnc)
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

func (this _store) postData(url string, data []byte) (bb []byte, err error) {
	buf := bytes.NewBuffer(data)
	resp, err := this.client.Post(url, "application/octet-stream", buf)
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
