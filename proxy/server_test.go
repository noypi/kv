package proxy_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/noypi/kv/leveldb"
	"github.com/noypi/kv/proxy"
	assertpkg "github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	assert := assertpkg.New(t)
	tmpdir := os.TempDir() + "/testing-proxyserver"
	os.RemoveAll(tmpdir)
	defer os.RemoveAll(tmpdir)

	kvstore, err := leveldb.GetDefault(tmpdir)
	assert.Nil(err)
	assert.NotNil(kvstore)

	const (
		port = 18081
		pass = "mamay"
	)

	fmt.Println("starting server...")
	srv, err := proxy.NewServer(kvstore, port, pass, false)
	assert.Nil(err)
	assert.NotNil(srv)
	defer srv.Close()

	// put some values
	fmt.Println("populating kv...")
	wrtr, _ := kvstore.Writer()
	batch := wrtr.NewBatch()
	batch.Set([]byte("somek0"), []byte("someval0"))
	batch.Set([]byte("somek1"), []byte("someval1"))
	batch.Set([]byte("somek2"), []byte("someval2"))
	err = wrtr.ExecuteBatch(batch)
	assert.Nil(err)

	time.Sleep(1 * time.Second)

	// create client
	fmt.Println("creating client...")
	client, err := proxy.NewClient(port, pass, false)
	assert.Nil(err)

	rdr, err := client.Reader()
	assert.Nil(err)

	fmt.Println("getting somek0...")
	bb, err := rdr.Get([]byte("somek0"))
	assert.Nil(err)

	assert.Equal([]byte("someval0"), bb)

}
