package cartotests

import "testing"

type base interface {
	Keys() []string
	Set(string, interface{})
	AbsorbMap(map[string]interface{})
	Delete(string)
	Clear()
}
type getter interface {
	base
	Get(string) interface{}
}
type getterDefault interface {
	base
	Get(string, interface{}) interface{}
}
type getterBool interface {
	base
	Get(string) (interface{}, bool)
}

type cartoTest struct {
	baseImpl base
	gImpl    getter
	gdImpl   getterDefault
	gbImpl   getterBool
	lazy     bool
	def      bool
	retB     bool
}

const (
	isBase int = 1 << iota
	isLZ
	isD
	isB
)

func newCartoTest(t *testing.T, data map[string]interface{}, attrs ...string) *cartoTest {
	t.Helper()
	var tMask int
	test := &cartoTest{}
	for _, a := range attrs {
		switch a {
		case "LZ":
			test.lazy = true
			tMask |= isLZ
		case "D":
			test.def = true
			tMask |= isD
		case "B":
			test.retB = true
			tMask |= isB
		}
	}
	if tMask == 0 {
		tMask = isBase
	}
	switch tMask {
	case isBase:
		impl := &Base{impl: data}
		test.baseImpl = impl
		test.gImpl = impl
	case isLZ:
		impl := &Base0LZ{impl: data}
		test.baseImpl = impl
		test.gImpl = impl
	case isD:
		impl := &Base0D{impl: data}
		test.baseImpl = impl
		test.gdImpl = impl
	case isB:
		impl := &Base0B{impl: data}
		test.baseImpl = impl
		test.gbImpl = impl
	case isLZ | isD:
		impl := &Base0LZ0D{impl: data}
		test.baseImpl = impl
		test.gdImpl = impl
	case isLZ | isB:
		impl := &Base0LZ0B{impl: data}
		test.baseImpl = impl
		test.gbImpl = impl
	default:
		t.Fatal("invalid combination")
	}

	return test
}
