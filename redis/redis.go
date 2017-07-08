//Package redis provides a cache implementation of onecache using redis as the backend
package redis

import (
	"time"

	"github.com/adelowo/onecache"
	"github.com/go-redis/redis"
)

//Default prefix to prevent collision with other key stored in redis
const defaultPrefix = "onecache:"

type RedisStore struct {
	client *redis.Client
	prefix string
}

func init() {
	onecache.Extend("redis", func() onecache.Store {
		//Default for most usage..
		//Can make use of NewRedisStore() for custom settings
		return NewRedisStore(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}, defaultPrefix)
	})
}

//Returns a new instance of the RedisStore
//If prefix is an empty string, the default cache prefix is used
func NewRedisStore(opts *redis.Options, prefix string) *RedisStore {

	var p string

	if prefix == "" {
		p = defaultPrefix
	} else {
		p = prefix
	}

	return &RedisStore{redis.NewClient(opts), p}
}

func (r *RedisStore) Set(k string, data []byte, expires time.Duration) error {
	return r.client.Set(r.key(k), data, expires).Err()
}

func (r *RedisStore) Get(key string) ([]byte, error) {

	val, err := r.client.Get(r.key(key)).Bytes()

	if err != nil {
		return nil, err
	}

	return val, nil
}

func (r *RedisStore) Delete(key string) error {
	return r.client.Del(r.key(key)).Err()
}

func (r *RedisStore) Flush() error {
	return r.client.FlushDb().Err()
}

func (r *RedisStore) Has(key string) bool {

	if _, err := r.Get(key); err != nil {
		return false
	}

	return true
}

func (r *RedisStore) key(k string) string {
	return r.prefix + k
}
