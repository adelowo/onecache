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

func TestMemcachedStore_Set(t *testing.T) {

	err := memcachedStore.Set("name", "memcached", time.Minute*1)

	if err != nil {
		t.Fatalf("An error occurred while trying to add some data to memcached.. \n %v",
			err)
	}
}

func TestMemcachedStore_Get(t *testing.T) {

	val, err := memcachedStore.Get("name")

	if err != nil {
		t.Fatalf("Could not get item with key %s when it in fact exists... \n%v",
			"name", err)
	}

	if !reflect.DeepEqual(val, "memcached") {
		t.Fatalf("Expected %v \n ..Got %v instead", "memcached", val)
	}

	memcachedStore.Set("number", 42, onecache.EXPIRES_DEFAULT)

	val, err = memcachedStore.Get("number")

	if err != nil {
		t.Fatalf("Could not get item with key %s when it in fact exists... \n%v",
			"number", err)
	}

	if !reflect.DeepEqual(val, 42) {
		t.Fatalf("Expected %v \n ..Got %v instead", 42, val)
	}
}

func TestMemcachedStore_Delete(t *testing.T) {

	err := memcachedStore.Delete("number")

	if err != nil {
		t.Fatalf("The key %s could not be deleted ... %v", "number", err)
	}
}

func TestMemcachedStore_Flush(t *testing.T) {

	err := memcachedStore.Flush()

	//If nil, we accept the cache is flushed..
	//We od this in best hope on the client library
	//As unlike redis, there isn't a way to get all stored keys in the db
	if err != nil {
		t.Fatalf("An error occurred while clearing the memcached db.. %v", err)
	}

	_, err = memcachedStore.Get("name")

	if err != onecache.ErrCacheMiss {
		t.Fatal("All data should have been cleared from the cache")
	}
}

func TestMemcachedStore_DefaultPrefixIsUsedWhenNoneIsSpecified(t *testing.T) {
	memcachedStore = NewMemcachedStore(memcache.New("127.0.0.1:11211"),
		"",
	)

	if !reflect.DeepEqual(memcachedStore.prefix, PREFIX) {
		t.Fatalf("Prefix doen't match. \n Expected %s \n.. Got %s", PREFIX, memcachedStore.prefix)
	}
}
