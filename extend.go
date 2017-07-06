package onecache

import (
	"errors"
	"sync"
)

//AdapterFunc defines the structure of a cache store/backend to be registered
type AdapterFunc func() Store

//Extend registers the given adapter
//Note that if an adapter with the specified name exists, it would be overridden
func Extend(name string, fn AdapterFunc) {
	adapters.add(name, fn)
}

//Get resolves a cache store by name. A non nil error would be returned if the store was found.
func Get(name string) (Store, error) {
	return adapters.get(name)
}

type registeredAdapters struct {
	stores map[string]AdapterFunc
	lock   sync.RWMutex
}

func (r *registeredAdapters) add(name string, fn AdapterFunc) {
	r.lock.Lock()
	r.stores[name] = fn
	r.lock.Unlock()
}

func (r *registeredAdapters) get(name string) (Store, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if fn, ok := r.stores[name]; ok {
		return fn(), nil
	}

	return nil, errors.New("Adapter not found")
}

var adapters *registeredAdapters

func init() {
	adapters = &registeredAdapters{
		stores: make(map[string]AdapterFunc, 10),
	}
}
