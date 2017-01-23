package proxy_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/noypi/kv"
	"github.com/noypi/kv/leveldb"
	"github.com/noypi/kv/proxy"
	assertpkg "github.com/stretchr/testify/assert"
)

type kvtype struct {
	k, v []byte
}

var g_testTable []kvtype

func init() {
	g_testTable = []kvtype{
		kvtype{[]byte("somek0"), []byte("someval0")},
		kvtype{[]byte("somek1"), []byte("someval1")},
		kvtype{[]byte("somek2"), []byte("someval2")},
		kvtype{[]byte("somek3"), []byte("someval3")},
	}

}

func set01(assert *assertpkg.Assertions) (srv *proxy.Server, clientStore kv.KVStore, closer func()) {
	tmpdir := os.TempDir() + "/testing-proxyserver"
	os.RemoveAll(tmpdir)

	kvstore, err := leveldb.GetDefault(tmpdir)
	assert.Nil(err)
	assert.NotNil(kvstore)

	const (
		port     = 18081
		pass     = "mamay"
		basename = "testing"
	)

	fmt.Println("starting server...")
	srv, err = proxy.NewServer(kvstore, port, basename, pass, false)
	assert.Nil(err)
	assert.NotNil(srv)

	// put some values
	fmt.Println("populating kv...")
	wrtr, _ := kvstore.Writer()
	batch := wrtr.NewBatch()
	for _, pair := range g_testTable {
		batch.Set(pair.k, pair.v)
	}
	err = wrtr.ExecuteBatch(batch)
	assert.Nil(err)

	time.Sleep(1 * time.Second)

	// create client
	fmt.Println("creating client...")
	clientStore, err = proxy.NewClient(port, basename, pass, false)
	assert.Nil(err)

	stat, err := proxy.Stat(clientStore)
	assert.Nil(err)
	assert.Equal(false, stat.TLS)

	closer = func() {
		srv.Close()
		kvstore.Close()
		os.RemoveAll(tmpdir)

		stat := srv.Stat()
		assert.Equal(0, stat.ReadersCount)
		assert.Equal(0, stat.IteratorCount)
	}

	return
}

func TestServerGet(t *testing.T) {
	assert := assertpkg.New(t)
	_, client, closer := set01(assert)
	defer closer()

	rdr, err := client.Reader()
	defer rdr.Close()
	assert.Nil(err)

	for _, pair := range g_testTable {
		bb, err := rdr.Get(pair.k)
		assert.Nil(err)
		assert.Equal(pair.v, bb)
	}

	stat, err := proxy.Stat(client)
	assert.Nil(err)
	assert.Equal(false, stat.TLS)
	assert.Equal(1, stat.ReadersCount)
	assert.Equal(0, stat.IteratorCount)

}

func TestPrefixIter(t *testing.T) {
	assert := assertpkg.New(t)
	_, client, closer := set01(assert)
	defer closer()

	rdr, err := client.Reader()
	assert.Nil(err)

	iter := rdr.PrefixIterator([]byte("some"))
	i := 0
	for ; iter.Valid(); iter.Next() {
		assert.Equal(g_testTable[i].k, iter.Key())
		assert.Equal(g_testTable[i].v, iter.Value())
		i++
	}
	assert.Equal(len(g_testTable), i)

	stat, err := proxy.Stat(client)
	assert.Nil(err)
	assert.Equal(false, stat.TLS)
	assert.Equal(1, stat.ReadersCount)
	assert.Equal(1, stat.IteratorCount)

	iter.Close()
	stat, err = proxy.Stat(client)
	assert.Nil(err)
	assert.Equal(0, stat.IteratorCount)
	assert.Equal(1, stat.ReadersCount)

	rdr.Close()
	stat, err = proxy.Stat(client)
	assert.Nil(err)
	assert.Equal(0, stat.IteratorCount)
	assert.Equal(0, stat.ReadersCount)

}

func TestRangeIter(t *testing.T) {
	assert := assertpkg.New(t)
	_, client, closer := set01(assert)
	defer closer()

	rdr, err := client.Reader()
	assert.Nil(err)

	iter := rdr.RangeIterator([]byte("somek1"), []byte("somek2"))
	i := 1
	for ; iter.Valid(); iter.Next() {
		assert.Equal(g_testTable[i].k, iter.Key())
		assert.Equal(g_testTable[i].v, iter.Value())
		i++
	}
	assert.Equal(2, i)

	stat, err := proxy.Stat(client)
	assert.Nil(err)
	assert.Equal(false, stat.TLS)
	assert.Equal(1, stat.ReadersCount)
	assert.Equal(1, stat.IteratorCount)

	iter.Close()
	stat, err = proxy.Stat(client)
	assert.Nil(err)
	assert.Equal(0, stat.IteratorCount)
	assert.Equal(1, stat.ReadersCount)

	rdr.Close()
	stat, err = proxy.Stat(client)
	assert.Nil(err)
	assert.Equal(0, stat.IteratorCount)
	assert.Equal(0, stat.ReadersCount)

}

func TestMultiGet(t *testing.T) {
	assert := assertpkg.New(t)
	_, client, closer := set01(assert)
	defer closer()

	rdr, err := client.Reader()
	defer rdr.Close()
	assert.Nil(err)

	bbRes, err := rdr.MultiGet([][]byte{
		[]byte("somek0"),
		[]byte("somek3"),
	})
	assert.Nil(err)
	assert.Equal(2, len(bbRes))
	assert.Equal(g_testTable[0].v, bbRes[0])
	assert.Equal(g_testTable[3].v, bbRes[1])

	stat, err := proxy.Stat(client)
	assert.Nil(err)
	assert.Equal(false, stat.TLS)
	assert.Equal(1, stat.ReadersCount)
	assert.Equal(0, stat.IteratorCount)

}
