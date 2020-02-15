package zip

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"os"
	"path"
)

func IsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CompressZip(filename string) ([]byte, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// Create a buffer to write our archive to.
	var buff bytes.Buffer
	// Create a new zip archive.
	w := zip.NewWriter(&buff)
	// Add some files to the archive.
	f, err := w.Create(path.Base(filename))
	if err != nil {
		return nil, err
	}
	_, err = f.Write(bs)
	if err != nil {
		return nil, err
	}
	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
