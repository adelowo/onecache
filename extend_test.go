package onecache

import (
	"reflect"
	"testing"
	"time"
)

func TestExtensibility(t *testing.T) {

	dummyStore := &mock{}

	Extend("dummy", func() Store {
		return dummyStore
	})

	s, err := Get("dummy")

	if err != nil {
		t.Fatalf("An error occurred...%v", err)
	}

	if !reflect.DeepEqual(dummyStore, s) {
		t.Fatalf("Stores differ...Expected %v \n Got %v", dummyStore, s)
	}
}

type mock struct {
}

func (_m *mock) Delete(key string) error {
	return nil
}

func (_m *mock) Flush() error {
	return nil
}

// Get provides a mock function with given fields: key
func (_m *mock) Get(key string) ([]byte, error) {
	return nil, nil
}

// Has provides a mock function with given fields: key
func (_m *mock) Has(key string) bool {
	return true
}

// Set provides a mock function with given fields: key, data, expires
func (_m *mock) Set(key string, data []byte, expires time.Duration) error {
	return nil
}
