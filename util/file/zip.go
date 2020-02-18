package file

import (
	"archive/zip"
	"bytes"
	"github.com/ystyle/kas/util/config"
	"io/ioutil"
	"path"
)

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

func CompressZipToFile(source, zipfiename string) error {
	buff, err := CompressZip(source)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(zipfiename, buff, config.Perm)
	return err
}
