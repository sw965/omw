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
	var ret T
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return ret, err
	}
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	if err := json.Unmarshal(file, &ret); err != nil {
		return ret, err
	}
	return ret, nil
}

func WriteJSON[T any](data *T, path string) error {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, file, 0644)
	return err
}