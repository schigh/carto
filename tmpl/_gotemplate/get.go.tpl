{{if .GetDefault}}\\\
// Get gets the {{.ValueType}} keyed by {{.KeyType}}.  If the key does not exist, a default {{.ValueType}} will be returned
func ({{.ReceiverName}} *{{.StructName}}) Get(key {{.KeyType}}, dflt {{.ValueType}})(value {{.ValueType}}) {
defer {{.ReceiverName}}.mx.RUnlock()
    {{.ReceiverName}}.mx.RLock()

	var ok bool
	value, ok = {{.ReceiverName}}.impl[key]
	if !ok {
		value = dflt
	}

	return
}
{{else}}\\\
// Get gets the {{.ValueType}} keyed by {{.KeyType}}. {{if .GetReturnsBool}}Also returns bool value indicating whether the key exists in the map{{end}}
func ({{.ReceiverName}} *{{.StructName}}) Get(key {{.KeyType}}) {{if .GetReturnsBool}}(value {{.ValueType}}, ok bool){{else}}(value {{.ValueType}}){{end}} {
	defer {{.ReceiverName}}.mx.RUnlock()
    {{.ReceiverName}}.mx.RLock()

	value{{if .GetReturnsBool}}, ok{{end}} = {{.ReceiverName}}.impl[key]

	return
}
{{end}}\\\
