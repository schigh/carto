// Set will add an element to the {{.StructName}}'s internal map with the specified key
func ({{.ReceiverName}} *{{.StructName}}) Set(key {{.KeyType}}, value {{.ValueType}}) {
	{{.ReceiverName}}.mx.Lock()
	defer {{.ReceiverName}}.mx.Unlock()

{{if .LazyInstantiates}}\\\
	{{.ReceiverName}}.onceToken.Do(func() {
		{{.ReceiverName}}.impl = make(map[{{.KeyType}}]{{.ValueType}})
	})
{{end}}\\\
	{{.ReceiverName}}.impl[key] = value
}
