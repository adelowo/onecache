//Package memory provides a lightweight in memory store for onecache
//Do take a look at other stores
package memory

import (
	"sync"
	"time"

	"github.com/adelowo/onecache"
)

// New returns a configured in memory store
func New(opts ...Option) *InMemoryStore {
	i := &InMemoryStore{}

	for _, opt := range opts {
		opt(i)
	}

	if i.keyfn == nil {
		i.keyfn = onecache.DefaultKeyFunc
	}

	if i.data == nil {
		var n int
		if i.bufferSize <= 0 {
			n = 100
			i.bufferSize = 100
		} else {
			n = i.bufferSize
		}

		i.data = make(map[string]*onecache.Item,n)
	}

	return i
}


//Represents an in-memory store
type InMemoryStore struct {
	lock sync.RWMutex
	data map[string]*onecache.Item

	bufferSize int
	keyfn onecache.KeyFunc
}

// NewInMemoryStore returns a new instance of the Inmemory store
// Deprecated... Use New() instead
func NewInMemoryStore() *InMemoryStore {
	return New()
}

func (i *InMemoryStore) Set(key string, data []byte, expires time.Duration) error {
	i.lock.Lock()

	i.data[i.keyfn(key)] = &onecache.Item{
		ExpiresAt: time.Now().Add(expires),
		Data:      copyData(data),
	}

	i.lock.Unlock()
	return nil
}

func (i *InMemoryStore) Get(key string) ([]byte, error) {
	i.lock.RLock()

	item := i.data[i.keyfn(key)]
	if item == nil {
		i.lock.RUnlock()
		return nil, onecache.ErrCacheMiss
	}

	if item.IsExpired() {
		i.lock.RUnlock()
		i.Delete(key)
		return nil, onecache.ErrCacheMiss
	}

	i.lock.RUnlock()
	return copyData(item.Data), nil
}

func (i *InMemoryStore) Delete(key string) error {
	i.lock.Lock()

	_, ok := i.data[i.keyfn(key)]
	if !ok {
		i.lock.Unlock()
		return onecache.ErrCacheMiss
	}

	i.lock.Unlock()
	delete(i.data, i.keyfn(key))
	return nil
}

func (i *InMemoryStore) Flush() error {
	i.lock.Lock()

	i.data = make(map[string]*onecache.Item, i.bufferSize)
	i.lock.Unlock()
	return nil
}

func (i *InMemoryStore) Has(key string) bool {
	i.lock.RLock()

	_, ok := i.data[i.keyfn(key)]
	i.lock.RUnlock()
	return ok
}

func (i *InMemoryStore) GC() {
	i.lock.Lock()

	for k, item := range i.data {
		if item.IsExpired() {
			//No need to spawn a new goroutine since we
			//still have the lock here
			delete(i.data, k)
		}
	}

	i.lock.Unlock()
}

func (i *InMemoryStore) count() int {
	i.lock.Lock()
	n := len(i.data)
	i.lock.Unlock()

	return n
}

func copyData(data []byte) []byte {
	result := make([]byte, len(data))
	copy(result, data)

	return result
}
