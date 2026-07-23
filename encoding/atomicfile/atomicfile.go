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

// WriteFile writes data to a temporary file in the destination directory and
// replaces the destination only after the complete contents reach the file.
func WriteFile(path string, data []byte, perm fs.FileMode) error {
	return WriteFrom(path, bytes.NewReader(data), perm)
}

// WriteFrom copies all of r into a temporary file in the destination directory and
// replaces the destination only after the complete contents reach the file.
// Unlike WriteFile, it streams the content instead of holding it all in memory at once.
func WriteFrom(path string, r io.Reader, perm fs.FileMode) (err error) {
	if r == nil {
		return fmt.Errorf("rがnilです")
	}

	dir := filepath.Dir(path)
	tmp, createErr := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp-*")
	if createErr != nil {
		return createErr
	}
	tmpPath := tmp.Name()
	committed := false
	defer func() {
		if committed {
			return
		}
		if closeErr := tmp.Close(); closeErr != nil && !errors.Is(closeErr, os.ErrClosed) {
			err = errors.Join(err, closeErr)
		}
		if removeErr := os.Remove(tmpPath); removeErr != nil && !os.IsNotExist(removeErr) {
			err = errors.Join(err, removeErr)
		}
	}()

	if err := tmp.Chmod(perm); err != nil {
		return err
	}
	if _, err := io.Copy(tmp, r); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}
	committed = true
	return nil
}
