// Code generated by /tools/tmpl.go. DO NOT EDIT.
package tmpl

/*
If you want to edit the templates directly, you can do so in
under tmpl/_gotemplate
The templates are maintained there for greater readability.
*/

// MapTmpl contains all the necessary data to generate a carto map
type MapTmpl struct {
	GenDate               string
	PackageName           string
	StructName            string
	KeyTypePackage        string
	KeyType               string
	ValueType             string
	ValueTypePackage      string
	ReceiverName          string
	GetReturnsBool        bool
	LazyInstantiates      bool
	GetDefault            bool
}

// HeadTmpl is the file header, including imports and struct declaration.
// If lazy map instantiation is _not_ enabled, this also wraps the New... func.
const HeadTmpl = `{{.HeadTmpl}}`

// GetTmpl wraps the `Get` func
const GetTmpl = `{{.GetTmpl}}`

// KeysTmpl wraps the `Keys` func
const KeysTmpl = `{{.KeysTmpl}}`

// SetTmpl wraps the `Set` func
const SetTmpl = `{{.SetTmpl}}`

// AbsorbTmpl wraps the `Absorb` func
const AbsorbTmpl = `{{.AbsorbTmpl}}`

// AbsorbMapTmpl wraps the `AbsorbMap` func
const AbsorbMapTmpl = `{{.AbsorbMapTmpl}}`

// DeleteTmpl wraps the `Delete` func
const DeleteTmpl = `{{.DeleteTmpl}}`

// ClearTmpl wraps the `Clear` func
const ClearTmpl = `{{.ClearTmpl}}`
