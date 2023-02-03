package omw

import (
	"encoding/json"
	"io/ioutil"
	"bytes"
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

func LoadJson[T any](v *T, path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	if err := json.Unmarshal(file, v); err != nil {
		return err
	}
	return nil
}

func WriteJson[T any](v *T, path string) error {
	file, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, file, 0644)
}
