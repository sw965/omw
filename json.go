package omw

import (
	"encoding/json"
	"io/ioutil"
	"bytes"
)

const (
	JSON_EXTENSION = ".json"
)

func LoadJSON[T any](path string) (T, error) {
	var y T
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return y, err
	}
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	if err := json.Unmarshal(file, &y); err != nil {
		return y, err
	}
	return y, nil
}

func WriteJSON[T any](y *T, path string) error {
	file, err := json.MarshalIndent(y, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, file, 0644)
	return err
}