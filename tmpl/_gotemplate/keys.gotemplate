// Keys will return all keys in the {{.StructName}}'s internal map
func ({{.ReceiverName}} *{{.StructName}}) Keys() (keys []{{.KeyType}}) {
	{{.ReceiverName}}.mx.RLock()
	defer {{.ReceiverName}}.mx.RUnlock()

	keys = make([]{{.KeyType}}, len({{.ReceiverName}}.impl))
	var i int
	for k := range {{.ReceiverName}}.impl {
		keys[i] = k
		i++
	}

	return
}
