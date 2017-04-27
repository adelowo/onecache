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

//identifes a cached piece of data
type Item struct {
	ExpiresAt time.Time
	Data      interface{}
}

//Interface for all onecache store implementations
type Store interface {
	Set(key string, data interface{}, expires time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Flush() error
	Increment(key string, steps int) error
	Decrement(key string, steps int) error
}
