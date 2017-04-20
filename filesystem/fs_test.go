package filesystem

import (
	"flag"
	"os"
	"testing"
	"time"

	"crypto/md5"
	"encoding/hex"
	"path/filepath"

	"github.com/adelowo/onecache"
)

var _ onecache.CacheStore = MustNewFSStore("./")

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

func TestFSStore_Set(t *testing.T) {

	err := fileCache.Set("name", "Lanre", time.Minute*2)

	if err != nil {
		t.Fatal(err)
	}
}

func TestFSStore_Get(t *testing.T) {

	val, err := fileCache.Get("name")

	if err != nil {
		t.Fatal(err)
	}

	data, ok := val.(string)

	if !ok {
		t.Fatal("Cached data should return a string")
	}

	if data != "Lanre" {
		t.Fatal("OOPS")
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

	err := fileCache.Set("xyz", "Elon Musk", onecache.EXPIRES_DEFAULT)

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

func TestFSStore_Delete(t *testing.T) {

	if err := fileCache.Delete("name"); err != nil {
		t.Fatalf("Could not delete the cached data... %v", err)
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
