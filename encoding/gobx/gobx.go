package gobx

import (
	"bytes"
	"encoding/gob"
	"os"
)

func Load[T any](path string) (T, error) {
	var result T

	data, err := os.ReadFile(path)
	if err != nil {
		return result, err
	}

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}

func Save[T any](data T, path string) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0644)
}
