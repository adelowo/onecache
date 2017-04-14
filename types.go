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
	ErrCacheMiss         = errors.New("Key not found")
	ErrCacheNotStored    = errors.New("Data not stored")
	ErrCacheNotSupported = errors.New("Operation not supported")
)

