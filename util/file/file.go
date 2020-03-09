package file

import (
	"fmt"
	"github.com/ystyle/kas/util/config"
	"os"
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

func CheckDir(dir string) {
	if ok, _ := IsExists(dir); !ok {
		err := os.MkdirAll(dir, config.Perm)
		fmt.Println(err)
	}
}
