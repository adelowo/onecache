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

err := fs.Set("profile", "Used to be Human", time.Second*60)
	 
fmt.Println(err)
	 
fmt.Println(fs.Get("profile")) // "Used to be Human, nil"

fmt.Println(fs.Flush())

```

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

fmt.Println(redisStore.Set("profile", "Github\n", time.Minute*10))

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

fmt.Println(m.Set("name", "Lanre", onecache.EXPIRES_DEFAULT))

fmt.Println(m.Get("name"))

fmt.Println(m.Flush())


```

- InMemory

```go

import (
	"github.com/adelowo/onecache/memory"
)

mem := memory.newInMemoryStore()
	
fmt.Println(mem.Set("name", "Lanre", time.Minute*10))

fmt.Println(mem.Get("name"))
fmt.Println(mem.Get("unknown"))

fmt.Println(mem.Set("occupation", "Troll", time.Second))

fmt.Println(mem.Flush())

fmt.Println(mem.Set("n", 2, time.Minute*1))
fmt.Println(mem.Get("n"))
fmt.Println(mem.Increment("n", 10))
fmt.Println(mem.Get("n"))

```

### LICENSE
MIT
