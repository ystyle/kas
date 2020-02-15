package hcomic

import (
	"fmt"
	"github.com/ystyle/hcc/util/tpl"
	"io"
	"os"
	"path"
)
import "github.com/satori/go.uuid"

func GenDoc(title, author, url string, images []string) error {
	comicId, _ := GetComicID(url)
	workdir := path.Join("cache", comicId)
	var imgList []map[string]string

	for i, image := range images {
		imgFilename := path.Base(image)
		item := map[string]string{
			"i":        fmt.Sprintf("%d", i+1),
			"id":       fmt.Sprintf("%d", i),
			"item":       fmt.Sprintf("%d", i+2),
			"title":    fmt.Sprintf("%d", i),
			"filename": fmt.Sprintf("html/Page-%d.html", i),
			"image":    fmt.Sprintf("scaled-images/%s", imgFilename),
		}
		htmlfile := path.Join(workdir, "html", fmt.Sprintf("Page-%d.html", i))
		err := tpl.Render(htmlfile, "page", item)
		if err != nil {
			return err
		}
		if i == 0 {
			// 生成封面
			coverimage := path.Join(workdir, "cover-image.jpg")
			first := path.Join(workdir, "html", "scaled-images", imgFilename)
			source, err := os.Open(first)
			if err != nil {
				return err
			}
			destination, err := os.Create(coverimage)
			if err != nil {
				return err
			}
			io.Copy(destination, source)
			source.Close()
			destination.Close()
			// 在第一个文件生成bom
			file, _ := os.Open(htmlfile)
			file.WriteString("\xEF\xBB\xBF")
			file.Close()
		}
		imgList = append(imgList, item)
	}

	data := map[string]interface{}{
		"id":     uuid.NewV4().String(),
		"title":  title,
		"author": author,
		"images": imgList,
	}

	err := tpl.Render(path.Join(workdir, "content.opf"), "opf", data)
	if err != nil {
		return err
	}
	err = tpl.Render(path.Join(workdir, "toc.ncx"), "toc", data)
	if err != nil {
		return err
	}
	return nil
}
