// Size returns the number of elements in the underlying map[{{.KeyType}}]{{.ValueType}}
func ({{.ReceiverName}} *{{.StructName}}) Size() int {
	defer {{.ReceiverName}}.mx.RUnlock()
	{{.ReceiverName}}.mx.RLock()

	return len({{.ReceiverName}}.impl)
}
