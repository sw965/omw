package jsonx

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/sw965/omw/encoding/atomicfile"
)

func Load[T any](path string) (T, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		var zero T
		return zero, err
	}

	// BOM除去
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))

	var data T
	if err := json.Unmarshal(file, &data); err != nil {
		var zero T
		return zero, err
	}
	return data, nil
}

func Save[T any](data T, path string) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return atomicfile.WriteFile(path, file, 0o644)
}
