// Package meta provides an assortment of coding-related functions.
package meta

import (
	"fmt"
	"os"
)

// RecoverFromPanicAt wraps a function to recover upon its panic.
// https://golang.org/ref/spec#Handling_panics
func RecoverFromPanicAt(g func()) {
	defer func() {
		recover()
	}()
	g()
}

// ----------------------------------------------------------------------------
// DEBUG

type PrintDebug func(...interface{})

// NewDebug returns a fmt.Println(...interface{}) that prints only when in debug mode (on=true).
func NewDebug(on bool) PrintDebug {
	return func(msgs ...interface{}) {
		if !on {
			return
		}
		fmt.Println(msgs)
	}
}

// ----------------------------------------------------------------------------
// FILESYSTEM

// IsDir ...
func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	mode := fi.Mode()
	return mode.IsDir()
}

// IsRegFile ...
func IsRegFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	mode := fi.Mode()
	return mode.IsRegular()
}
