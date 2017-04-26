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
