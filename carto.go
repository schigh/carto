package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
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
}

var (
	packageName      string
	structName       string
	keyType          string
	valueType        string
	byValue          bool
	receiverName     string
	getReturnsBool   bool
	lazyInstantiates bool
	outFileName      string
	internalMapName  string
	keyTypePackage   string
	valueTypePackage string
	getDefault       bool

	version bool
	Version string

	reserved = []string{"i", "k", "v", "keys", "onceToken", "value", "ok", "otherMap", "mx"}
)

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
-i    (string)      Variable name for internal map (defaults to internal)
-rv   (bool)        Receivers are by value
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

	packageName, structName, keyType, valueType, errList := ensureRequired(packageName, structName, keyType, valueType)
	if len(errList) > 0 {
		err := errors.New(strings.Join(errList, "\n"))
		handleError(err)
	}

	keyTypePackage, keyType, err := parsePackage(keyType)
	if err != nil {
		handleError(err)
	}

	valueTypePackage, valueType, err := parsePackage(valueType)
	if err != nil {
		handleError(err)
	}

	receiverName = reservedKwds(defaultReceiver(receiverName, structName), reserved)
	internalMapName = reservedKwds(internalMapName, reserved)

	mt := tmpl.MapTmpl{
		GenDate:          time.Now().Format(time.RFC1123),
		PackageName:      filepath.Base(packageName),
		StructName:       structName,
		KeyType:          keyType,
		KeyTypePackage:   keyTypePackage,
		ValueType:        valueType,
		ValueTypePackage: valueTypePackage,
		InternalMapName:  internalMapName,
		ByReference:      !byValue,
		ReceiverName:     receiverName,
		GetReturnsBool:   getReturnsBool,
		LazyInstantiates: lazyInstantiates,
		GetDefault:       getDefault,
	}

	b, err := parseAndExecTemplates(&mt, allTemplates)
	if err != nil {
		handleError(err)
	}
	formatted, err := applyFormatting(b.Bytes())
	if err != nil {
		handleError(err)
	}

	fileCreated, err := createOutFile(outFileName, formatted)
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
	flag.BoolVar(&byValue, "rv", false, "")
	flag.BoolVar(&getReturnsBool, "b", false, "")
	flag.BoolVar(&lazyInstantiates, "lz", false, "")
	flag.StringVar(&outFileName, "o", "", "")
	flag.StringVar(&internalMapName, "i", "internal", "")
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

func ensureRequired(p, s, k, v string) (pVal, sVal, kVal, vVal string, errList []string) {

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

	pVal = p
	sVal = s
	kVal = k
	vVal = v

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

func parsePackage(ppath string) (packageName string, typeName string, err error) {
	if ppath == "" {
		err = errors.New("type or package declaration was empty")
		return
	}
	isPointerType := ppath[0] == '*'
	if isPointerType {
		ppath = ppath[1:]
	}

	pathParts := strings.Split(ppath, ".")
	numParts := len(pathParts)

	if numParts == 1 {
		typeName = pathParts[0]
		return
	}

	// two parts
	packageName = strings.Join(pathParts[:numParts-1], ".")
	typeName = path.Base(packageName) + "." + pathParts[numParts-1]
	if isPointerType {
		typeName = "*" + typeName
	}

	return
}
