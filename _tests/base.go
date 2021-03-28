// Code generated Sat, 27 Mar 2021 22:17:33 EDT by carto.  DO NOT EDIT.
package cartotests

import (
	"sync"
)

// Base wraps map[string]interface{}, and locks reads and writes with a mutex
type Base struct {
	mx   sync.RWMutex
	impl map[string]interface{}
}

// NewBase generates a new Base with a non-nil map
func NewBase() *Base {
	b := &Base{}
	b.impl = make(map[string]interface{})

	return b
}

// Get gets the interface{} keyed by string.
func (b *Base) Get(key string) (value interface{}) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	value = b.impl[key]

	return
}

// Keys will return all keys in the Base's internal map
func (b *Base) Keys() (keys []string) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	keys = make([]string, len(b.impl))
	var i int
	for k := range b.impl {
		keys[i] = k
		i++
	}

	return
}

// Set will add an element to the Base's internal map with the specified key
func (b *Base) Set(key string, value interface{}) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.impl[key] = value
}

// Absorb will take all the keys and values from another Base's internal map and
// overwrite any existing keys
func (b *Base) Absorb(otherMap *Base) {
	b.mx.Lock()
	otherMap.mx.RLock()
	defer otherMap.mx.RUnlock()
	defer b.mx.Unlock()

	for k, v := range otherMap.impl {
		b.impl[k] = v
	}
}

// AbsorbMap will take all the keys and values from another map and overwrite any existing keys
func (b *Base) AbsorbMap(regularMap map[string]interface{}) {
	b.mx.Lock()
	defer b.mx.Unlock()

	for k, v := range regularMap {
		b.impl[k] = v
	}
}

// Delete will remove a interface{} from the map by key
func (b *Base) Delete(key string) {
	b.mx.Lock()
	defer b.mx.Unlock()

	delete(b.impl, key)
}

// Clear will remove all elements from the map
func (b *Base) Clear() {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.impl = make(map[string]interface{})
}
