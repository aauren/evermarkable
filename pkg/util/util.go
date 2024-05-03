package util

import (
	"bytes"
	"encoding/json"
	"io"
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
