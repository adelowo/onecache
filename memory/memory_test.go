package memory

import (
	"errors"
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

var sampleData = []byte("Lanre")

func TestInMemoryStore_Set(t *testing.T) {

	err := memoryStore.Set("name", sampleData, time.Minute*1)

	if err != nil {
		t.Fatalf("Data could not be stored in the inmemory store.. \n%v", err)
	}
}

type bytesItemMarshallerMock struct {
}

func (b *bytesItemMarshallerMock) MarshalBytes(i *onecache.Item) ([]byte, error) {
	return nil, errors.New("Yup an error occurred")
}

func (b *bytesItemMarshallerMock) UnMarshallBytes(data []byte) (*onecache.Item, error) {
	return nil, errors.New("Yet another error")
}

func TestInMemoryStore_SetErrorOccursWhenMarshallingItemToByte(t *testing.T) {

	m := &InMemoryStore{data: make(map[string][]byte, 100), b: &bytesItemMarshallerMock{}}

	err := m.Set("n", []byte("ERROR"), time.Second*2)

	if err == nil {
		t.Fatalf(
			`Error should be nil as the item could not be marshalled into
			bytes.. Got %v`,
			err)
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

func TestInMemoryStore_Get_ErrorOccurrsWhileUnmarshallingToBytes(t *testing.T) {

	//Since the mock marshaller fails everything,
	//we create an inmemory map contianing sample data Get can read from

	type d map[string][]byte

	f := make(d, 100)

	f["name"] = []byte("Onecache")

	m := &InMemoryStore{data: f, b: &bytesItemMarshallerMock{}}

	val, err := m.Get("name")

	if err == nil {
		t.Fatalf(
			`An error is supposed to occur if bytes marshalling fails.. Got %v`,
			err)
	}

	if val != nil {
		t.Fatalf("Value is supposed to be nil.. Got %v instead", val)
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
}

//func TestInMemoryStore_Flush(t *testing.T) {
//	//Add some more data
//
//	//Number of items that should be left in the store after flushing
//	expectedNumberOfItems := 0
//
//	expiresAt := time.Minute * 10
//
//	memoryStore.Set("name", "onecache", expiresAt)
//	memoryStore.Set("me", "you", expiresAt)
//	memoryStore.Set("something", "else", expiresAt)
//
//	err := memoryStore.Flush()
//
//	if err != nil {
//		t.Fatalf("An error occurred while the store was being flushed... %v", err)
//	}
//
//	if x := len(memoryStore.data); x != expectedNumberOfItems {
//		t.Fatalf(
//			"Store was not flushed..\n Expected %d.. Got %d ",
//			expectedNumberOfItems, x)
//	}
//}
//
//func TestInMemoryStore_GetUnknownKey(t *testing.T) {
//
//	_, err := memoryStore.Get("unknownKey")
//
//	if err != onecache.ErrCacheMiss {
//		t.Fatal("Item does not exist in the store, yet it was found")
//	}
//}
//
//func TestInMemoryStore_Get_GarbageCollectsExpiredKeys(t *testing.T) {
//	memoryStore.Set("animal", "Gopher", time.Nanosecond)
//
//	_, err := memoryStore.Get("animal")
//
//	if err != onecache.ErrCacheMiss {
//		t.Fatalf("Key is expired and should be garbage collected.. %v", err)
//	}
//}
//
//func TestInMemoryStore_Delete_UnknownKey(t *testing.T) {
//
//	err := memoryStore.Delete("animal")
//
//	if err != onecache.ErrCacheMiss {
//		t.Fatalf(
//			"An unknown key should return a cache miss error.. Received %v",
//			err)
//	}
//}
//
//func TestInMemoryStore_Increment(t *testing.T) {
//
//	expected := int32(52)
//
//	memoryStore.Set("number", int32(42), time.Second*10)
//
//	err := memoryStore.Increment("number", 10)
//
//	if err != nil {
//		t.Fatalf("An error occured while trying to increment the data.. %v", err)
//	}
//
//	val, _ := memoryStore.Get("number")
//
//	if !reflect.DeepEqual(expected, val) {
//		t.Fatalf(
//			"Incrementing cache data failed..\n Expected %d, got %d instead",
//			expected, val)
//	}
//}
//
//func TestInMemoryStore_Decrement(t *testing.T) {
//
//	expected := int32(42)
//
//	memoryStore.Set("number", int32(52), time.Second*10)
//
//	err := memoryStore.Decrement("number", 10)
//
//	if err != nil {
//		t.Fatalf("An error occured while trying to increment the data.. %v", err)
//	}
//
//	val, _ := memoryStore.Get("number")
//
//	if !reflect.DeepEqual(expected, val) {
//		t.Fatalf(
//			"Incrementing cache data failed..\n Expected %d, got %d instead",
//			expected, val)
//	}
//}
