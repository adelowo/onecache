package onecache

import (
	"reflect"
	"testing"

	"github.com/adelowo/onecache/mocks"
)

func TestExtensibility(t *testing.T) {

	dummyStore := &mocks.Store{}

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
