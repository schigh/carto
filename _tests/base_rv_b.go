// Code generated Mon, 26 Aug 2019 13:49:26 EDT by carto.  DO NOT EDIT.
package cartotests

// Base0RV0B wraps map[string]interface{}
type Base0RV0B struct {
	internal map[string]interface{}
}

// NewBase0RV0B generates a new Base0RV0B with a non-nil map
func NewBase0RV0B() Base0RV0B {
	b := Base0RV0B{}
	b.internal = make(map[string]interface{})

	return b
}

// Get gets the interface{} keyed by string. Also returns bool value indicating whether the key exists in the map
func (b Base0RV0B) Get(key string) (value interface{}, ok bool) {
	value, ok = b.internal[key]

	return
}

// Keys will return all keys in the Base0RV0B's internal map
func (b Base0RV0B) Keys() (keys []string) {
	keys = make([]string, len(b.internal))
	var i int
	for k := range b.internal {
		keys[i] = k
		i++
	}

	return
}

// Set will add an element to the Base0RV0B's internal map with the specified key
func (b Base0RV0B) Set(key string, value interface{}) {
	b.internal[key] = value
}

// Absorb will take all the keys and values from another Base0RV0B's internal map and
// overwrite any existing keys
func (b Base0RV0B) Absorb(otherMap Base0RV0B) {
	for k, v := range otherMap.internal {
		b.internal[k] = v
	}
}

// AbsorbMap will take all the keys and values from another map and overwrite any existing keys
func (b Base0RV0B) AbsorbMap(regularMap map[string]interface{}) {
	for k, v := range regularMap {
		b.internal[k] = v
	}
}

// Delete will remove a interface{} from the map by key
func (b Base0RV0B) Delete(key string) {
	delete(b.internal, key)
}

// Clear will remove all elements from the map
func (b Base0RV0B) Clear() {
	b.internal = make(map[string]interface{})
}
