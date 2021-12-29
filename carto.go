package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/schigh/carto/io"
	"github.com/schigh/carto/tmpl"
)

var allTemplates = []string{
	tmpl.HeadTmpl,
	tmpl.GetTmpl,
	tmpl.KeysTmpl,
	tmpl.SetTmpl,
	tmpl.AbsorbTmpl,
	tmpl.AbsorbMapTmpl,
	tmpl.DeleteTmpl,
	tmpl.ClearTmpl,
	tmpl.ValueTmpl,
	tmpl.SizeTmpl,
	tmpl.EachTmpl,
}

var (
	packageName      string
	structName       string
	keyType          string
	valueType        string
	receiverName     string
	getReturnsBool   bool
	lazyInstantiates bool
	outFileName      string
	getDefault       bool

	version bool
	Version = "2021-12-29"

	reserved = []string{"i", "k", "v", "keys", "onceToken", "value", "ok", "otherMap", "mx", "impl"}

	keyRx          = regexp.MustCompile(`^(?:(?P<pkg>(?:[\w.-]+/)*\w+)\.)?(?P<type>\w+)$`)
	defaultValueRx = regexp.MustCompile(`^(?P<slc>\[(?P<sz>\d*)])?(?P<ptr>\*)?(?:(?P<pkg>(?:[\w.-]+/)*\w+)\.)?(?P<type>(\w+|interface{}))$`)
	mapValueRx     = regexp.MustCompile(`^map\[(?:(?P<map_pkg>(?:[\w.-]+/)*\w+)\.)?(?P<map_type>\w+)](?P<slc>\[(?P<sz>\d*)])?(?P<ptr>\*)?(?:(?P<pkg>(?:[\w.-]+/)*\w+)\.)?(?P<type>(\w+|interface{}))$`)
)

type pkgCtx byte

const (
	_ pkgCtx = iota
	keyCtx
	valueCtx
)

func (c pkgCtx) valid() bool {
	switch c {
	case keyCtx, valueCtx:
		return true
	}

	return false
}

func usage() {
	io.PrintBold("C A R T O\n\n")
	io.PrintInfo("🌐 Maps made easy")
	usg := `
usage:
carto -p <package> -s <structname> -k <keytype> -v <valuetype> [options]

-p    (string)      Package name !
-s    (string)      Struct name  !
-k    (string)      Key type     !
-v    (string)      Value type   !
-r    (string)      Receiver name (defaults to lowercase first char of 
                    struct name)
-o    (string)      Output file path (if omitted, prints to STDOUT)
-b    (bool)        "Get" return signature includes a bool value indicating 
                    if the key exists in the internal map
-d    (bool)        "Get" signature has second parameter for default return 
                    value when key does not exist in the internal map
-lz   (bool)        Will lazy-instantiate the internal map when a write 
                    operation is used
-version            Print version and exit
`
	usg = strings.Replace(usg, "!", "\033[33m(required)\033[0m", -1)
	io.PrintPlain(usg)
}

func main() {
	if handleFlags() {
		os.Exit(0)
	}

	var errList []string
	packageName, structName, keyType, valueType, errList = ensureRequired(packageName, structName, keyType, valueType)
	if len(errList) > 0 {
		err := errors.New(strings.Join(errList, "\n"))
		handleError(err)
	}

	var keyTypePackage, valueTypePackage string
	var err error

	keyTypePackage, keyType, err = parsePackage(keyType, keyCtx)
	if err != nil {
		handleError(err)
	}

	valueTypePackage, valueType, err = parsePackage(valueType, valueCtx)
	if err != nil {
		handleError(err)
	}

	receiverName = reservedKwds(defaultReceiver(receiverName, structName), reserved)

	mt := tmpl.MapTmpl{
		GenDate:          time.Now().Format(time.RFC1123),
		PackageName:      filepath.Base(packageName),
		StructName:       structName,
		KeyType:          keyType,
		KeyTypePackage:   keyTypePackage,
		ValueType:        valueType,
		ValueTypePackage: valueTypePackage,
		ReceiverName:     receiverName,
		GetReturnsBool:   getReturnsBool,
		LazyInstantiates: lazyInstantiates,
		GetDefault:       getDefault,
	}

	var b *bytes.Buffer
	b, err = parseAndExecTemplates(&mt, allTemplates)
	if err != nil {
		handleError(err)
	}

	var formatted []byte
	formatted, err = applyFormatting(b.Bytes())
	if err != nil {
		handleError(err)
	}

	var fileCreated bool
	fileCreated, err = createOutFile(outFileName, formatted)
	if err != nil {
		handleError(err)
	}
	if !fileCreated {
		io.PrintSuccess("struct created")
		io.PrintPlain(string(formatted))
	}
}

func handleError(err error) {
	io.PrintErr("CARTO generate struct failed.\n%s", err.Error())
	os.Exit(1)
}

func handleFlags() bool {
	flag.StringVar(&packageName, "p", "", "")
	flag.StringVar(&structName, "s", "", "")
	flag.StringVar(&keyType, "k", "", "")
	flag.StringVar(&valueType, "v", "", "")
	flag.StringVar(&receiverName, "r", "", "")
	flag.BoolVar(&getReturnsBool, "b", false, "")
	flag.BoolVar(&lazyInstantiates, "lz", false, "")
	flag.StringVar(&outFileName, "o", "", "")
	flag.BoolVar(&getDefault, "d", false, "")
	flag.BoolVar(&version, "version", false, "")
	flag.Usage = usage
	flag.Parse()

	if version {
		io.PrintInfo("CARTO version: %s", Version)
		return true
	}

	return false
}

func ensureRequired(p, s, k, v string) (pkgVal, structVal, keyVal, valVal string, errList []string) {

	// check package
	if p == "" {
		errList = append(errList, "   - package is required ('-p')")
	}

	// check struct
	if s == "" {
		errList = append(errList, "   - struct name is required ('-s')")
	}

	// check key type
	if k == "" {
		errList = append(errList, "   - key type is required ('-k')")
	}

	// check value type
	if v == "" {
		errList = append(errList, "   - value type is required ('-v')")
	}

	pkgVal = p
	structVal = s
	keyVal = k
	valVal = v

	return
}

func defaultReceiver(r, s string) string {
	// default receiver name
	if r == "" {
		r = strings.ToLower(string(s[0]))
	}

	return r
}

func reservedKwds(r string, kwds []string) string {
	// prefix any "reserved" keywords with an underscore
	for _, rw := range kwds {
		if r == rw {
			r = "_" + r
			break
		}
	}
	return r
}

func parseAndExecTemplates(sat interface{}, tmpls []string) (*bytes.Buffer, error) {
	var buf []byte
	b := bytes.NewBuffer(buf)

	for i, tplt := range tmpls {
		t, err := template.New(fmt.Sprintf("tmpl_%d", i)).Parse(tplt)
		if err != nil {
			return nil, err
		}

		if err := t.Execute(b, sat); err != nil {
			return nil, err
		}
	}

	return b, nil
}

func applyFormatting(b []byte) ([]byte, error) {
	data, err := format.Source(b)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func createOutFile(fn string, data []byte) (bool, error) {
	if fn != "" {
		if err := ioutil.WriteFile(fn, data, 0644); err != nil {
			return false, err
		}
		io.PrintSuccess("struct created")
		return true, nil
	}
	return false, nil
}

func parsePackage(ppath string, pCtx pkgCtx) (packageName string, typeName string, err error) {
	if !pCtx.valid() {
		err = errors.New("invalid context for package parsing")
	}

	if ppath == "" {
		err = errors.New("type or package declaration was empty")
		return
	}

	if pCtx == keyCtx {
		packageName, typeName, err = parseKeyPackage(ppath)
		return
	}

	var m map[string]string
	if defaultValueRx.MatchString(ppath) {
		m = rxNamedExtract(defaultValueRx, ppath)
	} else if mapValueRx.MatchString(ppath) {
		m = rxNamedExtract(mapValueRx, ppath)
		typeName = typeName + "map["
		if m["map_pkg"] != "" {
			typeName = typeName + filepath.Base(m["map_pkg"]) + "."
		}
		typeName = typeName + m["map_type"] + "]"
	} else {
		err = fmt.Errorf("the package path '%s' could not be parsed", ppath)
		return
	}

	if m["slc"] != "" {
		typeName = typeName + "["
		if m["sz"] != "" {
			typeName = typeName + m["sz"]
		}
		typeName = typeName + "]"
	}

	if m["ptr"] != "" {
		typeName = typeName + "*"
	}

	if m["pkg"] != "" {
		packageName = m["pkg"]
		typeName = typeName + filepath.Base(packageName) + "." + m["type"]
		return
	}

	typeName = typeName + m["type"]

	return
}

func parseKeyPackage(ppath string) (packageName string, typeName string, err error) {
	if !keyRx.MatchString(ppath) {
		err = fmt.Errorf("path '%s' is an invalid key type", ppath)
		return
	}
	m := rxNamedExtract(keyRx, ppath)
	packageName = m["pkg"]
	typeName = m["type"]
	return
}

// this extracts regex values only for named captures...
// unnamed captures are skipped
func rxNamedExtract(rx *regexp.Regexp, s string) map[string]string {
	// the item at index zero is always the full match without captures
	names := rx.SubexpNames()
	if len(names) < 1 {
		return nil
	}
	names = names[1:]
	sm := rx.FindStringSubmatch(s)
	if len(sm) < 1 {
		return nil
	}
	sm = sm[1:]

	// assert that length of names == length of sm
	// see stdlib regexp implementation for details

	out := make(map[string]string)

	for i := range names {
		if names[i] == "" {
			// unnamed capture group
			continue
		}
		out[names[i]] = sm[i]
	}

	return out
}
