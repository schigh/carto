// Value returns a copy of the underlying map[{{.KeyType}}]{{.ValueType}}
func ({{.ReceiverName}} *{{.StructName}}) Value() map[{{.KeyType}}]{{.ValueType}} {
	defer {{.ReceiverName}}.mx.RUnlock()
{{.ReceiverName}}.mx.RLock()

	out := make(map[{{.KeyType}}]{{.ValueType}}, len({{.ReceiverName}}.impl))
	for k, v := range {{.ReceiverName}}.impl {
		out[k] = v
	}

	return out
}
