//Package memory provides a lightweight in memory store for onecache
//Do take a look at other stores
package memory

import (
	"sync"
	"time"

	"github.com/adelowo/onecache"
)

func init() {
	onecache.Extend("memory", func() onecache.Store {
		return &InMemoryStore{
			data: make(map[string]*onecache.Item),
		}
	})
}

//Represents an inmemory store
type InMemoryStore struct {
	lock sync.RWMutex
	data map[string]*onecache.Item
}

//Returns a new instance of the Inmemory store
func NewInMemoryStore(gcInterval time.Duration) *InMemoryStore {
	i := &InMemoryStore{
		data: make(map[string]*onecache.Item),
	}

	go i.GC(gcInterval)
	return i
}

func (i *InMemoryStore) Set(key string, data []byte, expires time.Duration) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.data[key] = &onecache.Item{
		ExpiresAt: time.Now().Add(expires),
		Data:      copyData(data),
	}

	return nil
}

func (i *InMemoryStore) Get(key string) ([]byte, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	item := i.data[key]
	if item == nil {
		return nil, onecache.ErrCacheMiss
	}

	if item.IsExpired() {
		go i.Delete(key) //Prevent a deadlock since the mutex is still locked here
		return nil, onecache.ErrCacheMiss
	}

	return copyData(item.Data), nil
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

	i.data = make(map[string]*onecache.Item)
	return nil
}

func (i *InMemoryStore) Has(key string) bool {
	i.lock.Lock()
	defer i.lock.Unlock()

	_, ok := i.data[key]
	return ok
}

func (i *InMemoryStore) GC(gcInterval time.Duration) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if gcInterval <= (time.Second * 1) {
		return
	}

	for k, item := range i.data {
		if item.IsExpired() {
			//No need to spawn a new goroutine since we
			//still have the lock here
			delete(i.data, k)
		}
	}

	time.AfterFunc(gcInterval, func() {
		i.GC(gcInterval)
	})
}

func (i *InMemoryStore) count() int {
	i.lock.Lock()
	defer i.lock.Unlock()

	return len(i.data)
}

func copyData(data []byte) []byte {
	result := make([]byte, len(data))
	copy(result, data)

	return result
}
