// Clear will remove all elements from the map
func ({{.ReceiverName}} *{{.StructName}}) Clear() {
	{{.ReceiverName}}.mx.Lock()
	defer {{.ReceiverName}}.mx.Unlock()

	{{.ReceiverName}}.impl = make(map[{{.KeyType}}]{{.ValueType}})
}
