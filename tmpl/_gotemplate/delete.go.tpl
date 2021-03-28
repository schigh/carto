// Delete will remove a {{.ValueType}} from the map by key
func ({{.ReceiverName}} *{{.StructName}}) Delete(key {{.KeyType}}) {
	{{.ReceiverName}}.mx.Lock()
	defer {{.ReceiverName}}.mx.Unlock()

{{if .LazyInstantiates}}\\\
    {{.ReceiverName}}.onceToken.Do(func() {
    	{{.ReceiverName}}.impl = make(map[{{.KeyType}}]{{.ValueType}})
	})
{{end}}\\\
	delete({{.ReceiverName}}.impl, key)
}
