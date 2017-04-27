package memory

import (
	"flag"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/adelowo/onecache"
)

var _ onecache.Store = &InMemoryStore{}

var memoryStore *InMemoryStore

func TestMain(t *testing.M) {

	memoryStore = NewInMemoryStore()

	flag.Parse()

	os.Exit(t.Run())
}

func TestInMemoryStore_Set(t *testing.T) {

	err := memoryStore.Set("name", "Lanre", time.Minute*1)

	if err != nil {
		t.Fatalf("Data could not be stored in the inmemory store.. \n%v", err)
	}
}

func TestInMemoryStore_Get(t *testing.T) {

	val, err := memoryStore.Get("name")

	if err != nil {
		t.Fatalf("Key %s should exist in the store... \n %v", "name", err)
	}

	if !reflect.DeepEqual("Lanre", val) {
		t.Fatalf(
			`Data returned from the store does not match what was returned..
			\n.Expected %v \n.. Got %v instead`,
			"Lanre",
			val)
	}
}

func TestInMemoryStore_Delete(t *testing.T) {

	err := memoryStore.Delete("name")

	if err != nil {
		t.Fatalf(
			"An error occurred while trying to delete the data from the store... %v",
			err)
	}
}

func TestInMemoryStore_Flush(t *testing.T) {
	//Add some more data

	//Number of items that should be left in the store after flushing
	expectedNumberOfItems := 0

	expiresAt := time.Minute * 10

	memoryStore.Set("name", "onecache", expiresAt)
	memoryStore.Set("me", "you", expiresAt)
	memoryStore.Set("something", "else", expiresAt)

	err := memoryStore.Flush()

	if err != nil {
		t.Fatalf("An error occurred while the store was being flushed... %v", err)
	}

	if x := len(memoryStore.data); x != expectedNumberOfItems {
		t.Fatalf(
			"Store was not flushed..\n Expected %d.. Got %d ",
			expectedNumberOfItems, x)
	}
}

func TestInMemoryStore_GetUnknownKey(t *testing.T) {

	_, err := memoryStore.Get("unknownKey")

	if err != onecache.ErrCacheMiss {
		t.Fatal("Item does not exist in the store, yet it was found")
	}
}

func TestInMemoryStore_Get_GarbageCollectsExpiredKeys(t *testing.T) {
	memoryStore.Set("animal", "Gopher", time.Nanosecond)

	_, err := memoryStore.Get("animal")

	if err != onecache.ErrCacheMiss {
		t.Fatalf("Key is expired and should be garbage collected.. %v", err)
	}
}

func TestInMemoryStore_Delete_UnknownKey(t *testing.T) {

	err := memoryStore.Delete("animal")

	if err != onecache.ErrCacheMiss {
		t.Fatalf(
			"An unknown key should return a cache miss error.. Received %v",
			err)
	}
}
