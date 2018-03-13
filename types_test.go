package onecache

import (
	"strings"
	"testing"
)

func TestDefaultKeyFunc(t *testing.T) {

	tt := []struct {
		expected, original string
	}{
		{"onecache:lanre", "lanre"},
		{"onecache:onecache", "onecache"},
	}

	for _, v := range tt {
		if !strings.EqualFold(v.expected, DefaultKeyFunc(v.original)) {
			t.Fatalf("an error occurred while checking if strings equal....")
		}
	}
}
