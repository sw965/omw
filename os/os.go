package os

import (
	"os"
)

type DirEntries []os.DirEntry

func NewDirEntries(path string) (DirEntries, error) {
	return os.ReadDir(path)
}

func (es DirEntries) Names() []string {
	y := make([]string, len(es))
	for i, e := range es {
		y[i] = e.Name()
	}
	return y
}
