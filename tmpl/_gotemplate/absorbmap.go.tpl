// AbsorbMap will take all the keys and values from another map and overwrite any existing keys
func ({{.ReceiverName}} *{{.StructName}}) AbsorbMap(regularMap map[{{.KeyType}}]{{.ValueType}}) {
defer {{.ReceiverName}}.mx.Unlock()
	{{.ReceiverName}}.mx.Lock()

{{if .LazyInstantiates}}\\\
	{{.ReceiverName}}.onceToken.Do(func() {
		{{.ReceiverName}}.impl = make(map[{{.KeyType}}]{{.ValueType}})
	})
{{end}}\\\
	for k, v := range regularMap {
		{{.ReceiverName}}.impl[k] = v
	}
}
