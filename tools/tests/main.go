package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var testStructs = map[string]string{
	"Base":      `-p cartotests -k string -v interface{} -s Base -o %s/_tests/base.go`,
	"Base0D":    `-p cartotests -k string -v interface{} -s Base0D -d -o %s/_tests/base_d.go`,
	"Base0B":    `-p cartotests -k string -v interface{} -s Base0B -b -o %s/_tests/base_b.go`,
	"Base0LZ":   `-p cartotests -k string -v interface{} -s Base0LZ -lz -o %s/_tests/base_lz.go`,
	"Base0LZ0D": `-p cartotests -k string -v interface{} -s Base0LZ0D -lz -d -o %s/_tests/base_lz_d.go`,
	"Base0LZ0B": `-p cartotests -k string -v interface{} -s Base0LZ0B -lz -b -o %s/_tests/base_lz_b.go`,
}

func createCommands(pwd string) map[string]*exec.Cmd {
	commands := make(map[string]*exec.Cmd)
	exe := fmt.Sprintf("%s/carto", pwd)
	for k, v := range testStructs {
		c := fmt.Sprintf(v, pwd)
		args := strings.Split(c, " ")
		cmd := exec.Command(exe, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		commands[k] = cmd
	}

	return commands
}

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	commands := createCommands(pwd)
	for _, cmd := range commands {
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
