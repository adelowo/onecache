//Package redis provides a cache implementation of onecache using redis as the backend
package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/adelowo/onecache"
)

// Option is a redis option type
type Option func(r *RedisStore)

// ClientOptions is an Option type that allows configuring a redis client
func ClientOptions(opts *redis.Options) Option {
	return func(r *RedisStore) {
		r.client = redis.NewClient(opts)
	}
}

// CacheKeyGenerator allows configuring the cache key generation process
func CacheKeyGenerator(fn onecache.KeyFunc) Option {
	return func(r *RedisStore) {
		r.keyFn = fn
	}
}


type RedisStore struct {
	client *redis.Client

	keyFn onecache.KeyFunc
}

// New returns a new RedisStore by applying all options passed into it
// It also sets sensible defaults too
func New(opts ...Option) *RedisStore {
	r := &RedisStore{}

	for _, opt := range opts {
		opt(r)
	}

	if r.client == nil {
		redisOpts := &redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}

		ClientOptions(redisOpts)(r)
	}

	if r.keyFn == nil {
		r.keyFn = onecache.DefaultKeyFunc
	}


	return r
}

// Deprecated -- Use New instead
// Returns a new instance of the RedisStore
// If prefix is an empty string, the default cache prefix is used
func NewRedisStore(opts *redis.Options, prefix string) *RedisStore {
	return New(ClientOptions(opts))
}

func (r *RedisStore) Set(k string, data []byte, expires time.Duration) error {
	return r.client.Set(r.key(k), data, expires).Err()
}

func (r *RedisStore) Get(key string) ([]byte, error) {
	return r.client.Get(r.key(key)).Bytes()
}

func (r *RedisStore) Delete(key string) error {
	return r.client.Del(r.key(key)).Err()
}

func (r *RedisStore) Flush() error {
	return r.client.FlushDB().Err()
}

func (r *RedisStore) Has(key string) bool {

	if _, err := r.Get(key); err != nil {
		return false
	}

	return true
}

func (r *RedisStore) key(k string) string {
	return r.keyFn(k)
}
