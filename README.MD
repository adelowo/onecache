### OneCache - A Go caching Library

[![Coverage Status](https://coveralls.io/repos/github/adelowo/onecache/badge.svg)](https://coveralls.io/github/adelowo/onecache)
[![Build Status](https://img.shields.io/travis/adelowo/onecache/master.svg?style=flat-square)](https://travis-ci.org/adelowo/onecache.svg?branch=master)

- [Installation](#install)
- [Supported stores](#stores)
- [Examples](#eg)
  - [3rd party stores](#extra)

<div id="install"></div>

### Installation

```go

$ go get -u github.com/adelowo/onecache

```

<div id="stores"></div>

Supported cache stores

- [x] InMemory - "memory"
- [x] Filesystem - "fs"
- [x] Memcached - "memcached"
- [x] Redis  - "redis"

> You can make use of `onecache.Get("memcached")` to access the memcached adapter.. Refer to the list above for other adapters' aliases. 


<div id="eg"></div>

```go

store, err := onecache.Get("fs") //filesystem store

if err != nil {
	return err
}

_ := store.Set("number", []byte("42"), time.Hour * 10)

return store.Get("number")

```

Take a look at the [_examples](https://github.com/adelowo/onecache/blob/master/types.go#L27-L35) for all stores to get an idea of all available methods

<div id="extra"></div>

```go
import (
          "github.com/adelowo/onecache"
)

//This is lazy loaded..The function only  gets executed when you request for the boltdb store
onecache.Extend("boltdb", func() onecache.Store {
	cfg, _ := config.New("path.json")
	return &Bolt{cfg.DBPath}
})

//Somewhere else
boltStore, err := onecache.Get("boltdb")

```

### LICENSE
MIT
