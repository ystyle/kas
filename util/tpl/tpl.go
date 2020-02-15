package tpl

import (
	"log"
	"os"
	"text/template"
)

var tpl  *template.Template

func init() {
	t, err :=  template.New("page").Parse(PAGE)
	if err != nil {
		log.Fatal(err)
	}
	t, err =  t.New("opf").Parse(OPF)
	if err != nil {
		log.Fatal(err)
	}
	t, err =  t.New("toc").Parse(TOC)
	if err != nil {
		log.Fatal(err)
	}
	tpl = t
}

func Render(filename, code string, data interface{})  error {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	err = tpl.ExecuteTemplate(f, code, data)
	return err
}