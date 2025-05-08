package os

import (
	"bufio"
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

func ReadLines(path string, c int) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0, c)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
