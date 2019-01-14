package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var (
	packageName           string
	structName            string
	keyType               string
	valueType             string
	byValue               bool
	receiverName          string
	getReturnsBool        bool
	lazyInstantiates      bool
	outFileName           string
	internalMapName       string
	noMutex               bool
	isGoGenerate          bool
	keyTypePackage        string
	keyTypePackageShort   string
	keyTypeIsPointer      bool
	valueTypePackage      string
	valueTypePackageShort string
	valueTypeIsPointer    bool
	getDefault            bool

	reserved = []string{
		"i", "k", "v", "keys", "onceToken", "value", "ok", "otherMap",
	}
)

func init() {
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
	flag.BoolVar(&noMutex, "xm", false, "")
	flag.BoolVar(&getDefault, "d", false, "")
}

func usage() {
	printBold("C A R T O\n")
	printInfo("Maps made easy")
	usg := `
usage:
-p      package name !
-s      struct name  !
-k      key type     !
-v      value type   !
-r      receiver name (defaults to lowercase first char of struct name)
-rv     receivers are by value
-b      "Get" return signature includes a bool value indicating if the key exists in the internal map
-d      "Get" signature has second parameter for default return value when key does not exist in the internal map
-lz     will lazy-instantiate the internal map when a write operation is used
-o      output file name (if omitted, prints to STDOUT)
-i      variable name for internal map (defaults to internal)
-xm     operations will not be mutexed
`
	usg = strings.Replace(usg, "!", "\033[33m(required)\033[0m", -1)
	printPlain(usg)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var errors []string

	// check package
	if packageName == "" {
		packageName = os.Getenv("GOPACKAGE")
		isGoGenerate = true
	}
	if packageName == "" {
		errors = append(errors, "   - no package specified")
	}

	// check struct
	if structName == "" {
		errors = append(errors, "   - no struct specified")
	}

	// check key type
	if keyType == "" {
		errors = append(errors, "   - no key type specified")
	} else {
		var keyTypeErr error
		keyTypePackage, keyType, keyTypeErr = parsePackage(keyType)
		if keyTypeErr != nil {
			errors = append(errors, keyTypeErr.Error())
		}
	}

	// check value type
	if valueType == "" {
		errors = append(errors, "   - no value type specified")
	} else {
		var valueTypeErr error
		valueTypePackage, valueType, valueTypeErr = parsePackage(valueType)
		if valueTypeErr != nil {
			errors = append(errors, valueTypeErr.Error())
		}
	}

	if len(errors) > 0 {
		printErr("unable to generate CARTO struct:\n" + strings.Join(errors, "\n"))
		os.Exit(1)
	}

	// default receiver name
	if receiverName == "" {
		receiverName = strings.ToLower(string(structName[0]))
	}
	for _, r := range reserved {
		if receiverName == r {
			receiverName = "_" + receiverName
			break
		}
	}

	mt := &mapTmpl{
		GenDate:            time.Now().Format(time.RFC1123),
		PackageName:        filepath.Base(packageName),
		StructName:         structName,
		Sync:               lazyInstantiates || !noMutex,
		Mutex:              !noMutex,
		KeyType:            keyType,
		KeyTypePackage:     keyTypePackage,
		KeyTypeIsPointer:   keyTypeIsPointer,
		ValueType:          valueType,
		ValueTypePackage:   valueTypePackage,
		ValueTypeIsPointer: valueTypeIsPointer,
		InternalMapName:    internalMapName,
		ByReference:        !byValue,
		ReceiverName:       receiverName,
		GetReturnsBool:     getReturnsBool,
		LazyInstantiates:   lazyInstantiates,
		GetDefault:         getDefault,
	}

	templates := []string{
		headTmpl,
		getTmpl,
		keysTmpl,
		setTmpl,
		absorbTmpl,
		absorbMapTmpl,
		deleteTmpl,
		clearTmpl,
	}
	var buf []byte
	b := bytes.NewBuffer(buf)

	for i, tmpl := range templates {
		t, err := template.New(fmt.Sprintf("tmpl_%d", i)).Parse(tmpl)
		if err != nil {
			printErr("template error: %s", err.Error())
			os.Exit(1)
		}

		if err := t.Execute(b, mt); err != nil {
			printErr("template execute error: %s", err.Error())
			os.Exit(1)
		}
	}

	formatted, err := format.Source(b.Bytes())
	if err != nil {
		printErr("formatting error: %s", err.Error())
		os.Exit(1)
	}

	printSuccess("struct created")
	printPlain(string(formatted))
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
