// Keys will return all keys in the {{.StructName}}'s internal map
func ({{.ReceiverName}} *{{.StructName}}) Keys() (keys []{{.KeyType}}) {
defer {{.ReceiverName}}.mx.RUnlock()
	{{.ReceiverName}}.mx.RLock()

	keys = make([]{{.KeyType}}, len({{.ReceiverName}}.impl))
	var i int
	for k := range {{.ReceiverName}}.impl {
		keys[i] = k
		i++
	}

	return
}
