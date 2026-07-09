package gobx

import (
	"bytes"
	"encoding/gob"
	"os"
)

func Load[T any](path string) (T, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		var zero T
		return zero, err
	}

	buf := bytes.NewBuffer(file)
	dec := gob.NewDecoder(buf)

	var data T
	if err := dec.Decode(&data); err != nil {
		var zero T
		return zero, err
	}
	return data, nil
}

func Save[T any](data T, path string) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0644)
}
