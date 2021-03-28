// Code generated {{.GenDate}} by carto.  DO NOT EDIT.
package {{.PackageName}}
import (
	"sync"

{{if .KeyTypePackage}}\\\
	"{{.KeyTypePackage}}"
{{end}}\\\
{{if .ValueTypePackage}}\\\
	"{{.ValueTypePackage}}"
{{end}}\\\
)

// {{.StructName}} wraps map[{{.KeyType}}]{{.ValueType}}, and locks reads and writes with a mutex
type {{.StructName}} struct {
	mx sync.RWMutex
	impl map[{{.KeyType}}]{{.ValueType}}
{{if .LazyInstantiates}}\\\
	onceToken sync.Once
{{end}}
}
{{if .LazyInstantiates}}{{else}}\\\

// New{{.StructName}} generates a new {{.StructName}} with a non-nil map
func New{{.StructName}}() *{{.StructName}} {
	{{.ReceiverName}} := &{{.StructName}}{}
	{{.ReceiverName}}.impl = make(map[{{.KeyType}}]{{.ValueType}})

	return {{.ReceiverName}}
}
{{end}}\\\
