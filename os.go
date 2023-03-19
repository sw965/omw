package omw

import (
	"io/ioutil"
)

func ListDir(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return []string{}, err
	}
	y := make([]string, len(files))
	for i, file := range files {
		y[i] = file.Name()
	}
	return y, nil
}
