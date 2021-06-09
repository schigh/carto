// Value returns a copy of the underlying map[{{.KeyType}}]{{.ValueType}}
func ({{.ReceiverName}} *{{.StructName}}) Value() map[{{.KeyType}}]{{.ValueType}} {
	{{.ReceiverName}}.mx.RLock()
	defer {{.ReceiverName}}.mx.RUnlock()

	out := make(map[{{.KeyType}}]{{.ValueType}}, len({{.ReceiverName}}.impl))
	for k, v := range {{.ReceiverName}}.impl {
		out[k] = v
	}

	return out
}
