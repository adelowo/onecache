package onecache

import (
	"errors"
	"time"
)

const (
	EXPIRES_DEFAULT = time.Duration(0)
	EXPIRES_FOREVER = time.Duration(-1)
)

var (
	ErrCacheMiss                             = errors.New("Key not found")
	ErrCacheNotStored                        = errors.New("Data not stored")
	ErrCacheNotSupported                     = errors.New("Operation not supported")
	ErrCacheDataCannotBeIncreasedOrDecreased = errors.New(`
		Data isn't an integer/string type. Hence, it cannot be increased or decreased`)
)

//Item identifes a cached piece of data
type Item struct {
	ExpiresAt time.Time
	Data      []byte
}

//Interface for all onecache store implementations
type Store interface {
	Set(key string, data []byte, expires time.Duration) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	Flush() error
	Has(key string) bool
}

//Some stores like redis and memcache automatically clear out the cache
//But for the filesystem and in memory, this cannot be said.
//Stores that have to manually clear out the cached data should implement this method.
//It's implementation should re run this function everytime the interval is reached
//Say every 5 minutes.
type GarbageCollector interface {
	GC(interval time.Duration)
}
