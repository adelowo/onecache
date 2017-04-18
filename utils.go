package onecache

import (
	"bytes"
	"encoding/gob"
	"time"
)

//Converts an item into bytes
func (i *Item) ToBytes() ([]byte, error) {

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(i); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

//Helper method to check if an item is expired.
//Current usecase for this is for garbage collection
func (i *Item) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

//Decodes bytes into an item struct
func BytesToItem(data []byte) (*Item, error) {

	i := new(Item)

	err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(i)

	return i, err
}
