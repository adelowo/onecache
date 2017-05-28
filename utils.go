package onecache

import (
	"bytes"
	"encoding/gob"
	"strconv"
	"time"
)

//Helper method to check if an item is expired.
//Current usecase for this is for garbage collection
func (i *Item) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

type Serializer interface {
	Serialize(i interface{}) ([]byte, error)
	DeSerialize(data []byte, i interface{}) error
}

func NewCacheSerializer() *CacheSerializer {
	return &CacheSerializer{}
}

//Helper to serialize and deserialize types
type CacheSerializer struct {
}

//Convert a given type into a byte array
//Caveat -> Types you create might have to be registered with the encoding/gob package
func (b *CacheSerializer) Serialize(i interface{}) ([]byte, error) {

	if b, ok := i.([]byte); ok {
		return b, nil
	}

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(i); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

//Writes a byte array into a type.
func (b *CacheSerializer) DeSerialize(data []byte, i interface{}) error {

	if b, ok := i.(*[]byte); ok {
		*b = data
		return nil
	}

	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(i)
}

//Increment increases the value of an item by the specified number of steps
func Increment(val interface{}, steps int) (interface{}, error) {

	var ret interface{}

	switch val.(type) {

	case int:
		ret = val.(int) + steps

	case int32:
		ret = val.(int32) + int32(steps)

	case int64:
		ret = val.(int64) + int64(steps)

	case uint:
		ret = val.(uint) + uint(steps)

	case uint8:
		ret = val.(uint8) + uint8(steps)

	case uint16:
		ret = val.(uint16) + uint16(steps)

	case uint32:
		ret = val.(uint32) + uint32(steps)

	case uint64:
		ret = val.(uint64) + uint64(steps)

	case string:

		num, err := strconv.Atoi(val.(string))

		if err != nil {
			return -0, err
		}

		num += steps

		ret = strconv.Itoa(num)

	default:
		return -0, ErrCacheDataCannotBeIncreasedOrDecreased
	}

	return ret, nil

}

//Decrement decreases the value of an item by the specified number of steps
func Decrement(val interface{}, steps int) (interface{}, error) {
	return Increment(val, steps*-1)
}
