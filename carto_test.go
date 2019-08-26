package main

import (
	"os"
	"path"
	"reflect"
	"testing"
)

func Test_handleFlags(t *testing.T) {

}

func Test_ensureRequired(t *testing.T) {
	tests := []struct {
		name    string
		pIn     string
		sIn     string
		kIn     string
		vIn     string
		errList []string
	}{
		{
			name:    "happy path",
			pIn:     "foo",
			sIn:     "foo",
			kIn:     "foo",
			vIn:     "foo",
			errList: []string(nil),
		},
		{
			name:    "p missing",
			sIn:     "foo",
			kIn:     "foo",
			vIn:     "foo",
			errList: []string{"   - package is required ('-p')"},
		},
		{
			name:    "s missing",
			pIn:     "foo",
			kIn:     "foo",
			vIn:     "foo",
			errList: []string{"   - struct name is required ('-s')"},
		},
		{
			name:    "k missing",
			pIn:     "foo",
			sIn:     "foo",
			vIn:     "foo",
			errList: []string{"   - key type is required ('-k')"},
		},
		{
			name:    "v missing",
			pIn:     "foo",
			sIn:     "foo",
			kIn:     "foo",
			errList: []string{"   - value type is required ('-v')"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, errList := ensureRequired(tt.pIn, tt.sIn, tt.kIn, tt.vIn)
			if !reflect.DeepEqual(errList, tt.errList) {
				t.Fatalf("ensureRequired() => expected: %v, got %v", tt.errList, errList)
			}
		})
	}
}

func Test_defaultReceiver(t *testing.T) {
	tests := []struct {
		name   string
		r      string
		s      string
		expect string
	}{
		{
			name:   "not set for Foo",
			r:      "",
			s:      "Foo",
			expect: "f",
		},
		{
			name:   "set for Foo",
			r:      "ff",
			s:      "Foo",
			expect: "ff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := defaultReceiver(tt.r, tt.s)
			if got != tt.expect {
				t.Fatalf("defaultReceiver() => expected %s, got %s", tt.expect, got)
			}
		})
	}
}

func Test_reservedKwds(t *testing.T) {
	allKwds := []string{"i", "k", "v", "keys", "onceToken", "value", "ok", "otherMap", "mx"}
	tests := []struct {
		name   string
		r      string
		kwds   []string
		expect string
	}{
		{
			name:   "no conflict",
			r:      "foo",
			kwds:   allKwds,
			expect: "foo",
		},
		{
			name:   "i",
			r:      "i",
			kwds:   allKwds,
			expect: "_i",
		},
		{
			name:   "ok",
			r:      "ok",
			kwds:   allKwds,
			expect: "_ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reservedKwds(tt.r, tt.kwds)
			if got != tt.expect {
				t.Fatalf("reservedKwds() => expected %s, got %s", tt.expect, got)
			}
		})
	}
}

func Test_parseAndExecTemplates(t *testing.T) {
	type s struct {
		Foo string
		Bar string
	}
	tests := []struct {
		name     string
		sat      *s
		tmpls    []string
		wantErr  bool
		expected []byte
	}{
		{
			name: "test",
			sat: &s{
				Foo: "test",
				Bar: "test",
			},
			tmpls:    []string{`{{.Foo}}-{{.Bar}}`},
			wantErr:  false,
			expected: []byte(`test-test`),
		},
		{
			name: "unknown property",
			sat: &s{
				Foo: "test",
				Bar: "test",
			},
			tmpls:    []string{`{{Foo}}-{{.Bar}}`},
			wantErr:  true,
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := parseAndExecTemplates(tt.sat, tt.tmpls)
			if err != nil && !tt.wantErr {
				t.Fatalf("parseAndExecTemplates{} => unexpected error: %v", err)
			}
			if tt.wantErr && err == nil {
				t.Fatal("parseAndExecTemplates() => error expected")
			}
			var bt []byte
			if b != nil {
				bt = b.Bytes()
			}
			if !reflect.DeepEqual(tt.expected, bt) {
				t.Fatalf("parseAndExecTemplates() => expected %s, got %s", tt.expected, bt)
			}
		})
	}
}

func Test_createOutFile(t *testing.T) {
	tests := []struct {
		name    string
		fn      string
		data    []byte
		b       bool
		wantErr bool
	}{
		{
			name:    "no outfile",
			fn:      "",
			data:    nil,
			b:       false,
			wantErr: false,
		},
		{
			name:    "with outfile",
			fn:      path.Join(os.TempDir(), "test.txt"),
			data:    []byte("this is test data"),
			b:       true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := createOutFile(tt.fn, tt.data)
			if err != nil && !tt.wantErr {
				t.Fatalf("createOutFile{} => unexpected error: %v", err)
			}
			if tt.wantErr && err == nil {
				t.Fatal("createOutFile() => error expected")
			}
			if b != tt.b {
				t.Fatalf("createOutFile() => unexpected value %t", b)
			}
		})
	}
}

func Test_parsePackage(t *testing.T) {
	tests := []struct {
		name    string
		ppath   string
		pName   string
		tName   string
		wantErr bool
	}{
		{
			name:    "interface{}",
			ppath:   "interface{}",
			pName:   "",
			tName:   "interface{}",
			wantErr: false,
		},
		{
			name:    "package path",
			ppath:   "github.com/schigh/carto.Foo",
			pName:   "github.com/schigh/carto",
			tName:   "carto.Foo",
			wantErr: false,
		},
		{
			name:    "package path, pointer",
			ppath:   "*github.com/schigh/carto.Foo",
			pName:   "github.com/schigh/carto",
			tName:   "*carto.Foo",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pn, tn, err := parsePackage(tt.ppath)
			if err != nil && !tt.wantErr {
				t.Fatalf("parsePackage{} => unexpected error: %v", err)
			}
			if tt.wantErr && err == nil {
				t.Fatal("parsePackage() => error expected")
			}
			if pn != tt.pName {
				t.Errorf("parsePackage() => wanted %s, got %s", tt.pName, pn)
			}
			if tn != tt.tName {
				t.Errorf("parsePackage() => wanted %s, got %s", tt.tName, tn)
			}
		})
	}
}
