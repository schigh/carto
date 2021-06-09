// Delete will remove a {{.ValueType}} from the map by key
func ({{.ReceiverName}} *{{.StructName}}) Delete(key {{.KeyType}}) {
defer {{.ReceiverName}}.mx.Unlock()
	{{.ReceiverName}}.mx.Lock()

{{if .LazyInstantiates}}\\\
    {{.ReceiverName}}.onceToken.Do(func() {
    	{{.ReceiverName}}.impl = make(map[{{.KeyType}}]{{.ValueType}})
	})
{{end}}\\\
	delete({{.ReceiverName}}.impl, key)
}
