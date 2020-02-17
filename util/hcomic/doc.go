package hcomic

import (
	"github.com/ystyle/kas/model"
	"github.com/ystyle/kas/util/tpl"
	"io"
	"os"
)

func GenDoc(book model.HcomicInfo) error {
	for i, section := range book.Sections {
		err := tpl.Render(section.HtmlFile, "page", section)
		if err != nil {
			return err
		}
		if i == 0 {
			// 生成封面
			source, err := os.Open(section.ImgFile)
			if err != nil {
				return err
			}
			destination, err := os.Create(book.CoverFile)
			if err != nil {
				return err
			}
			io.Copy(destination, source)
			source.Close()
			destination.Close()
			// 在第一个文件生成bom
			file, _ := os.Open(section.HtmlFile)
			file.WriteString("\xEF\xBB\xBF")
			file.Close()
		}
	}

	err := tpl.Render(book.OpfFile, "opf", book)
	if err != nil {
		return err
	}
	err = tpl.Render(book.NxcFile, "toc", book)
	if err != nil {
		return err
	}
	return nil
}
