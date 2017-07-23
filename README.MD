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

OneCache also comes with ___garbage collection___. This is used by the filesystem and memory adapter to purge out expired items automatically. Please refer to the [examples][eg]

> You can make use of `onecache.Get("memcached")` to access the memcached adapter.. Refer to the list above for other adapters' aliases. 


<div id="eg"></div>

[Examples containing all adapters can be found here][eg]

```go


var store onecache.Store

store = filesystem.MustNewFSStore("/home/adez/onecache_tmp", time.Minute*10)

err := store.Set("profile", []byte("Lanre"), time.Second*60)

fmt.Println(err)

fmt.Println(store.Get("profile"))

```

You can also make use of the `Get` to fetch a store..

```go

store, err := onecache.Get("fs") //filesystem store

if err != nil {
	return err
}

_ := store.Set("number", []byte("42"), time.Hour * 10)

return store.Get("number")

```

A couple things to keep in mind while making use of `onecache.Get("adapter_alias")` :

- For the filesystem and memory adapters, ___garbage collection___ isn't run automatically, you would have to manually do this by calling `GC()`.

```go

fs, _ := onecache.Get("fs")

go fs.GC(time.Minute * 10) //Check for and delete expired items every 10 minutes

//Or make use of the constructors provided by the adapters (auto GC) -- NewXXX()
```

- Redis : THe connection used is `localhost:6379` with the default database and an empty password.. While this would suffce for most usage, if you have custom connection requirements (like a password protected redis database, refer to the [examples][eg] too see how to make use of custom connection parameters.

- Memcached - The Memcached server connected to is located at the address `127.0.0.1:11211` which is usually the default. If you have custom requirements, please refer to the [examples][eg] to see how to make use of custom connection parameters.

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

[eg]: https://github.com/adelowo/onecache/blob/master/_examples/main.go
