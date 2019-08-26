// Code generated Mon, 26 Aug 2019 13:49:26 EDT by carto.  DO NOT EDIT.
package cartotests

// Base0RV0D wraps map[string]interface{}
type Base0RV0D struct {
	internal map[string]interface{}
}

// NewBase0RV0D generates a new Base0RV0D with a non-nil map
func NewBase0RV0D() Base0RV0D {
	b := Base0RV0D{}
	b.internal = make(map[string]interface{})

	return b
}

// Get gets the interface{} keyed by string.  If the key does not exist, a default interface{} will be returned
func (b Base0RV0D) Get(key string, dflt interface{}) (value interface{}) {
	var ok bool
	value, ok = b.internal[key]
	if !ok {
		value = dflt
	}

	return
}

// Keys will return all keys in the Base0RV0D's internal map
func (b Base0RV0D) Keys() (keys []string) {
	keys = make([]string, len(b.internal))
	var i int
	for k := range b.internal {
		keys[i] = k
		i++
	}

	return
}

// Set will add an element to the Base0RV0D's internal map with the specified key
func (b Base0RV0D) Set(key string, value interface{}) {
	b.internal[key] = value
}

// Absorb will take all the keys and values from another Base0RV0D's internal map and
// overwrite any existing keys
func (b Base0RV0D) Absorb(otherMap Base0RV0D) {
	for k, v := range otherMap.internal {
		b.internal[k] = v
	}
}

// AbsorbMap will take all the keys and values from another map and overwrite any existing keys
func (b Base0RV0D) AbsorbMap(regularMap map[string]interface{}) {
	for k, v := range regularMap {
		b.internal[k] = v
	}
}

// Delete will remove a interface{} from the map by key
func (b Base0RV0D) Delete(key string) {
	delete(b.internal, key)
}

// Clear will remove all elements from the map
func (b Base0RV0D) Clear() {
	b.internal = make(map[string]interface{})
}
