package util

import (
	"os"
	"path/filepath"
)

func ExePath(s ...string) string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exeDir := filepath.Dir(exePath)
	return filepath.Join(append([]string{exeDir}, s...)...)
}
