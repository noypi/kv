package kvutil

import (
	"log"
	"os"

	"fmt"

	"github.com/howeyc/gopass"
	"github.com/noypi/kv"
	"github.com/noypi/kv/leveldb"
	"github.com/noypi/kv/proxy"
)

func LoadLeveldbFallbackProxyPass(fpath string, proxyport int, proxyname, pass string) (store kv.KVStore, err error) {
	store, err = leveldb.GetDefault(fpath)
	if nil == err {
		// ok
		return
	}

	log.Println("Failed to open kv=", fpath, ", err=", err, ". trying to open proxy...")
	store, err = proxy.NewClient(proxyport, proxyname, pass, false)
	if nil != err {
		return
	}

	return
}

func LoadLeveldbFallbackProxy(fpath string, proxyport int, proxyname string) (store kv.KVStore, err error) {
	store, err = leveldb.GetDefault(fpath)
	if nil == err {
		// ok
		return
	}

	var pass []byte
	log.Println("Failed to open kv=", fpath, ", err=", err, ". trying to open proxy...")
	fmt.Fprintf(os.Stderr, "Password(%s): ", proxyname)
	pass, err = gopass.GetPasswdMasked()
	if nil != err {
		return
	}
	store, err = proxy.NewClient(proxyport, proxyname, string(pass), false)
	if nil != err {
		return
	}

	return
}
