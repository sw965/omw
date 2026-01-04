package jsonx

import (
	"bytes"
	"encoding/json"
	"os"
)

func Load[T any](path string) (T, error) {
	var result T
	file, err := os.ReadFile(path)
	if err != nil {
		return result, err
	}

	// BOM除去
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))

	if err := json.Unmarshal(file, &result); err != nil {
		return result, err
	}
	return result, nil
}

func Save[T any](data T, path string) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, file, 0644)
}
