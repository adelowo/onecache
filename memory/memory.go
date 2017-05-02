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
	b    onecache.Serializer
	lock sync.RWMutex
	data map[string][]byte
}

//Returns a new instance of the Inmemory store
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{data: make(map[string][]byte), b: onecache.NewCacheSerializer()}
}

func (i *InMemoryStore) Set(key string, data []byte, expires time.Duration) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	item := &onecache.Item{ExpiresAt: time.Now().Add(expires), Data: data}

	b, err := i.b.Serialize(item)

	if err != nil {
		return err
	}

	i.data[key] = b

	return nil
}

func (i *InMemoryStore) Get(key string) ([]byte, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	bytes, ok := i.data[key]

	if !ok {
		return nil, onecache.ErrCacheMiss
	}

	item := new(onecache.Item)

	err := i.b.DeSerialize(bytes, item)

	if err != nil {
		return nil, err
	}

	if item.IsExpired() {
		go i.Delete(key) //Prevent a deadlock since the mutex is still locked here
		return nil, onecache.ErrCacheMiss
	}

	return item.Data, nil
}

func (i *InMemoryStore) Delete(key string) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	_, ok := i.data[key]

	if !ok {
		return onecache.ErrCacheMiss
	}

	delete(i.data, key)

	return nil
}

func (i *InMemoryStore) Flush() error {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.data = make(map[string][]byte)

	return nil
}
