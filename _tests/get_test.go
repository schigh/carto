package cartotests

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	type ts struct {
		name     string
		testArgs []string
		internal map[string]interface{}
		key      string
		dflt     interface{}
		rBool    bool
		expect   interface{}
	}
	tests := []ts{
		{
			name:   "base with nil map",
			key:    "foo",
			expect: nil,
		},
		{
			name: "base with value",
			internal: map[string]interface{}{
				"foo": 42,
			},
			key:    "foo",
			expect: 42,
		},
		{
			name:     "lazy with nil map",
			testArgs: []string{"LZ"},
			key:      "foo",
			expect:   nil,
		},
		{
			name:     "lazy with value",
			testArgs: []string{"LZ"},
			internal: map[string]interface{}{
				"foo": 42,
			},
			key:    "foo",
			expect: 42,
		},
		{
			name:     "D with nil map",
			testArgs: []string{"D"},
			dflt:     "bar",
			key:      "foo",
			expect:   "bar",
		},
		{
			name:     "D with value",
			testArgs: []string{"D"},
			dflt:     "bar",
			internal: map[string]interface{}{
				"foo": 42,
			},
			key:    "foo",
			expect: 42,
		},
		{
			name:     "B with nil map",
			testArgs: []string{"B"},
			key:      "foo",
			expect:   nil,
			rBool:    false,
		},
		{
			name:     "B with value",
			testArgs: []string{"B"},
			internal: map[string]interface{}{
				"foo": 42,
			},
			key:    "foo",
			expect: 42,
			rBool:  true,
		},
		{
			name:     "LZD with nil map",
			testArgs: []string{"LZ", "D"},
			dflt:     "bar",
			key:      "foo",
			expect:   "bar",
		},
		{
			name:     "LZD with value",
			testArgs: []string{"LZ", "D"},
			dflt:     "bar",
			internal: map[string]interface{}{
				"foo": 42,
			},
			key:    "foo",
			expect: 42,
		},
		{
			name:     "LZB with nil map",
			testArgs: []string{"LZ", "B"},
			key:      "foo",
			expect:   nil,
			rBool:    false,
		},
		{
			name:     "LZB with value",
			testArgs: []string{"LZ", "B"},
			internal: map[string]interface{}{
				"foo": 42,
			},
			key:    "foo",
			expect: 42,
			rBool:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := newCartoTest(t, tt.internal, tt.testArgs...)
			var got interface{}
			var exists bool
			switch {
			case test.retB:
				got, exists = test.gbImpl.Get(tt.key)
			case test.def:
				got = test.gdImpl.Get(tt.key, tt.dflt)
			default:
				got = test.gImpl.Get(tt.key)
			}

			if exists != tt.rBool {
				t.Fatalf("unexpected existence flag returned: %t", exists)
			}
			if !reflect.DeepEqual(got, tt.expect) {
				t.Fatalf("Get() => expected %v, got %v", tt.expect, got)
			}

			tt2 := tt
			for i := 0; i < 20; i++ {
				go func(test *cartoTest, tt *ts) {
					switch {
					case test.retB:
						_, _ = test.gbImpl.Get(tt.key)
					case test.def:
						_ = test.gdImpl.Get(tt.key, tt.dflt)
					default:
						_ = test.gImpl.Get(tt.key)
					}
				}(test, &tt2)
			}
		})
	}
}
