package kindlegen

import (
	"os"
	"os/exec"
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
	return err
}
