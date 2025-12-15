// Package jsonx provides generic helper functions for JSON file I/O.
// It acts as a wrapper around encoding/json, adding support for
// automatically stripping UTF-8 BOMs during generic Load operations.
//
// JSONファイルの読み書きを行うためのジェネリックなヘルパー関数を提供します。
// encoding/json のラッパーとして機能し、Load 時の UTF-8 BOM 自動除去などをサポートします。
package jsonx

import (
	"bytes"
	"encoding/json"
	"os"
)

// Load reads a JSON file from the specified path and unmarshals it into type T.
// It automatically strips the UTF-8 Byte Order Mark (BOM) if present.
//
// 指定されたパスからJSONファイルを読み込み、型 T に変換（アンマーシャル）します。
// ファイルの先頭に UTF-8 BOM (Byte Order Mark) が存在する場合は、自動的に除去します。
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

// Save marshals the provided data into JSON format with indentation and writes it to the specified path.
// The file is created with 0644 permissions.
//
// 渡されたデータをインデント付きのJSON形式に変換（マーシャル）し、指定されたパスに保存します。
// ファイルはパーミッション 0644 で作成されます。
func Save[T any](data T, path string) error {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, file, 0644)
}
