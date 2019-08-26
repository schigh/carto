package io

import (
	"fmt"
	"os"
)

//  PrintBold will print bold text
func PrintBold(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "\033[34m\033[1m"+msg+"\033[0m", args...)
}

//  PrintPlain will print plain text
func PrintPlain(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, msg, args...)
}

// PrintErr will print error text
func PrintErr(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "\033[31m💥 "+msg+"\033[0m\n", args...)
}

// PrintSuccess will print success text
func PrintSuccess(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "\033[32m👍 "+msg+"\033[0m\n", args...)
}

// PrintInfo will print info text
func PrintInfo(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "\033[36m"+msg+"\033[0m\n", args...)
}
