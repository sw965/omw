package os

import (
	"os"
)

type DirEntries []os.DirEntry

func NewDirEntries(path string) (DirEntries, error) {
	dirs, err := os.ReadDir(path)
	return dirs, err
}

func (d DirEntries) Names() []string {
	y := make([]string, len(d))
	for i, dir := range d {
		y[i] = dir.Name()
	}
	return y
}
