// Absorb will take all the keys and values from another {{.StructName}}'s internal map and
// overwrite any existing keys
func ({{.ReceiverName}} *{{.StructName}}) Absorb(otherMap *{{.StructName}}) {
	defer otherMap.mx.RUnlock()
	defer {{.ReceiverName}}.mx.Unlock()
{{.ReceiverName}}.mx.Lock()
otherMap.mx.RLock()

{{if .LazyInstantiates}}\\\
	{{.ReceiverName}}.onceToken.Do(func() {
		{{.ReceiverName}}.impl = make(map[{{.KeyType}}]{{.ValueType}})
	})
{{end}}\\\
	for k, v := range otherMap.impl {
		{{.ReceiverName}}.impl[k] = v
	}
}
