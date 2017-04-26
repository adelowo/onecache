//Package memcached is a memcached store for onecache
package memcached

import (
	"time"

	"github.com/adelowo/onecache"
	"github.com/bradfitz/gomemcache/memcache"
)

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

func (m *MemcachedStore) Set(k string, data interface{}, expires time.Duration) error {

	i := &onecache.Item{Data: data}

	b, err := i.Bytes()

	if err != nil {
		return err
	}

	item := &memcache.Item{
		Key:        m.key(k),
		Value:      b,
		Expiration: int32(expires / time.Second),
	}

	return m.client.Set(item)
}
