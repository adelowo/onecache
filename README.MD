### OneCache - A Go caching Library

[![Coverage Status](https://coveralls.io/repos/github/adelowo/onecache/badge.svg)](https://coveralls.io/github/adelowo/onecache)
[![Build Status](https://img.shields.io/travis/adelowo/onecache/master.svg?style=flat-square)](https://travis-ci.org/adelowo/onecache.svg?branch=master)

- [Installation](#install)
- [Supported stores](#stores)
- [Examples](#eg)

<div id="install"></div>

### Installation

```go

$ go get -u github.com/adelowo/onecache

```

> Version 2 was recently released and it fixed a number of issues as per how data was cached (the main bone of contention being the use of `interface{}`. Take a look at this [issue](https://github.com/adelowo/onecache/issues/7)). Version 1 would still be maintained going forward ([branch here](https://github.com/adelowo/onecache/tree/1.x)) but newer development would be based on v2.

<div id="stores"></div>

Supported cache stores

- [x] InMemory
- [x] Filesystem
- [x] Memcached
- [x] Redis


<div id="eg"></div>

> Have a look at the [interface](https://github.com/adelowo/onecache/blob/master/types.go#L27-L35) for all stores to get an idea of all available methods

- Filesystem

```go

import (
	"github.com/adelowo/onecache/filesystem"
)

var fs onecache.Store

fs = filesystem.MustNewFSStore("/some/writeable/directory") //Would be created if it does not exists

err := fs.Set("profile", []byte("Used to be Human"), time.Second*60)
	 
fmt.Println(err)
	 
fmt.Println(fs.Get("profile"))

fmt.Println(fs.Flush()) //Clear out all cached data

```

> You are not limited to storing primitives ...Please check out the [examples](https://github.com/adelowo/onecache/blob/master/_examples/main.go)

- Redis 

```go

import (
	"github.com/adelowo/onecache/redis"
	r "github.com/go-redis/redis"
)

opt := &r.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
}

redisStore := redis.NewRedisStore(opt, "")

fmt.Println(redisStore.Set("profile", []byte("Github\n"), time.Minute*10))

val, err := redisStore.Get("profile")

fmt.Println(val, err)

fmt.Println(redisStore.Flush())


```

- Memcached

```go

import(
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/adelowo/onecache/memcached"
)

m := memcached.NewMemcachedStore(memcache.New("127.0.0.1:11211"),"",)

fmt.Println(m.Set("name", []byte("Lanre"), onecache.EXPIRES_DEFAULT))

fmt.Println(m.Get("name"))

fmt.Println(m.Flush())


```

- InMemory

```go

import (
	"github.com/adelowo/onecache/memory"
)

mem := memory.newInMemoryStore()
	
fmt.Println(mem.Set("name", []byte("Lanre"), time.Minute*10))

fmt.Println(mem.Get("name"))
fmt.Println(mem.Get("unknown"))

fmt.Println(mem.Set("hobby", []byte("Trolling"), time.Second))

fmt.Println(mem.Flush())

```

### LICENSE
MIT
