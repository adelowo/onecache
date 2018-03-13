package memory

import (
	"bytes"
	"flag"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/adelowo/onecache"
)

var _ onecache.Store = &InMemoryStore{}

var _ onecache.GarbageCollector = &InMemoryStore{}

var memoryStore *InMemoryStore

func TestMain(t *testing.M) {

	memoryStore = NewInMemoryStore()

	flag.Parse()

	os.Exit(t.Run())
}

var sampleData = []byte("Lanre")

func TestInMemoryStore_Set(t *testing.T) {

	err := memoryStore.Set("name", sampleData, time.Minute*1)

	if err != nil {
		t.Fatalf("Data could not be stored in the inmemory store.. \n%v", err)
	}
}

func TestInMemoryStore_StoresCopy(t *testing.T) {
	data := []byte("abcdef")
	err := memoryStore.Set("key", data, time.Minute*1)
	if err != nil {
		t.Fatalf("Data could not be stored in the inmemory store.. \n%v", err)
	}

	// modify the set input
	data[0] = 'z'
	data[1] = 'z'

	val, err := memoryStore.Get("key")
	if err != nil {
		t.Fatalf("Key %s should exist in the store... \n %v", "name", err)
	}

	if !bytes.Equal(val, []byte("abcdef")) {
		t.Fatalf("Data was not as expected: %v", val)
	}

	// modify the get output
	val[0] = 'z'
	val[1] = 'z'

	val, err = memoryStore.Get("key")
	if err != nil {
		t.Fatalf("Key %s should exist in the store... \n %v", "name", err)
	}

	if !bytes.Equal(val, []byte("abcdef")) {
		t.Fatalf("Data was not as expected: %v", val)
	}
}

func TestInMemoryStore_Get(t *testing.T) {

	val, err := memoryStore.Get("name")

	if err != nil {
		t.Fatalf("Key %s should exist in the store... \n %v", "name", err)
	}

	if !reflect.DeepEqual(sampleData, val) {
		t.Fatalf(
			`Data returned from the store does not match what was returned..
			\n.Expected %v \n.. Got %v instead`,
			sampleData,
			val)
	}
}

func TestInMemoryStore_Get_GarabageCollection(t *testing.T) {
	memoryStore.Set("expiredItem", []byte("I just set this"), time.Nanosecond*1)

	val, err := memoryStore.Get("expiredItem")

	if val != nil {
		t.Fatalf(
			`Expected data to have a nil value.. Got %v instead`,
			val)
	}

	if err != onecache.ErrCacheMiss {
		t.Fatalf(
			`Exoected error to be a cache miss..
			\n Expected %v \n Got %v instead`,
			onecache.ErrCacheMiss,
			err)
	}
}

func TestInMemoryStore_GetUnknownKey(t *testing.T) {

	val, err := memoryStore.Get("unknown")

	if err != onecache.ErrCacheMiss {
		t.Fatalf(
			`Expeted to get a cache miss error.. \n
			Got %v instead`,
			err)
	}

	if val != nil {
		t.Fatalf(
			`Expected %v to be a nil value`,
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

	_, err = memoryStore.Get("name")

	if err != onecache.ErrCacheMiss {
		t.Fatalf(
			"Expected an error of %v \n Got %v",
			onecache.ErrCacheMiss,
			err)
	}

	//Test no-op on non-existent key

	err = memoryStore.Delete("unknown")

	if err != onecache.ErrCacheMiss {
		t.Fatalf(
			`Error should be a missed cache.
			 \n. Expected %v.\n Got %v`,
			onecache.ErrCacheMiss,
			err)
	}
}

func TestInMemoryStore_Flush(t *testing.T) {
	//Add some more data

	//Number of items that should be left in the store after flushing
	expectedNumberOfItems := 0

	expiresAt := time.Minute * 10

	memoryStore.Set("name", []byte("onecache"), expiresAt)
	memoryStore.Set("me", []byte("you"), expiresAt)
	memoryStore.Set("something", []byte("else"), expiresAt)

	err := memoryStore.Flush()

	if err != nil {
		t.Fatalf("An error occurred while the store was being flushed... %v", err)
	}

	if x := memoryStore.count(); x != expectedNumberOfItems {
		t.Fatalf(
			"Store was not flushed..\n Expected %d.. Got %d ",
			expectedNumberOfItems, x)
	}
}

func TestInMemoryStore_GC(t *testing.T) {

	store := NewInMemoryStore()

	tableTests := []struct {
		key, value string
		expires    time.Duration
	}{
		{"name", "Onecache", time.Microsecond},
		{"number", "Fourty two", time.Microsecond},
		{"x", "yz", time.Microsecond},
	}

	for _, v := range tableTests {
		if err := store.Set(v.key, []byte(v.value), v.expires); err != nil {
			t.Fatalf("an error occurred while trying to write to store.. %v", err)
		}
	}

	ticker := time.NewTicker(time.Millisecond)

	<-ticker.C
	store.GC()
	//GC should wipe everything off since they are well past
	//their expiration time
	expectedNumberOfItemsInStore := 0

	if x := store.count(); x != expectedNumberOfItemsInStore {
		t.Fatalf(
			`Expected %d items in the store. %d found`,
			expectedNumberOfItemsInStore, x)
	}

	ticker.Stop()
}

func TestInMemoryStore_Has(t *testing.T) {

	store := NewInMemoryStore()

	if ok := store.Has("name"); ok {
		t.Fatalf("Key %s does not exist", "name")
	}

	store.Set("name", []byte("Lanre"), time.Second*2)

	if ok := store.Has("name"); !ok {
		t.Fatalf("Key %s was set and is supposed to exist", "name")
	}
}
