package onecache

import (
	"bytes"
	"encoding/gob"
)

func (i *Item) ToBytes() ([]byte, error) {

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(i); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
