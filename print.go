package main

import (
	"fmt"
	"os"
)

//  print bold
func printBold(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, "\033[34m\033[1m"+msg+"\033[0m", args...)
}

//  print plain
func printPlain(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, "\033[30m"+msg+"\033[0m", args...)
}

// printErr will print in red
func printErr(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "\033[31m💥 "+msg+"\033[0m\n", args...)
}

// printSuccess will print in green
func printSuccess(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, "\033[32m👍 "+msg+"\033[0m\n", args...)
}

// printInfo will print in cyan
func printInfo(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, "\033[36m🌐 "+msg+"\033[0m\n", args...)
}
