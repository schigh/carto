// Code generated Mon, 26 Aug 2019 13:49:26 EDT by carto.  DO NOT EDIT.
package cartotests

import (
	"sync"
)

// Base0LZ wraps map[string]interface{}, and locks reads and writes with a mutex
type Base0LZ struct {
	mx        sync.RWMutex
	internal  map[string]interface{}
	onceToken sync.Once
}

// Get gets the interface{} keyed by string.
func (b *Base0LZ) Get(key string) (value interface{}) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	value = b.internal[key]

	return
}

// Keys will return all keys in the Base0LZ's internal map
func (b *Base0LZ) Keys() (keys []string) {
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

// Set will add an element to the Base0LZ's internal map with the specified key
func (b *Base0LZ) Set(key string, value interface{}) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.onceToken.Do(func() {
		b.internal = make(map[string]interface{})
	})
	b.internal[key] = value
}

// Absorb will take all the keys and values from another Base0LZ's internal map and
// overwrite any existing keys
func (b *Base0LZ) Absorb(otherMap *Base0LZ) {
	b.mx.Lock()
	otherMap.mx.RLock()
	defer otherMap.mx.RUnlock()
	defer b.mx.Unlock()

	b.onceToken.Do(func() {
		b.internal = make(map[string]interface{})
	})
	for k, v := range otherMap.internal {
		b.internal[k] = v
	}
}

// AbsorbMap will take all the keys and values from another map and overwrite any existing keys
func (b *Base0LZ) AbsorbMap(regularMap map[string]interface{}) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.onceToken.Do(func() {
		b.internal = make(map[string]interface{})
	})
	for k, v := range regularMap {
		b.internal[k] = v
	}
}

// Delete will remove a interface{} from the map by key
func (b *Base0LZ) Delete(key string) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.onceToken.Do(func() {
		b.internal = make(map[string]interface{})
	})
	delete(b.internal, key)
}

// Clear will remove all elements from the map
func (b *Base0LZ) Clear() {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.internal = make(map[string]interface{})
}
