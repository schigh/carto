package main

//go:generate env

type mapTmpl struct {
	GenDate               string
	PackageName           string
	StructName            string
	Sync                  bool
	Mutex                 bool
	KeyTypePackage        string
	KeyTypePackageShort   string
	KeyTypeIsPointer      bool
	KeyType               string
	ValueType             string
	ValueTypePackage      string
	ValueTypePackageShort string
	ValueTypeIsPointer    bool
	InternalMapName       string
	ByReference           bool
	ReceiverName          string
	GetReturnsBool        bool
	LazyInstantiates      bool
	GetDefault            bool
}

const headTmpl = `// Code generated {{.GenDate}} by carto.  DO NOT EDIT.
package {{.PackageName}}
{{if .Sync}}
import (
	"sync"

	{{if .KeyTypePackage}}"{{.KeyTypePackage}}"{{end}}
	{{if .ValueTypePackage}}"{{.ValueTypePackage}}"{{end}}
)
{{end}}

// {{.StructName}} wraps map[{{.KeyType}}]{{.ValueType}}{{if .Mutex}}, and locks reads and writes with a mutex{{end}}
type {{.StructName}} struct {
	{{if .Mutex}}sync.RWMutex{{end}}
	{{.InternalMapName}} map[{{.KeyType}}]{{.ValueType}}
	{{if .LazyInstantiates}}onceToken sync.Once{{end}}
}{{if .LazyInstantiates}}{{else}}

// New{{.StructName}} generates a new {{.StructName}} with a non-nil map
func New{{.StructName}}() {{if .ByReference}}*{{end}}{{.StructName}} {
	{{.ReceiverName}} := {{if .ByReference}}&{{end}}{{.StructName}}{}
	{{.ReceiverName}}.{{.InternalMapName}} = make(map[{{.KeyType}}]{{.ValueType}})

	return {{.ReceiverName}}
}
{{end}}
`

const getTmpl = `
{{if .GetDefault}}// Get gets the {{.ValueType}} keyed by {{.KeyType}}.  If the key does not exist, a default {{.ValueType}} will be returned
func ({{.ReceiverName}} {{if .ByReference}}*{{end}}{{.StructName}}) Get(key {{.KeyType}}, dflt {{.ValueType}})(value {{.ValueType}}) {
	{{if .Mutex}}{{.ReceiverName}}.RLock()
	{{end}}		var ok bool
	value, ok = {{.ReceiverName}}.{{.InternalMapName}}[key]{{if .Mutex}}
	{{.ReceiverName}}.RUnlock(){{end}}
	if !ok {
		value = dflt
	}
	return
}{{else}}// Get gets the {{.ValueType}} keyed by {{.KeyType}}. {{if .GetReturnsBool}}Also returns bool value indicating whether the key exists in the map{{end}}
func ({{.ReceiverName}} {{if .ByReference}}*{{end}}{{.StructName}}) Get(key {{.KeyType}}) {{if .GetReturnsBool}}(value {{.ValueType}}, ok bool){{else}}(value {{.ValueType}}){{end}} {
	{{if .Mutex}}{{.ReceiverName}}.RLock()
	{{end}}		value{{if .GetReturnsBool}}, ok{{end}} = {{.ReceiverName}}.{{.InternalMapName}}[key]	{{if .Mutex}}
	{{.ReceiverName}}.RUnlock(){{end}}
	return
}{{end}}
`

const keysTmpl = `
// Keys will return all keys in the {{.StructName}}'s internal map
func ({{.ReceiverName}} {{if .ByReference}}*{{end}}{{.StructName}}) Keys() (keys []{{.KeyType}}) {
	{{if .Mutex}}{{.ReceiverName}}.RLock()
	{{end}}		keys = make([]{{.KeyType}}, len({{.ReceiverName}}.{{.InternalMapName}}))
	var i int
	for k := range {{.ReceiverName}}.{{.InternalMapName}} {
		keys[i] = k
		i++
	}{{if .Mutex}}
	{{.ReceiverName}}.RUnlock(){{end}}
	return
}
`

const setTmpl = `
// Set will add an element to the {{.StructName}}'s internal map with the specified key
func ({{.ReceiverName}} {{if .ByReference}}*{{end}}{{.StructName}}) Set(key {{.KeyType}}, value {{.ValueType}}) {
	{{if .Mutex}}{{.ReceiverName}}.Lock()
	{{end}}{{if .LazyInstantiates}}{{.ReceiverName}}.onceToken.Do(func() {
		{{.ReceiverName}}.{{.InternalMapName}} = make(map[{{.KeyType}}]{{.ValueType}})
	})
	{{end}}{{.ReceiverName}}.{{.InternalMapName}}[key] = value		{{if .Mutex}}
	{{.ReceiverName}}.Unlock(){{end}}
}
`

const absorbTmpl = `
// Absorb will take all the keys and values from another {{.StructName}}'s internal map and
// overwrite any existing keys
func ({{.ReceiverName}} {{if .ByReference}}*{{end}}{{.StructName}}) Absorb(otherMap {{if .ByReference}}*{{end}}{{.StructName}}) {
	{{if .Mutex}}{{.ReceiverName}}.Lock()
	otherMap.RLock()
	{{end}}{{if .LazyInstantiates}}{{.ReceiverName}}.onceToken.Do(func() {
		{{.ReceiverName}}.{{.InternalMapName}} = make(map[{{.KeyType}}]{{.ValueType}})
	})
	{{end}}for k, v := range otherMap.{{.InternalMapName}} {
		{{.ReceiverName}}.{{.InternalMapName}}[k] = v
	}{{if .Mutex}}
	otherMap.RUnlock()
	{{.ReceiverName}}.Unlock(){{end}}
}
`

const absorbMapTmpl = `
// AbsorbMap will take all the keys and values from another map and overwrite any existing keys
func ({{.ReceiverName}} {{if .ByReference}}*{{end}}{{.StructName}}) AbsorbMap(regularMap map[{{.KeyType}}]{{.ValueType}}) {
	{{if .Mutex}}{{.ReceiverName}}.Lock()
	{{end}}{{if .LazyInstantiates}}{{.ReceiverName}}.onceToken.Do(func() {
		{{.ReceiverName}}.{{.InternalMapName}} = make(map[{{.KeyType}}]{{.ValueType}})
	})
	{{end}}		for k, v := range regularMap {
		{{.ReceiverName}}.{{.InternalMapName}}[k] = v
	}{{if .Mutex}}
	{{.ReceiverName}}.Unlock(){{end}}
}
`

const deleteTmpl = `
// Delete will remove an item from the map by key
func ({{.ReceiverName}} {{if .ByReference}}*{{end}}{{.StructName}}) Delete(key {{.KeyType}}) {
	{{if .Mutex}}{{.ReceiverName}}.Lock()
	{{end}}		delete({{.ReceiverName}}.{{.InternalMapName}}, key)		{{if .Mutex}}
	{{.ReceiverName}}.Unlock(){{end}}
}
`

const clearTmpl = `
// Clear will remove all elements from the map
func ({{.ReceiverName}} {{if .ByReference}}*{{end}}{{.StructName}}) Clear() {
	{{if .Mutex}}{{.ReceiverName}}.Lock()
	{{end}}keys := make([]{{.KeyType}}, len({{.ReceiverName}}.{{.InternalMapName}}))
	var i int
	for k := range {{.ReceiverName}}.{{.InternalMapName}} {
		keys[i] = k
		i++
	}
	for _, k := range keys {
		delete({{.ReceiverName}}.{{.InternalMapName}}, k)
	}{{if .Mutex}}
	{{.ReceiverName}}.Unlock(){{end}}
}
`
