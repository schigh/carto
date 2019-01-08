package main

import (
	"bytes"
	"flag"
	"fmt"
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
}

func usage() {
	printBold("C A R T O")
	usg := `
usage:
-p      package name __
-s      struct name __
-k      key type __
-v      value type __
-r      receiver name (defaults to lowercase first char of struct name)
-rv     receivers are by value
-b      "Get" also returns a bool value indicating if the key exists in the internal map
-lz     will lazy-instantiate the internal map when a write operation is used
-o      output file name (if omitted, prints to STDOUT)
-i      variable name for internal map (defaults to internal)
-xm     operations will not be mutexed
`
	usg = strings.Replace(usg, "__", "\033[33m(required)\033[30m", -1)
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
		keyTypeIsPointer = keyType[0] == '*'
		if keyTypeIsPointer {
			keyType = keyType[1:]
		}
		keyTypePackage = path.Dir(keyType)
		if len(keyTypePackage) > 0 {
			keyTypePackageShort = path.Base(keyTypePackage) + "."
		}
		keyType = path.Base(keyType)
	}

	// check value type
	if valueType == "" {
		errors = append(errors, "   - no value type specified")
	} else {
		valueTypeIsPointer = valueType[0] == '*'
		if valueTypeIsPointer {
			valueType = valueType[1:]
		}
		valueTypePackage = path.Dir(valueType)
		if len(valueTypePackage) > 0 {
			valueTypePackageShort = path.Base(valueTypePackage) + "."
		}
		valueType = path.Base(valueType)
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

	//formatted, err := format.Source(b.Bytes())
	//if err != nil {
	//	printErr("formatting error: %s", err.Error())
	//	os.Exit(1)
	//}

	printSuccess("struct created")
	//printPlain(string(formatted))
	printPlain(string(b.Bytes()))
}
