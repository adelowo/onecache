package memcached

import (
	"flag"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/adelowo/onecache"
	"github.com/bradfitz/gomemcache/memcache"
)

var _ onecache.Store = &MemcachedStore{}

var memcachedStore *MemcachedStore

func TestMain(m *testing.M) {

	memcachedStore = NewMemcachedStore(memcache.New("127.0.0.1:11211"),
		"test:",
	)

	flag.Parse()
	os.Exit(m.Run())
}

var sampleData = []byte("Onecache")

func TestMemcachedStore_Set(t *testing.T) {

	err := memcachedStore.Set("name", sampleData, time.Minute*1)

	if err != nil {
		t.Fatalf(
			`An error occurred while trying to add
			 some data to memcached.. \n %v`, err)
	}
}

func TestMemcachedStore_Get(t *testing.T) {

	val, err := memcachedStore.Get("name")

	if err != nil {
		t.Fatalf(
			`Could not get item with key %s when it in fact
			exists... \n%v`, "name", err)
	}

	if !reflect.DeepEqual(val, sampleData) {
		t.Fatalf(
			`Expected %v \n ..Got %v instead`,
			sampleData, val)
	}
}

func TestMemcachedStore_Delete(t *testing.T) {

	err := memcachedStore.Delete("name")

	if err != nil {
		t.Fatalf(
			`The key %s could not be deleted ... %v`,
			"number", err)
	}
}

func TestMemcachedStore_Flush(t *testing.T) {

	err := memcachedStore.Flush()

	//If nil, we accept the cache is flushed..
	//We od this in best hope on the client library
	//As unlike redis, there isn't a way to get all stored keys in the db
	if err != nil {
		t.Fatalf(
			`An error occurred while clearing the memcached db.. %v`,
			err)
	}

	_, err = memcachedStore.Get("name")

	if err != onecache.ErrCacheMiss {
		t.Fatal("All data should have been cleared from the cache")
	}
}

func TestMemcachedStore_adaptError(t *testing.T) {

	if err := memcachedStore.adaptError(nil); err != nil {
		t.Fatalf(
			`Expected %v.. Got %v`, nil, err)
	}

}

func TestMemcachedStore_DefaultPrefixIsUsedWhenNoneIsSpecified(t *testing.T) {
	memcachedStore = NewMemcachedStore(memcache.New("127.0.0.1:11211"),
		"",
	)

	if !reflect.DeepEqual(memcachedStore.prefix, PREFIX) {
		t.Fatalf(
			`Prefix doen't match.
			\n Expected %s \n.. Got %s`,
			PREFIX, memcachedStore.prefix)
	}
}

func TestMemcachedStore_Has(t *testing.T) {
	memcachedStore = NewMemcachedStore(memcache.New("127.0.0.1:11211"),
		"",
	)

	if ok := memcachedStore.Has("name"); ok {
		t.Fatalf("Key %s is not supposed to exist in the cache", "name")
	}

	memcachedStore.Set("name", []byte("Lanre"), time.Minute*10)

	if ok := memcachedStore.Has("name"); !ok {
		t.Fatalf(`Key %s is supposed to exist in the cache
			as it was recently added`, "name")
	}

}

func TestExtensibility(t *testing.T) {

	_, err := onecache.Get("memcached")

	if err != nil {
		t.Fatalf("An error occured.. %v", err)
	}
}
