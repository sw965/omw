package json

import (
	"encoding/json"
	"io/ioutil"
	"bytes"
)

func Load[T any](path string) (T, error) {
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

func Write[T any](path string) (T, error) {
	var y T
	file, err := json.MarshalIndent(&y, "", " ")
	if err != nil {
		return y, err
	}
	err = ioutil.WriteFile(path, file, 0644)
	return y, err
}