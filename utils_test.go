package onecache

import (
	"testing"
	"time"
	"reflect"
)

func TestItem_IsExpired(t *testing.T) {

	item := &Item{ExpiresAt: time.Now().Add(-2 * time.Minute), Data: "Ping-Pong"}

	if !item.IsExpired() {
		t.Fatal("Item should be expired since it's expiration date is set 2 minutes backwards")
	}
}

func TestBytesToItem(t *testing.T) {

	item := &Item{ExpiresAt: time.Now(), Data: "Ping-Pong"}

	b, err := item.Bytes()

	if err != nil {
		t.Error(err)
	}

	i, err := BytesToItem(b)

	if err != nil {
		t.Error(err)
	}

	if !item.ExpiresAt.Equal(i.ExpiresAt) {
		t.Fatalf("Time should equal.. Expected %v \n Got %v", item.ExpiresAt, i.ExpiresAt)
	}

	if !reflect.DeepEqual(item.Data, i.Data) {
		t.Fatalf("Data not equal.. Expected %v \n. Got %v", item.ExpiresAt, i.ExpiresAt)
	}

}
