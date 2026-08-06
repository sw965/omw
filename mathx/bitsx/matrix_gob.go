package bitsx

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type gobEncodedMatrix struct {
	Rows int
	Cols int
	Data []uint64
}

func (m *Matrix) GobEncode() ([]byte, error) {
	buf := &bytes.Buffer{}
	payload := gobEncodedMatrix{Rows: m.rows, Cols: m.cols, Data: m.data}
	if err := gob.NewEncoder(buf).Encode(payload); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Matrix) GobDecode(b []byte) error {
	var payload gobEncodedMatrix
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&payload); err != nil {
		return err
	}

	decoded := &Matrix{rows: payload.Rows, cols: payload.Cols, data: payload.Data}
	if err := decoded.validateDotAVX512Family(); err != nil {
		return fmt.Errorf("デコードされたMatrixが不正: %w", err)
	}

	// 端数ビットが常に0になるように、強制する。
	decoded.ApplyTailMask()
	*m = *decoded
	return nil
}
