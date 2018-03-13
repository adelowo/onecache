package filesystem

import "github.com/adelowo/onecache"

// Option is an optional type
type Option func(fs *FSStore)

func BaseDirectory(base string) Option {
	return func(fs *FSStore) {
		fs.baseDir = base
	}
}

func Serializer(serializer onecache.Serializer) Option {
	return func(fs *FSStore) {
		fs.b = serializer
	}
}

func CacheKeyGenerator(fn onecache.KeyFunc) Option {
	return func(fs *FSStore) {
		fs.keyFn = fn
	}
}
