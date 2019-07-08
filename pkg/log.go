package dbjumper

import (
	"fmt"
	"os"
)

// Log prints to console
func Log(format string, str ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf(format, str...)
	}
}

// Error prints errors
func Error(err error) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf("%e", err)

	}
}
