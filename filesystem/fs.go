//Package filesystem provides a filesystem cache implementation for onecache
package filesystem

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/adelowo/onecache"
)

const (
	defaultFilePerm          os.FileMode = 0666
	defaultDirectoryFilePerm             = 0755
)


func createDirectory(dir string) error {

	return os.MkdirAll(dir, defaultDirectoryFilePerm)
}


type FSStore struct {
	baseDir string
}

//Returns an initialized Filesystem Cache
//If a non-existent directory is passed, it would be created automatically.
//This function Panics if the directory could not be created
func MustNewFSStore(baseDir string) *FSStore {

	_, err := os.Stat(baseDir)

	if err != nil { //Directory does not exist..Let's create it
		if err := createDirectory(baseDir); err != nil {
			panic(fmt.Errorf("Base directory could not be created : %s", err))
		}
	}

	return &FSStore{baseDir}
}

func (fs *FSStore) Set(key string, data interface{}, expiresAt time.Duration) error {

	return onecache.ErrCacheNotStored
}

func (fs *FSStore) Get(key string) ([]byte, error) {

	return nil, nil
}

func (fs *FSStore) getFileNameFor(key string) string {
	return ""
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
