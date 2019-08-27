package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"text/template"
	"time"
)

func ensurePath() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if path.Base(pwd) != "carto" {
		return "", errors.New("this generator must be run from the carto project root")
	}

	return pwd, nil
}

type tmplTemplate struct {
	GenDate       string
	HeadTmpl      string
	GetTmpl       string
	KeysTmpl      string
	SetTmpl       string
	AbsorbTmpl    string
	AbsorbMapTmpl string
	DeleteTmpl    string
	ClearTmpl     string
}

func getFileStr(pwd, filename string) string {
	re := regexp.MustCompile(`(?m)\\\\\\\n`)
	filename = path.Join(pwd, "tmpl", "_gotemplate", filename) + ".gotemplate"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return string(re.ReplaceAll(data, []byte{}))
}

func main() {
	pwd, err := ensurePath()
	if err != nil {
		log.Fatal(err)
	}

	headBytes := getFileStr(pwd, "head")
	getBytes := getFileStr(pwd, "get")
	keysBytes := getFileStr(pwd, "keys")
	setBytes := getFileStr(pwd, "set")
	absorbBytes := getFileStr(pwd, "absorb")
	absorbMapBytes := getFileStr(pwd, "absorbmap")
	deleteBytes := getFileStr(pwd, "delete")
	clearBytes := getFileStr(pwd, "clear")
	tmplBytes := getFileStr(pwd, "tmpl")

	tmpl := tmplTemplate{
		GenDate:       time.Now().Format(time.RFC1123),
		HeadTmpl:      headBytes,
		GetTmpl:       getBytes,
		KeysTmpl:      keysBytes,
		SetTmpl:       setBytes,
		AbsorbTmpl:    absorbBytes,
		AbsorbMapTmpl: absorbMapBytes,
		DeleteTmpl:    deleteBytes,
		ClearTmpl:     clearBytes,
	}

	var b bytes.Buffer
	t, err := template.New("tmpl_main").Parse(tmplBytes)
	if err != nil {
		log.Fatal(err)
	}
	if err := t.Execute(&b, &tmpl); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(path.Join(pwd, "tmpl", "tmpl.go"), b.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
}
