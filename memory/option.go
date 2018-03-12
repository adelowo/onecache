package memory

import "github.com/adelowo/onecache"

// Option defines options for creating a memory store
type Option func(i *InMemoryStore)

// BufferSize configures the store to allow a maximum of n
func BufferSize(n int) Option {
	return func(i *InMemoryStore) {
		i.bufferSize = n
	}
}

// KeyFunc allows for dynamic generation of cache keys
func KeyFunc(fn onecache.KeyFunc) Option {
	return func(i *InMemoryStore) {
		i.keyfn = fn
	}
}