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

OneCache also comes with ___garbage collection___. This is used by the filesystem and memory adapter to purge out expired items automatically. Please refer to the [examples][eg]


<div id="eg"></div>

[Examples containing all adapters can be found here][eg]

```go


var store onecache.Store

store = filesystem.MustNewFSStore("/home/adez/onecache_tmp", time.Minute*10)

err := store.Set("profile", []byte("Lanre"), time.Second*60)

fmt.Println(err)

fmt.Println(store.Get("profile"))

```

Some adapters like the `filesystem` and `memory` have a ___Garbage collection___ implementation. All
that is needed to call is `store.GC()`. Ideally, this should be called in a `ticker.C`. 

### LICENSE
MIT

[eg]: https://github.com/adelowo/onecache/blob/master/_examples/main.go
