![carto](./_img/carto.png)

Carto is a command line tool that can help you make typed Go maps quickly and easily.  

```
> carto --help
```

```
C A R T O

🌐 Maps made easy

usage:
carto -p <package> -s <structname> -k <keytype> -v <valuetype> [options]

-p    (string)      Package name (required)
-s    (string)      Struct name  (required)
-k    (string)      Key type     (required)
-v    (string)      Value type   (required)
-r    (string)      Receiver name (defaults to lowercase first char of
                    struct name)
-o    (string)      Output file path (if omitted, prints to STDOUT)
-b    (bool)        "Get" return signature includes a bool value indicating
                    if the key exists in the internal map
-d    (bool)        "Get" signature has second parameter for default return
                    value when key does not exist in the internal map
-lz   (bool)        Will lazy-instantiate the internal map when a write
                    operation is used
-version            Print version and exit
```

## Getting Carto

You can install carto by using `go get`

```
> go get -u github.com/schigh/carto
```

## Options

* `p` (string, required) - The package of the generated struct

* `s` (string, required) - The name of the generated struct

* `k` (string, required) - The key type (see below regarding key and value syntax)

* `v` (string, required) - The value type (see below regarding key and value syntax)

* `r` (string) - The receiver name.  By default, it is the lowercase first letter of the generated struct name

* `o` (string) - Output file.  By default, carto prints to STDOUT

* `b` (bool) - `Get` will return an additional boolean value indicating whether or not the value was found within the internal map
* `d` (bool) - `Get` will take an additional default value to be returned if no value exists for the specified key
* `lz` (bool) - The generated struct will lazily (via  `sync.Once`) instantiate its internal map

## Examples

A struct `MyMap` that wraps `map[string]*zerolog.Logger`:

```
> carto -p foo -s MyMap -k string -v '*github.com/rs/zerolog.Logger'
```
<details>
<summary>Generated Source</summary>
<p>

```go
// Code generated Tue, 27 Aug 2019 10:27:43 EDT by carto.  DO NOT EDIT.
package foo

import (
	"sync"

	"github.com/rs/zerolog"
)

// MyMap wraps map[string]*zerolog.Logger, and locks reads and writes with a mutex
type MyMap struct {
	mx   sync.RWMutex
	impl map[string]*zerolog.Logger
}

// NewMyMap generates a new MyMap with a non-nil map
func NewMyMap() *MyMap {
	m := &MyMap{}
	m.impl = make(map[string]*zerolog.Logger)

	return m
}

// Get gets the *zerolog.Logger keyed by string.
func (m *MyMap) Get(key string) (value *zerolog.Logger) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	value = m.impl[key]

	return
}

// Keys will return all keys in the MyMap's internal map
func (m *MyMap) Keys() (keys []string) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	keys = make([]string, len(m.impl))
	var i int
	for k := range m.impl {
		keys[i] = k
		i++
	}

	return
}

// Set will add an element to the MyMap's internal map with the specified key
func (m *MyMap) Set(key string, value *zerolog.Logger) {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.impl[key] = value
}

// Absorb will take all the keys and values from another MyMap's internal map and
// overwrite any existing keys
func (m *MyMap) Absorb(otherMap *MyMap) {
	m.mx.Lock()
	otherMap.mx.RLock()
	defer otherMap.mx.RUnlock()
	defer m.mx.Unlock()

	for k, v := range otherMap.impl {
		m.impl[k] = v
	}
}

// AbsorbMap will take all the keys and values from another map and overwrite any existing keys
func (m *MyMap) AbsorbMap(regularMap map[string]*zerolog.Logger) {
	m.mx.Lock()
	defer m.mx.Unlock()

	for k, v := range regularMap {
		m.impl[k] = v
	}
}

// Delete will remove a *zerolog.Logger from the map by key
func (m *MyMap) Delete(key string) {
	m.mx.Lock()
	defer m.mx.Unlock()

	delete(m.impl, key)
}

// Clear will remove all elements from the map
func (m *MyMap) Clear() {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.impl = make(map[string]*zerolog.Logger)
}
```

</p>
</details>

The previous struct, lazy-instantiated:

```
> carto -p foo -s MyMap -k string -v '*github.com/rs/zerolog.Logger' -lz
```
<details>
<summary>Generated Source</summary>
<p>

```go
// Code generated Tue, 27 Aug 2019 10:28:46 EDT by carto.  DO NOT EDIT.
package foo

import (
	"sync"

	"github.com/rs/zerolog"
)

// MyMap wraps map[string]*zerolog.Logger, and locks reads and writes with a mutex
type MyMap struct {
	mx        sync.RWMutex
	impl      map[string]*zerolog.Logger
	onceToken sync.Once
}

// Get gets the *zerolog.Logger keyed by string.
func (m *MyMap) Get(key string) (value *zerolog.Logger) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	value = m.impl[key]

	return
}

// Keys will return all keys in the MyMap's internal map
func (m *MyMap) Keys() (keys []string) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	keys = make([]string, len(m.impl))
	var i int
	for k := range m.impl {
		keys[i] = k
		i++
	}

	return
}

// Set will add an element to the MyMap's internal map with the specified key
func (m *MyMap) Set(key string, value *zerolog.Logger) {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.onceToken.Do(func() {
		m.impl = make(map[string]*zerolog.Logger)
	})
	m.impl[key] = value
}

// Absorb will take all the keys and values from another MyMap's internal map and
// overwrite any existing keys
func (m *MyMap) Absorb(otherMap *MyMap) {
	m.mx.Lock()
	otherMap.mx.RLock()
	defer otherMap.mx.RUnlock()
	defer m.mx.Unlock()

	m.onceToken.Do(func() {
		m.impl = make(map[string]*zerolog.Logger)
	})
	for k, v := range otherMap.impl {
		m.impl[k] = v
	}
}

// AbsorbMap will take all the keys and values from another map and overwrite any existing keys
func (m *MyMap) AbsorbMap(regularMap map[string]*zerolog.Logger) {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.onceToken.Do(func() {
		m.impl = make(map[string]*zerolog.Logger)
	})
	for k, v := range regularMap {
		m.impl[k] = v
	}
}

// Delete will remove a *zerolog.Logger from the map by key
func (m *MyMap) Delete(key string) {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.onceToken.Do(func() {
		m.impl = make(map[string]*zerolog.Logger)
	})
	delete(m.impl, key)
}

// Clear will remove all elements from the map
func (m *MyMap) Clear() {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.impl = make(map[string]*zerolog.Logger)
}
```

</p>
</details>

With a default return value:

```
> carto -p foo -s MyMap -k string -v '*github.com/rs/zerolog.Logger' -d
```
<details>
<summary>Generated Source (`Get` func)</summary>
<p>

```go
...

// Get gets the *zerolog.Logger keyed by string.  If the key does not exist, a default // Get gets the *zerolog.Logger keyed by string.  If the key does not exist, a default *zerolog.Logger will be returned
func (m *MyMap) Get(key string, dflt *zerolog.Logger) (value *zerolog.Logger) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	var ok bool
	value, ok = m.impl[key]
	if !ok {
		value = dflt
	}

	return
}

...
```

</p>
</details>

With a second boolean return value:

```
> carto -p foo -s MyMap -k string -v '*github.com/rs/zerolog.Logger' -b
```
<details>
<summary>Generated Source (`Get` func)</summary>
<p>

```go
...

// Get gets the *zerolog.Logger keyed by string. Also returns bool value indicating whether the key exists in the map
func (m *MyMap) Get(key string) (value *zerolog.Logger, ok bool) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	value, ok = m.impl[key]

	return
}

...
```

</p>
</details>