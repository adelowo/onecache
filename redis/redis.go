//Package redis provides a cache implementation of onecache using redis as the backend
package redis

import (
	"time"

	"github.com/go-redis/redis"
)

//Default prefix to prevent collision with other key stored in redis
const PREFIX = "onecache:"

type RedisStore struct {
	client *redis.Client
	prefix string
}

//Returns a new instance of the RedisStore
//If prefix is an empty string, the default cache prefix is used
func NewRedisStore(opts *redis.Options, prefix string) *RedisStore {

	var p string

	if prefix == "" {
		p = PREFIX
	} else {
		p = prefix
	}

	return &RedisStore{redis.NewClient(opts), p}
}

func (r *RedisStore) Set(key string, data interface{}, expires time.Duration) error {
	return r.client.Set(r.key(key), data, expires).Err()
}

func (r *RedisStore) Get(key string) (interface{}, error) {

	val, err := r.client.Get(r.key(key)).Result()

	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisStore) Delete(key string) error {
	return r.client.Del(r.key(key)).Err()
}

func (r *RedisStore) key(k string) string {
	return r.prefix + k
}
