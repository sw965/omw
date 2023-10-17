package os

import (
	"os"
	"github.com/sw965/omw"
)

type DirEntries []os.DirEntry

func NewDirEntries(path string) (DirEntries, error) {
	dirs, err := os.ReadDir(path)
	return dirs, err
}

func (ds DirEntries) Names() []string {
	y := make([]string, len(ds))
	for i, d := range ds {
		y[i] = d.Name()
	}
	return y
}