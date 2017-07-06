//Package memcached is a cache implementation for onecache which uses memcached
package memcached

import (
	"time"

	"github.com/adelowo/onecache"
	"github.com/bradfitz/gomemcache/memcache"
)

func init() {

	//Make use of NewMemcachedStore() for customization of the server list
	//and others
	onecache.Extend("memcached", func() onecache.Store {
		return NewMemcachedStore(memcache.New("127.0.0.1:11211"),
			"onecache:",
		)
	})
}

type MemcachedStore struct {
	client *memcache.Client
	prefix string
}

//PREFIX prevents collision with other items stored in the db
const PREFIX = "onecache:"

//Returns a new instance of the memached store.
//If prefix is an empty string, it defaults to the package's prefix constant
func NewMemcachedStore(c *memcache.Client, prefix string) *MemcachedStore {

	var p string

	if prefix == "" {
		p = PREFIX
	} else {
		p = prefix
	}

	return &MemcachedStore{client: c, prefix: p}
}

func (m *MemcachedStore) key(k string) string {
	return m.prefix + k
}

func (m *MemcachedStore) Set(k string, data []byte, expires time.Duration) error {

	item := &memcache.Item{
		Key:        m.key(k),
		Value:      data,
		Expiration: int32(expires / time.Second),
	}

	return m.client.Set(item)
}

func (m *MemcachedStore) Get(k string) ([]byte, error) {

	val, err := m.client.Get(m.key(k))

	if err != nil {
		return nil, m.adaptError(err)
	}

	return val.Value, nil

}

func (m *MemcachedStore) Delete(k string) error {
	return m.adaptError(
		m.client.Delete(
			m.key(k)))
}

//Converts errors into onecache's types...
//If the error doesn't have an equivalent in the onecache package, it is returned as is
func (m *MemcachedStore) adaptError(err error) error {

	switch err {

	case nil:
		return nil
	case memcache.ErrCacheMiss:
		return onecache.ErrCacheMiss
	}

	return err
}

func (m *MemcachedStore) Flush() error {
	return m.client.DeleteAll()
}

func (m *MemcachedStore) Has(key string) bool {

	if _, err := m.Get(key); err != nil {
		return false
	}

	return true
}
