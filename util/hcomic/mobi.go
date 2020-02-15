package hcomic

import (
	"os"
	"os/exec"
	"path"
	"runtime"
)

func Run(dir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ConverToMobi(opf string, bookname string) error {
	command := "kindlegen"
	if runtime.GOOS == "windows" {
		command = "kindlegen.exe"
	}
	kindlegen, err := exec.LookPath(command)
	if err != nil {
		return err
	}
	err = Run(path.Dir(opf), kindlegen, "-c1", "-dont_append_source", path.Base(opf), "-o", bookname)
	return err
}
