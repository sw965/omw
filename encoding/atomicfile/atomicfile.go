package atomicfile

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func WriteFile(path string, data []byte, perm fs.FileMode) error {
	return WriteFrom(path, bytes.NewReader(data), perm)
}

func WriteFrom(path string, r io.Reader, perm fs.FileMode) (err error) {
	if r == nil {
		return fmt.Errorf("io.Readerがnil")
	}

	// 指定された親ディレクトリのパスを取得する
	dir := filepath.Dir(path)

	// 指定されたパスの最後の要素を取得する
	fileName := filepath.Base(path)

	// Unix系では隠しファイルとして扱われるよう、ファイル名の先頭に「.」を付ける
	// Windowsでは隠し属性は付かないが、Unix系と共通の命名慣習に従い、先頭に「.」を付ける
	// 末尾の「*」はランダムな文字列に置き換えられる
	tmpFileName := "." + fileName + ".tmp-*"

	// 指定されたパスと同じディレクトリに一時ファイルを作る
	tmp, createErr := os.CreateTemp(dir, tmpFileName)
	if createErr != nil {
		return createErr
	}

	tmpPath := tmp.Name()
	committed := false

	defer func() {
		// 処理が最後まで完了している場合、何もしない
		if committed {
			return
		}
		// 途中で失敗した場合の後始末

		// 一時ファイルの削除を試みる前に、開いたままなら閉じて、ファイルハンドルを解放する
		if closeErr := tmp.Close(); closeErr != nil && !errors.Is(closeErr, os.ErrClosed) {
			err = errors.Join(err, closeErr)
		}

		// 一時ファイルを削除する
		if removeErr := os.Remove(tmpPath); removeErr != nil && !os.IsNotExist(removeErr) {
			err = errors.Join(err, removeErr)
		}
	}()

	// 一時ファイルに対して、指定された権限を設定する
	if err := tmp.Chmod(perm); err != nil {
		return err
	}

	// io.Readerから読み取った内容を、一時ファイルへ書き込む
	if _, err := io.Copy(tmp, r); err != nil {
		return err
	}

	// 一時ファイルの内容を、永続ストレージへ同期する
	if err := tmp.Sync(); err != nil {
		return err
	}

	// 一時ファイルを閉じ、名前変更に備えてファイルハンドルを解放する
	if err := tmp.Close(); err != nil {
		return err
	}

	// 一時ファイルの名前を、指定された名前に変更し、保存を完了する
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}

	committed = true
	return nil
}
