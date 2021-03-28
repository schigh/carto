// Code generated Sat, 27 Mar 2021 22:17:33 EDT by carto.  DO NOT EDIT.
package cartotests

import (
	"sync"
)

// Base0LZ0B wraps map[string]interface{}, and locks reads and writes with a mutex
type Base0LZ0B struct {
	mx        sync.RWMutex
	impl      map[string]interface{}
	onceToken sync.Once
}

// Get gets the interface{} keyed by string. Also returns bool value indicating whether the key exists in the map
func (b *Base0LZ0B) Get(key string) (value interface{}, ok bool) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	value, ok = b.impl[key]

	return
}

// Keys will return all keys in the Base0LZ0B's internal map
func (b *Base0LZ0B) Keys() (keys []string) {
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

// Set will add an element to the Base0LZ0B's internal map with the specified key
func (b *Base0LZ0B) Set(key string, value interface{}) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.onceToken.Do(func() {
		b.impl = make(map[string]interface{})
	})
	b.impl[key] = value
}

// Absorb will take all the keys and values from another Base0LZ0B's internal map and
// overwrite any existing keys
func (b *Base0LZ0B) Absorb(otherMap *Base0LZ0B) {
	b.mx.Lock()
	otherMap.mx.RLock()
	defer otherMap.mx.RUnlock()
	defer b.mx.Unlock()

	b.onceToken.Do(func() {
		b.impl = make(map[string]interface{})
	})
	for k, v := range otherMap.impl {
		b.impl[k] = v
	}
}

// AbsorbMap will take all the keys and values from another map and overwrite any existing keys
func (b *Base0LZ0B) AbsorbMap(regularMap map[string]interface{}) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.onceToken.Do(func() {
		b.impl = make(map[string]interface{})
	})
	for k, v := range regularMap {
		b.impl[k] = v
	}
}

// Delete will remove a interface{} from the map by key
func (b *Base0LZ0B) Delete(key string) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.onceToken.Do(func() {
		b.impl = make(map[string]interface{})
	})
	delete(b.impl, key)
}

// Clear will remove all elements from the map
func (b *Base0LZ0B) Clear() {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.impl = make(map[string]interface{})
}
