// Clear will remove all elements from the map
func ({{.ReceiverName}} *{{.StructName}}) Clear() {
defer {{.ReceiverName}}.mx.Unlock()
	{{.ReceiverName}}.mx.Lock()

	{{.ReceiverName}}.impl = make(map[{{.KeyType}}]{{.ValueType}})
}
