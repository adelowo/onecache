package redis

import (
	"bytes"
	"flag"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/adelowo/onecache"
	"github.com/go-redis/redis"
)

var _ onecache.Store = &RedisStore{}

var redisStore *RedisStore

const TEST_PREFIX = "onecache_test:"

func TestMain(m *testing.M) {

	redisStore = NewRedisStore(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}, TEST_PREFIX)

	flag.Parse()
	os.Exit(m.Run())
}

var sampleData = []byte("Onecache")

func TestRedisStore_Set(t *testing.T) {

	err := redisStore.Set("name", sampleData, time.Minute*1)

	if err != nil {
		t.Fatalf(
			`An error occurred while interacting with redis....
			 %v`, err)
	}
}

func TestRedisStore_Get(t *testing.T) {

	val, err := redisStore.Get("name")

	if err != nil {
		t.Fatalf(
			`Could not get an item with the key, %s
			 due to an error %v`,
			"name", err)
	}

	if !reflect.DeepEqual(sampleData, val) {
		t.Fatalf("Expected %v.. \nGot %v instead", sampleData, val)
	}
}

func TestRedisStore_Delete(t *testing.T) {

	if err := redisStore.Delete("name"); err != nil {
		t.Fatalf(`
		Could not delete the key,%s due to an error ...%v
		`, "name", err)
	}
}

func TestRedisStore_Flush(t *testing.T) {

	//Save some data

	redisStore.Set("me", []byte("lanre"), onecache.EXPIRES_DEFAULT)
	redisStore.Set("animalName", []byte("Gopher"), onecache.EXPIRES_FOREVER)

	if err := redisStore.Flush(); err != nil {
		t.Fatalf(`
		An error occured while flushing the redis database... %v
		`, err)
	}

	//Manually inspect all data left in redis after flushing it's database
	cmd := redisStore.client.Keys("*")

	if err := cmd.Err(); err != nil {
		t.Fatalf(
			`An error occured while trying to get all keys
			stored in redis.. %v`,
			err)

	}
	res, err := cmd.Result()

	if err != nil {
		t.Fatalf(
			`An error occured while trying to get the
			result from redis... %v`, err)
	}

	if x := len(res); x != 0 {
		t.Fatalf(
			`There should be no more data stored in
			REDIS since we flushed the database...
			\n Expected %d.. Got %d instead `, 0, x)
	}
}

func TestRedisStore_GetUnknownKey(t *testing.T) {

	val, err := redisStore.Get("oops")

	if err == nil {
		t.Fatal("An Unknown key was encountered.. Yet we were able to retrieve it")

	}

	if !bytes.Equal(make([]byte, 0), val) {
		t.Fatalf(
			`Cache store should return a nil value.
			Since an unknown key was requested.. \n
			Got %v instead`, val)
	}
}

func TestNewRedisStore_DefaultPrefixIsUsedIfNoneIsProvided(t *testing.T) {

	s := NewRedisStore(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}, "")

	if !reflect.DeepEqual(defaultPrefix, s.prefix) {
		t.Fatalf(`
		Redis store prefix is invalid..
		Expected %s \n... Got %s`, defaultPrefix, s.prefix)
	}

}

func TestRedisStore_Has(t *testing.T) {
	s := NewRedisStore(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}, "")

	if ok := s.Has("name"); ok {
		t.Fatalf("Key %s is not supposed to exist in the cache", "name")
	}

	s.Set("name", []byte("Lanre"), time.Second*19)

	if ok := s.Has("name"); !ok {
		t.Fatalf("Key %s is supposed to exist in the cache", "name")
	}
}

func TestExtensibility(t *testing.T) {

	_, err := onecache.Get("redis")

	if err != nil {
		t.Fatalf("An error occurred %v", err)
	}
}
