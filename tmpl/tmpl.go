// Code generated by /tools/tmpl.go on Mon, 26 Aug 2019 15:43:31 EDT. DO NOT EDIT
package tmpl

/*
If you want to edit the templates directly, you can do so in
under tmpl/_gotemplate
The templates are maintained there for greater readability.
*/

// MapTmpl contains all the necessary data to generate a carto map
type MapTmpl struct {
	GenDate               string
	PackageName           string
	StructName            string
	KeyTypePackage        string
	KeyType               string
	ValueType             string
	ValueTypePackage      string
	InternalMapName       string
	ByReference           bool
	ReceiverName          string
	GetReturnsBool        bool
	LazyInstantiates      bool
	GetDefault            bool
}

// HeadTmpl is the file header, including imports and struct declaration.
// If lazy map instantiation is _not_ enabled, this also wraps the New... func.
const HeadTmpl = `// Code generated {{.GenDate}} by carto.  DO NOT EDIT.
package {{.PackageName}}
import (
{{if or .ByReference .LazyInstantiates}}	"sync"
{{end}}{{if .KeyTypePackage}}	"{{.KeyTypePackage}}"
{{end}}{{if .ValueTypePackage}}	"{{.ValueTypePackage}}"
{{end}})

// {{.StructName}} wraps map[{{.KeyType}}]{{.ValueType}}{{if .ByReference}}, and locks reads and writes with a mutex{{end}}
type {{.StructName}} struct {
{{if .ByReference}}	mx sync.RWMutex
{{end}}	{{.InternalMapName}} map[{{.KeyType}}]{{.ValueType}}
{{if .LazyInstantiates}}	onceToken sync.Once
{{end}}
}
{{if .LazyInstantiates}}{{else}}
// New{{.StructName}} generates a new {{.StructName}} with a non-nil map
func New{{.StructName}}() {{if .ByReference}}*{{end}}{{.StructName}} {
	{{.ReceiverName}} := {{if .ByReference}}&{{end}}{{.StructName}}{}
	{{.ReceiverName}}.{{.InternalMapName}} = make(map[{{.KeyType}}]{{.ValueType}})

	return {{.ReceiverName}}
}
{{end}}`

// GetTmpl wraps the `Get` func
const GetTmpl = `{{if .GetDefault}}// Get gets the {{.ValueType}} keyed by {{.KeyType}}.  If the key does not exist, a default {{.ValueType}} will be returned
func ({{.ReceiverName}} {{if .ByReference}}*{{end}}{{.StructName}}) Get(key {{.KeyType}}, dflt {{.ValueType}})(value {{.ValueType}}) {
{{if .ByReference}}	{{.ReceiverName}}.mx.RLock()
	defer {{.ReceiverName}}.mx.RUnlock()

{{end}}	var ok bool
	value, ok = {{.ReceiverName}}.{{.InternalMapName}}[key]
	if !ok {
		value = dflt
	}

	return
}
{{else}}// Get gets the {{.ValueType}} keyed by {{.KeyType}}. {{if .GetReturnsBool}}Also returns bool value indicating whether the key exists in the map{{end}}
func ({{.ReceiverName}} {{if .ByReference}}*{{end}}{{.StructName}}) Get(key {{.KeyType}}) {{if .GetReturnsBool}}(value {{.ValueType}}, ok bool){{else}}(value {{.ValueType}}){{end}} {
{{if .ByReference}}	{{.ReceiverName}}.mx.RLock()
	defer {{.ReceiverName}}.mx.RUnlock()

{{end}}	value{{if .GetReturnsBool}}, ok{{end}} = {{.ReceiverName}}.{{.InternalMapName}}[key]

	return
}
{{end}}`

// KeysTmpl wraps the `Keys` func
const KeysTmpl = `// Keys will return all keys in the {{.StructName}}'s internal map
func ({{.ReceiverName}} {{if .ByReference}}*{{end}}{{.StructName}}) Keys() (keys []{{.KeyType}}) {
{{if .ByReference}}	{{.ReceiverName}}.mx.RLock()
	defer {{.ReceiverName}}.mx.RUnlock()

{{end}}	keys = make([]{{.KeyType}}, len({{.ReceiverName}}.{{.InternalMapName}}))
	var i int
	for k := range {{.ReceiverName}}.{{.InternalMapName}} {
		keys[i] = k
		i++
	}

	return
}
`

// SetTmpl wraps the `Set` func
const SetTmpl = `// Set will add an element to the {{.StructName}}'s internal map with the specified key
func ({{.ReceiverName}} {{if .ByReference}}*{{end}}{{.StructName}}) Set(key {{.KeyType}}, value {{.ValueType}}) {
{{if .ByReference}}	{{.ReceiverName}}.mx.Lock()
	defer {{.ReceiverName}}.mx.Unlock()

{{end}}{{if .LazyInstantiates}}	{{.ReceiverName}}.onceToken.Do(func() {
		{{.ReceiverName}}.{{.InternalMapName}} = make(map[{{.KeyType}}]{{.ValueType}})
	})
{{end}}	{{.ReceiverName}}.{{.InternalMapName}}[key] = value
}
`

// AbsorbTmpl wraps the `Absorb` func
const AbsorbTmpl = `// Absorb will take all the keys and values from another {{.StructName}}'s internal map and
// overwrite any existing keys
func ({{.ReceiverName}} {{if .ByReference}}*{{end}}{{.StructName}}) Absorb(otherMap {{if .ByReference}}*{{end}}{{.StructName}}) {
{{if .ByReference}}	{{.ReceiverName}}.mx.Lock()
	otherMap.mx.RLock()
	defer otherMap.mx.RUnlock()
	defer {{.ReceiverName}}.mx.Unlock()

{{end}}{{if .LazyInstantiates}}	{{.ReceiverName}}.onceToken.Do(func() {
		{{.ReceiverName}}.{{.InternalMapName}} = make(map[{{.KeyType}}]{{.ValueType}})
	})
{{end}}	for k, v := range otherMap.{{.InternalMapName}} {
		{{.ReceiverName}}.{{.InternalMapName}}[k] = v
	}
}
`

// AbsorbMapTmpl wraps the `AbsorbMap` func
const AbsorbMapTmpl = `// AbsorbMap will take all the keys and values from another map and overwrite any existing keys
func ({{.ReceiverName}} {{if .ByReference}}*{{end}}{{.StructName}}) AbsorbMap(regularMap map[{{.KeyType}}]{{.ValueType}}) {
{{if .ByReference}}	{{.ReceiverName}}.mx.Lock()
    defer {{.ReceiverName}}.mx.Unlock()

{{end}}{{if .LazyInstantiates}}	{{.ReceiverName}}.onceToken.Do(func() {
		{{.ReceiverName}}.{{.InternalMapName}} = make(map[{{.KeyType}}]{{.ValueType}})
	})
{{end}}	for k, v := range regularMap {
		{{.ReceiverName}}.{{.InternalMapName}}[k] = v
	}
}
`

// DeleteTmpl wraps the `Delete` func
const DeleteTmpl = `// Delete will remove a {{.ValueType}} from the map by key
func ({{.ReceiverName}} {{if .ByReference}}*{{end}}{{.StructName}}) Delete(key {{.KeyType}}) {
{{if .ByReference}}	{{.ReceiverName}}.mx.Lock()
	defer {{.ReceiverName}}.mx.Unlock()

{{end}}{{if .LazyInstantiates}}    {{.ReceiverName}}.onceToken.Do(func() {
    	{{.ReceiverName}}.{{.InternalMapName}} = make(map[{{.KeyType}}]{{.ValueType}})
	})
{{end}}	delete({{.ReceiverName}}.{{.InternalMapName}}, key)
}
`

// ClearTmpl wraps the `Clear` func
const ClearTmpl = `// Clear will remove all elements from the map
func ({{.ReceiverName}} {{if .ByReference}}*{{end}}{{.StructName}}) Clear() {
{{if .ByReference}}	{{.ReceiverName}}.mx.Lock()
	defer {{.ReceiverName}}.mx.Unlock()

{{end}}	{{.ReceiverName}}.{{.InternalMapName}} = make(map[{{.KeyType}}]{{.ValueType}})
}
`
