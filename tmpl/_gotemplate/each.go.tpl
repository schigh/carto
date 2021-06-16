// Each runs a function over each key/value pair in the {{.StructName}}
// If the function returns false, the interation through the underlying
// map will halt.
// This function does not mutate the underlying map, although the values
// of the map may be mutated in place
// 		!!! Warning: calls to any mutating functions of {{.StructName}}
//		!!! will deadlock if called from within the supplied function
func ({{.ReceiverName}} *{{.StructName}}) Each(f func(key {{.KeyType}}, value {{.ValueType}}) bool) {
	defer {{.ReceiverName}}.mx.Unlock()
	{{.ReceiverName}}.mx.Lock()

	for _k := range {{.ReceiverName}}.impl {
		_v := {{.ReceiverName}}.impl[_k]
		if !f(_k, _v) {
			return
		}
	}
}
