package filesystem

import (
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

	if x := fileCache.filePathFor("page_hits"); path != x {
		t.Fatalf("Path differs.. Expected %s. Got %s instead", path, x)
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

func TestFSStore_GetFailsBecauseOfBytesMarshalling(t *testing.T) {

	fileCache.Set("test", []byte("test"), time.Second*1)

	fs := &FSStore{"./../cache", &bytesItemMarshallerMock{}}

	_, err := fs.Get("test")

	if err == nil {
		t.Fatalf(
			`Expected a cache miss.. Got %v`, err)
	}

}

func TestFSStore_SetFailsBecauseOfBytesMarshalling(t *testing.T) {

	fs := &FSStore{"./../cache", &bytesItemMarshallerMock{}}

	err := fs.Set("test", []byte("test"), time.Nanosecond*4)

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

func TestFSStore_SetFailsWhenMakingUseOfAnUnwriteableDirectory(t *testing.T) {
	fileCache.baseDir = "/" //change directory to the OS root

	if err := fileCache.Set("test", []byte("test"), time.Microsecond*4); err == nil {
		t.Fatal(
			`An error was supposed to occur because the root directory isn't writeable`)
	}
}
