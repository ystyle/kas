package zlib

import (
	"bytes"
	"compress/zlib"
	"io"
)

func Decode(buff []byte) ([]byte, error) {
	b := bytes.NewReader(buff)
	var out bytes.Buffer
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(&out, r)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func Encode(buff []byte) ([]byte, error) {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	_, err := w.Write(buff)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return in.Bytes(), nil
}
