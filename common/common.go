package common

import (
	"log"
	"os"
)

// CheckError log a Fatal if error not nil
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// CloseFile allows to close an *os.File
func CloseFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

// OpenFile opens a path and returns an *os.File
func OpenFile(path string) *os.File {
	f, err := os.Open(path)
	CheckError(err)
	return f
}
