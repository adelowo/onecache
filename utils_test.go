package onecache

import (
	"reflect"
	"testing"
	"time"
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

func TestIncrease(t *testing.T) {

	var tests = []struct {
		expected interface{}
		give     interface{}
		steps    interface{}
	}{
		{52, 42, 10},
		{int32(32), int32(26), 6},
		{int64(30), int64(20), 10},
		{uint(8), uint(8), 0},
		{uint8(15), uint8(7), 8},
		{uint16(10), uint16(2), 8},
		{uint32(40), uint32(10), 30},
		{uint64(100), uint64(90), 10},
		{"42", "41", 1},
		{"30", "20", 10},
	}

	for _, v := range tests {
		val, err := Increment(v.give, v.steps.(int))

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

func TestIncrementForUnSupportedType(t *testing.T) {

	var tests = []struct {
		give  interface{}
		steps interface{}
	}{
		{true, 10},
		{"10.0", 2},
	}

	for _, v := range tests {
		_, err := Increment(v.give, v.steps.(int))

		if err == nil {
			t.Fatalf(
				`There should be an error on encountering an unsupported data type
				.. %v`, err)
		}
	}

}
