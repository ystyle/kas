package kindlegen

import (
	"github.com/ystyle/kas/util/file"
	"os"
	"os/exec"
	"path"
	"runtime"
)

func Run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Conver(source string, bookname string) error {
	command := "kindlegen"
	if runtime.GOOS == "windows" {
		command = "kindlegen.exe"
	}
	kindlegen, err := exec.LookPath(command)
	if err != nil {
		return err
	}
	err = Run(kindlegen, "-c1", "-dont_append_source", source, "-o", bookname)
	if err != nil {
		if ok, _ := file.IsExists(path.Join(path.Dir(source), bookname)); ok {
			return nil
		}
		return err
	}
	return nil
}
