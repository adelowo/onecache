//Package memory provides a lightweight in memory store for onecache
//Do take a look at other stores
package memory

import (
	"github.com/adelowo/onecache"
	"sync"
	"time"
)

//Represents an inmemory store
type InMemoryStore struct {
	sync.RWMutex
	data map[string][]byte
}

//Returns a new instance of the Inmemory store
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{data: make(map[string][]byte)}
}

func (i *InMemoryStore) Set(key string, data interface{}, expires time.Duration) error {
	i.RLock()

	defer i.RUnlock()

	item := &onecache.Item{ExpiresAt: time.Now().Add(expires), Data: data}

	b, err := item.Bytes()

	if err != nil {
		return err
	}

	i.data[key] = b

	return nil
}

func (i *InMemoryStore) Get(key string) (interface{}, error) {
	i.RLock()
	defer i.RUnlock()

	bytes, ok := i.data[key]

	if !ok {
		return nil, onecache.ErrCacheMiss
	}

	item, err := onecache.BytesToItem(bytes)

	if err != nil {
		return nil, err
	}

	return item.Data, nil
}
