package redis

import (
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

func TestRedisStore_Set(t *testing.T) {

	err := redisStore.Set("name", "Lanre", time.Minute*1)

	if err != nil {
		t.Fatalf("An error occurred while interacting with redis.... %v", err)
	}
}

func TestRedisStore_Get(t *testing.T) {

	val, err := redisStore.Get("name")

	if err != nil {
		t.Fatalf("Could not get an item with the key, %s due to an error %v", "name", err)
	}

	if !reflect.DeepEqual("Lanre", val) {
		t.Fatalf("Expected %s.. Got %s instead", "Lanre", val)
	}
}

func TestRedisStore_Delete(t *testing.T) {

	if err := redisStore.Delete("name"); err != nil {
		t.Fatalf("Could not delete the key,%s due to an error ...%v", "name", err)
	}
}

func TestRedisStore_Flush(t *testing.T) {

	//Save some data

	redisStore.Set("me", "you", onecache.EXPIRES_DEFAULT)
	redisStore.Set("animalName", "Gopher", onecache.EXPIRES_FOREVER)

	if err := redisStore.Flush(); err != nil {
		t.Fatalf("An error occured while flushing the redis database... %v", err)
	}

	cmd := redisStore.client.Keys("*")

	if err := cmd.Err(); err != nil {
		t.Fatalf("An error occured while trying to get all keys stored in redis.. %v", err)
	}

	res, err := cmd.Result()

	if err != nil {
		t.Fatalf("An error occured while trying to get the result from redis... %v", err)
	}

	if x := len(res); x != 0 {
		t.Fatalf("There should be no more data stored in REDIS since we flushed the database...\n Expected %d.. Got %d instead ", 0, x)
	}
}

func TestRedisStore_GetUnknownKey(t *testing.T) {

	val, err := redisStore.Get("oops")

	if err == nil {
		t.Fatal("An Unknown key was encountered.. Yet we were able to retrieve it")
	}

	if "" != val {
		t.Fatalf("Should return an empty string.. Got %s instead", val)
	}
}

func TestNewRedisStore_DefaultPrefixIsUsedIfNoneIsProvided(t *testing.T) {

	s := NewRedisStore(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}, "")

	if !reflect.DeepEqual(PREFIX, s.prefix) {
		t.Fatalf("Redis store prefix is invalid.. Expected %s \n... Got %s", PREFIX, s.prefix)
	}

}

func TestRedisStore_Increment(t *testing.T) {
	var tests = []struct {
		key      string
		give     interface{}
		expected interface{}
		steps    int
	}{
		{"name", "40", "42", 2},
		{"int", "100", "102", 2},
	}

	for _, v := range tests {
		redisStore.Set(v.key, v.give, time.Second*2)
		err := redisStore.Increment(v.key, v.steps)

		val, _ := redisStore.Get(v.key)

		if err != nil {
			t.Fatalf("An error occurred... %v", err)
		}

		if !reflect.DeepEqual(v.expected, val) {
			t.Fatalf(
				"Differs.. Expected %v .\n Got %v instead",
				v.expected, val)
		}
	}
}
