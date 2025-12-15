// Package gobx provides generic helper functions for Gob file I/O.
// It acts as a wrapper around encoding/gob, simplifying the encoding and decoding processes.
//
// Gob形式のファイルの読み書きを行うためのジェネリックなヘルパー関数を提供します。
// encoding/gob のラッパーとして機能し、エンコードとデコードの処理を簡略化します。
package gobx

import (
	"bytes"
	"encoding/gob"
	"os"
)

// Load reads a Gob-encoded file from the specified path and decodes it into type T.
// Note that type T must match the structure of the data that was saved.
//
// 指定されたパスからGobエンコードされたファイルを読み込み、型 T にデコードします。
// 型 T は保存されたデータの構造と一致している必要があります。
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

// Save encodes the provided data into Gob format and writes it to the specified path.
// Only exported fields will be encoded. The file is created with 0644 permissions.
//
// 渡されたデータをGob形式にエンコードし、指定されたパスに保存します。
// エクスポートされた（大文字で始まる）フィールドのみがエンコードされます。ファイルはパーミッション 0644 で作成されます。
func Save[T any](data T, path string) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0644)
}
