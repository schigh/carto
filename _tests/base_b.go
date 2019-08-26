// Code generated Mon, 26 Aug 2019 13:49:26 EDT by carto.  DO NOT EDIT.
package cartotests

import (
	"sync"
)

// Base0B wraps map[string]interface{}, and locks reads and writes with a mutex
type Base0B struct {
	mx       sync.RWMutex
	internal map[string]interface{}
}

// NewBase0B generates a new Base0B with a non-nil map
func NewBase0B() *Base0B {
	b := &Base0B{}
	b.internal = make(map[string]interface{})

	return b
}

// Get gets the interface{} keyed by string. Also returns bool value indicating whether the key exists in the map
func (b *Base0B) Get(key string) (value interface{}, ok bool) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	value, ok = b.internal[key]

	return
}

// Keys will return all keys in the Base0B's internal map
func (b *Base0B) Keys() (keys []string) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	keys = make([]string, len(b.internal))
	var i int
	for k := range b.internal {
		keys[i] = k
		i++
	}

	return
}

// Set will add an element to the Base0B's internal map with the specified key
func (b *Base0B) Set(key string, value interface{}) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.internal[key] = value
}

// Absorb will take all the keys and values from another Base0B's internal map and
// overwrite any existing keys
func (b *Base0B) Absorb(otherMap *Base0B) {
	b.mx.Lock()
	otherMap.mx.RLock()
	defer otherMap.mx.RUnlock()
	defer b.mx.Unlock()

	for k, v := range otherMap.internal {
		b.internal[k] = v
	}
}

// AbsorbMap will take all the keys and values from another map and overwrite any existing keys
func (b *Base0B) AbsorbMap(regularMap map[string]interface{}) {
	b.mx.Lock()
	defer b.mx.Unlock()

	for k, v := range regularMap {
		b.internal[k] = v
	}
}

// Delete will remove a interface{} from the map by key
func (b *Base0B) Delete(key string) {
	b.mx.Lock()
	defer b.mx.Unlock()

	delete(b.internal, key)
}

// Clear will remove all elements from the map
func (b *Base0B) Clear() {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.internal = make(map[string]interface{})
}
