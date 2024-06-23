package util

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

func ToIOReader(source interface{}) (io.Reader, error) {
	var content []byte
	var err error

	if source == nil {
		return bytes.NewReader(nil), nil
	}

	content, err = json.Marshal(source)

	return bytes.NewReader(content), err
}

func CopyFile(src, dst string) (int64, error) {
	r, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	w, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer w.Close()

	n, err := io.Copy(w, r)
	if err != nil {
		return 0, err
	}

	return n, nil
}
