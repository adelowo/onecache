package main

import (
	"bytes"
	"fmt"
	"reflect"
	"time"

	"github.com/adelowo/onecache"
	"github.com/adelowo/onecache/filesystem"
	"github.com/adelowo/onecache/memcached"
	"github.com/adelowo/onecache/memory"
	"github.com/adelowo/onecache/redis"
	"github.com/bradfitz/gomemcache/memcache"
	r "github.com/go-redis/redis"
)

//A custom type
//Might want to register this type with encoding/gob (Not really necessary though)
type user struct {
	Name string
}

func main() {

	marshal := onecache.NewCacheSerializer()

	i := &onecache.Item{time.Now().Add(time.Minute * 10), []byte("Lanre")}

	fileSystemCache(marshal)

	redisStore(marshal, i)

	memcachedStore()

	inMemoryStore()

}

func inMemoryStore() {
	//Setting a larger interval for gc would be better.
	//Something like time.Minute * 10
	//In other to prevent doing the same thing over and over again
	// if nothing really changed
	mem := memory.NewInMemoryStore(time.Minute * 2)
	fmt.Println(mem.Set("name", []byte("Lanre"), time.Minute*10))
	fmt.Println(mem.Get("name"))
	fmt.Println(mem.Set("occupation", []byte("What ?"), time.Second))
	fmt.Println(mem.Has("occupation"))
	fmt.Println(mem.Flush())
	fmt.Println(mem.Set("n", []byte("42"), time.Minute*1))
	fmt.Println(mem.Get("n"))
	fmt.Println(mem.Get("n"))
}

func memcachedStore() {
	m := memcached.NewMemcachedStore(
		memcache.New("127.0.0.1:11211"),
		"",
	)
	fmt.Println(m.Set("name", []byte("Rob Pike"), onecache.EXPIRES_DEFAULT))
	fmt.Println(m.Get("name"))
	fmt.Println(m.Flush())
}

func redisStore(marshal *onecache.CacheSerializer, i *onecache.Item) {
	opt := &r.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}
	redisStore := redis.NewRedisStore(opt, "")
	byt, _ := marshal.Serialize(i)
	fmt.Println(redisStore.Set("name", byt, time.Second*10))
	val, _ := redisStore.Get("name")
	fmt.Println(marshal.Serialize(val))
	fmt.Println(redisStore.Flush())
	fmt.Println(redisStore.Get("name"))
}

func fileSystemCache(marshal *onecache.CacheSerializer) {

	var store onecache.Store

	store = filesystem.MustNewFSStore("/home/adez/onecache_tmp", time.Minute*10)

	err := store.Set("profile", []byte("Lanre"), time.Second*60)

	fmt.Println(err)

	fmt.Println(store.Get("profile"))

	u := &user{"Lanre"}

	b, _ := marshal.Serialize(u)
	//Handle error

	fmt.Println(store.Set("user", b, time.Minute*1))

	newB, _ := store.Get("user")

	if !bytes.Equal(b, newB) {
		panic("OOPS")
	}
	//Convert back to a user struct
	newUser := new(user)

	marshal.DeSerialize(newB, newUser)

	//Check if the conversion was right
	if !reflect.DeepEqual(u, newUser) {
		panic("OOPS")
	}

	fmt.Println(store.Flush())

	fmt.Println(store.Get("unkownKey"))
}
