package main

import (
	"os"
	"path/filepath"
)

func getAppCurrentDir() string {
	absPath, _ := filepath.Abs(os.Args[0])

	return filepath.Dir(absPath)
}
