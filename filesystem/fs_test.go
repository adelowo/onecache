package filesystem

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/adelowo/onecache"
)

var _ onecache.Store = MustNewFSStore("./")

var _ onecache.GarbageCollector = MustNewFSStore("./")

var fileCache *FSStore

func TestMain(m *testing.M) {

	fileCache = MustNewFSStore("./../cache")

	flag.Parse()
	os.Exit(m.Run())
}

func TestMustNewFSStore(t *testing.T) {

	defer func() {
		recover()
	}()

	_ = MustNewFSStore("/hh")
}

var sampleData = []byte("Lanre")

func TestFSStore_Set(t *testing.T) {

	err := fileCache.Set("name", sampleData, time.Minute*2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFSStore_Get(t *testing.T) {

	val, err := fileCache.Get("name")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(val, sampleData) {
		t.Fatalf(
			`Values are not equal.. Expected %v \n
			Got %v`, sampleData, val)
	}
}

func TestFSStore_GetUnknownKey(t *testing.T) {
	val, err := fileCache.Get("unknown")

	if err == nil {
		t.Fatal("Expected an error for a file that doesn't exist on the filesystem")
	}

	if val != nil {
		t.Fatalf("Expected a nil item to be return ... Got %v instead", val)
	}

	if err != onecache.ErrCacheMiss {
		t.Fatalf("unexpected error found...Expected %v , got %v", onecache.ErrCacheMiss, err)
	}
}

func TestFSStore_GarbageCollection(t *testing.T) {

	err := fileCache.Set("xyz", []byte("Elon Musk"), onecache.EXPIRES_DEFAULT)

	if err != nil {
		t.Fatalf("An error occurred... %v", err)
	}

	data, err := fileCache.Get("xyz")

	if err != onecache.ErrCacheMiss {
		t.Fatal("Cached data is supposed to be expired")
	}

	if data != nil {
		t.Fatal("Garbage collected item is supposed to be empty")
	}
}

func TestFSStore_Flush(t *testing.T) {
	if err := fileCache.Flush(); err != nil {
		t.Fatalf("The cache directory, %s could not be flushed... %v", fileCache.baseDir, err)
	}
}

func TestFilePathForKey(t *testing.T) {

	key := "page_hits"

	b := md5.Sum([]byte(key))
	s := hex.EncodeToString(b[:])

	path := filepath.Join(fileCache.baseDir, s[0:2], s[2:4], s[4:6], s)

	if x := fileCache.filePathFor("page_hits"); path != x && FilePathKeyFunc(key) != path {
		t.Fatalf("Path differs.. Expected %s. Got %s instead", path, x)
	}
}

type mockSerializer struct {
}

func (b *mockSerializer) Serialize(i interface{}) ([]byte, error) {
	return nil, errors.New("Yup an error occurred")
}

func (b *mockSerializer) DeSerialize(data []byte, i interface{}) error {
	return errors.New("Yet another error")
}

func TestFSStore_GetFailsBecauseOfBytesMarshalling(t *testing.T) {

	fileCache.Set("test", []byte("test"), time.Second*1)

	fs, err := New(BaseDirectory("./cache"), Serializer(&mockSerializer{}))
	if err != nil {
		t.Fatal(err)
	}

	_, err = fs.Get("test")
	if err == nil {
		t.Fatalf(
			`Expected a cache miss.. Got %v`, err)
	}
}

func TestFSStore_SetFailsBecauseOfBytesMarshalling(t *testing.T) {

	fs, err := New(BaseDirectory("./cache"), Serializer(&mockSerializer{}))
	if err != nil {
		t.Fatal(err)
	}

	defer fs.Flush()

	err = fs.Set("test", []byte("test"), time.Nanosecond*4)
	if err == nil {
		t.Fatalf(
			`Expected an error from bytes marshalling.. Got %v`, err)
	}
}

func TestFSStore_Delete(t *testing.T) {
	if err := fileCache.Delete("name"); err != nil {
		t.Fatalf("Could not delete the cached data... %v", err)
	}
}

func TestFSStore_GC(t *testing.T) {
	store := MustNewFSStore("./../cache")
	defer store.Flush()

	tableTests := []struct {
		key, value string
		expires    time.Duration
	}{
		{"name", "Onecache", time.Microsecond},
		{"number", "Fourty two", time.Microsecond},
		{"x", "yz", time.Microsecond},
	}

	for _, v := range tableTests {
		store.Set(v.key, []byte(v.value), v.expires)
	}

	ticker := time.NewTicker(time.Millisecond)
	<-ticker.C

	store.GC()

	var filePath string

	for _, v := range tableTests {

		filePath = store.filePathFor(v.key)

		if _, err := os.Stat(filePath); err == nil {
			t.Fatal(
				`File exists when it isn't supposed to since there was
				a garbage collection`)
		}
	}
}

func TestFSStore_Has(t *testing.T) {
	store := MustNewFSStore("./../cache")

	defer store.Flush()

	if ok := store.Has("name"); ok {
		t.Fatalf("Key %s is not supposed to exist in the cache", "name")
	}

	store.Set("name", []byte("Lanre"), time.Hour*10)

	if ok := store.Has("name"); !ok {
		t.Fatalf(`Expected store to have an item with key %s
			since that key was persisted secs ago`, "name")
	}
}

func BenchmarkFSStore_Get(b *testing.B) {

	store := MustNewFSStore("./../cache")

	key := "life"
	answer := []byte("42")

	err := store.Set(key, answer, time.Minute*10)
	if err != nil {
		b.Fatalf("an error occurred while setting key in cache store... %v", err)
	}

	b.ResetTimer()

	equals := func(a, b []byte) bool {
		return bytes.Equal(a, b)
	}

	for i := 0; i <= b.N; i++ {
		buf, err := store.Get(key)
		if err != nil {
			b.Fatalf("an error happened while reading from cache.... %v", err)
		}

		if !equals(answer, buf) {
			b.Fatalf(`Data does not match... Expected %v \n
			Got %v`, answer, buf)
		}
	}
}
