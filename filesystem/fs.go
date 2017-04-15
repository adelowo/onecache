//Package filesystem provides a filesystem cache implementation for onecache
package filesystem

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"os"
	"time"

	"github.com/adelowo/onecache"
)

const (
	defaultFilePerm          os.FileMode = 0666
	defaultDirectoryFilePerm             = 0755
)

var (
	//Onecache's error that defines an error with json marshalling
	ErrCacheCouldNotBeSerialized = errors.New("Data could not be serialized")
)

func createDirectory(dir string) error {

	return os.MkdirAll(dir, defaultDirectoryFilePerm)
}

//identides a cached piece of data
type item struct {
	ExpiresAt time.Time   `json:"expires_at"`
	Data      interface{} `json:"data"`
}

type FSCache struct {
	baseDir string
	hasher  hash.Hash
}

//Returns an initialized Filesystem Cache
//If a non-existent directory is passed, it would be created automatically.
//This function Panics if the directory could not be created
func MustNewFSCache(baseDir string) *FSCache {

	_, err := os.Stat(baseDir)

	if err != nil { //Directory does not exist..Let's create it
		err = createDirectory(baseDir)

		if err != nil {
			panic(fmt.Errorf("Base directory could not be created : %s", err))
		}
	}

	return &FSCache{baseDir, md5.New()}
}

func (fs *FSCache) Set(key string, data interface{}, expiresAt time.Duration) error {

	i := item{time.Now().Add(expiresAt), data}

	b, err := json.Marshal(i)

	if err != nil {
		return ErrCacheCouldNotBeSerialized
	}

	if err = ioutil.WriteFile(fs.getFileNameFor(key), b, defaultFilePerm); err == nil {
		return nil
	}

	return onecache.ErrCacheNotStored
}

func (fs *FSCache) Get(key string) ([]byte, error) {
	panic("Not implemented")
}

func (fs *FSCache) getFileNameFor(key string) string {
	hashedFilePath := fmt.Sprintf("%x", string(fs.hasher.Sum([]byte(key))))

	return fs.baseDir + string(os.PathSeparator) + hashedFilePath
}

func toBytes(val interface{}) ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	err := enc.Encode(val)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
