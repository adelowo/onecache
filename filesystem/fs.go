//Package filesystem provides a filesystem cache implementation for onecache
package filesystem

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"errors"
	"strings"

	"github.com/adelowo/onecache"
)

const (
	defaultFilePerm          os.FileMode = 0666
	defaultDirectoryFilePerm             = 0755
)

func FilePathKeyFunc(s string) string {

	hashSum := md5.Sum([]byte(s))
	hashSumAsString := hex.EncodeToString(hashSum[:])

	return filepath.Join(string(hashSumAsString[0:2]),
		string(hashSumAsString[2:4]),
		string(hashSumAsString[4:6]), hashSumAsString)
}

func createDirectory(dir string) error {
	return os.MkdirAll(dir, defaultDirectoryFilePerm)
}

type FSStore struct {
	baseDir string
	b       onecache.Serializer
	keyFn   onecache.KeyFunc
}

func MustNewFSStore(baseDir string) *FSStore {

	_, err := os.Stat(baseDir)

	if err != nil { //Directory does not exist..Let's create it
		if err := createDirectory(baseDir); err != nil {
			panic(fmt.Errorf("Base directory could not be created : %s", err))
		}
	}

	fs, err := New(BaseDirectory(baseDir), CacheKeyGenerator(FilePathKeyFunc))
	if err != nil {
		panic(err)
	}

	return fs
}

func New(opts ...Option) (*FSStore, error) {

	store := &FSStore{}

	for _, opt := range opts {
		opt(store)
	}

	if store.b == nil {
		store.b = onecache.NewCacheSerializer()
	}

	if len(strings.TrimSpace(store.baseDir)) == 0 {
		return nil, errors.New("onecache : base directory not provided")
	}

	if store.keyFn == nil {
		store.keyFn = FilePathKeyFunc
	}

	return store, nil
}

func (fs *FSStore) Set(key string, data []byte, expiresAt time.Duration) error {

	path := fs.filePathFor(key)

	if err := createDirectory(filepath.Dir(path)); err != nil {
		return err
	}

	i := &onecache.Item{ExpiresAt: time.Now().Add(expiresAt), Data: data}

	b, err := fs.b.Serialize(i)
	if err != nil {
		return err
	}

	return writeFile(path, b)
}

func (fs *FSStore) Get(key string) ([]byte, error) {

	var b = new(bytes.Buffer)

	f, err := os.OpenFile(fs.filePathFor(key), os.O_RDONLY, 0644)
	if err != nil {
		pe, ok := err.(*os.PathError)
		if !ok {
			return nil, err
		}

		if pe.Err == syscall.ENOENT && pe.Op == "open" {
			return nil, onecache.ErrCacheMiss
		}

		return nil, err
	}

	if _, err := io.Copy(b, f); err != nil {
		f.Close()
		return nil, err
	}

	f.Close()

	i := new(onecache.Item)

	if err := fs.b.DeSerialize(b.Bytes(), i); err != nil {
		return nil, err
	}

	if i.IsExpired() {
		fs.Delete(key)
		return nil, onecache.ErrCacheMiss
	}

	return i.Data, nil
}

func (fs *FSStore) Delete(key string) error {
	return os.RemoveAll(fs.filePathFor(key))
}

func (fs *FSStore) Flush() error {
	return os.RemoveAll(fs.baseDir)
}

func (fs *FSStore) GC() {

	filepath.Walk(
		fs.baseDir,
		func(path string, finfo os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if finfo.IsDir() {
				return nil
			}

			currentItem := new(onecache.Item)

			byt, err := ioutil.ReadFile(path)

			if err != nil {
				return err
			}

			if err = fs.b.DeSerialize(byt, currentItem); err != nil {
				return err
			}

			if currentItem.IsExpired() {
				if err := os.Remove(path); !os.IsExist(err) {
					return err
				}
			}

			return nil
		})
}

func (fs *FSStore) Has(key string) bool {
	_, err := os.Stat(fs.filePathFor(key))
	return !os.IsNotExist(err)
}

func (fs *FSStore) filePathFor(key string) string {
	return filepath.Join(fs.baseDir, fs.keyFn(key))
}

func writeFile(path string, b []byte) error {
	return ioutil.WriteFile(path, b, defaultFilePerm)
}
