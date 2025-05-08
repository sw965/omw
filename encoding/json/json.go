package json

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

const (
	EXTENSION = ".json"
)

func Load[T any](path string) (T, error) {
	var result T
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return result, err
	}
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	if err := json.Unmarshal(file, &result); err != nil {
		return result, err
	}
	return result, nil
}

func Save[T any](data *T, path string) error {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, file, 0644)
	return err
}
