//Package memory provides a lightweight in memory store for onecache
//Do take a look at other stores
package memory

import (
	"sync"
	"time"

	"github.com/adelowo/onecache"
)

//Represents an inmemory store
type InMemoryStore struct {
	lock sync.RWMutex
	data map[string][]byte
}

//Returns a new instance of the Inmemory store
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{data: make(map[string][]byte)}
}

func (i *InMemoryStore) Set(key string, data interface{}, expires time.Duration) error {
	i.lock.RLock()

	defer i.lock.RUnlock()

	item := &onecache.Item{ExpiresAt: time.Now().Add(expires), Data: data}

	b, err := item.Bytes()

	if err != nil {
		return err
	}

	i.data[key] = b

	return nil
}

func (i *InMemoryStore) Get(key string) (interface{}, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	bytes, ok := i.data[key]

	if !ok {
		return nil, onecache.ErrCacheMiss
	}

	item, err := onecache.BytesToItem(bytes)

	if item.IsExpired() {
		i.Delete(key)
		return nil, onecache.ErrCacheMiss
	}

	if err != nil {
		return nil, err
	}

	return item.Data, nil
}

func (i *InMemoryStore) Delete(key string) error {
	i.lock.RLock()
	defer i.lock.RUnlock()

	_, ok := i.data[key]

	if !ok {
		return onecache.ErrCacheMiss
	}

	delete(i.data, key)

	return nil
}

func (i *InMemoryStore) Flush() error {
	i.lock.RLock()
	defer i.lock.RUnlock()

	i.data = make(map[string][]byte)

	return nil
}

func (i *InMemoryStore) Increment(key string, steps int) error {

	i.lock.RLock()
	defer i.lock.RUnlock()

	if !i.has(key) {
		return onecache.ErrCacheMiss
	}

	bytes := i.data[key]

	item, err := onecache.BytesToItem(bytes)

	if err != nil {
		return err
	}

	item.Data, err = onecache.Increment(item.Data, steps)

	if err != nil {
		return err
	}

	b, err := item.Bytes()

	if err != nil {
		return err
	}

	i.data[key] = b

	return nil

}

func (i *InMemoryStore) has(key string) bool {

	_, ok := i.data[key]

	return ok
}
