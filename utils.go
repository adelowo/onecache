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

type Marshaller interface {
	UnMarshallBytes(data []byte) (*Item, error)
	MarshalBytes(i *Item) ([]byte, error)
}

func NewBytesItemMarshaller() *BytesItemMarshaller {
	return &BytesItemMarshaller{}
}

type BytesItemMarshaller struct {

}

func (b *BytesItemMarshaller) MarshalBytes(i *Item) ([]byte, error){
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(i); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (b *BytesItemMarshaller) UnMarshallBytes(data []byte) (*Item, error) {
	i := new(Item)

	err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(i)

	return i, err
}

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

func Decrement(val interface{}, steps int) (interface{}, error) {
	return Increment(val, steps*-1)
}
