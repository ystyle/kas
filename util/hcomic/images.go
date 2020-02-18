package hcomic

import (
	"github.com/ystyle/kas/model"
	"github.com/ystyle/kas/util/web"
	"net/url"
	"path"
	"strings"
)

func GetAllImages(book *model.HcomicInfo) error {
	html, err := web.GetHtmlNode(book.Url)
	if err != nil {
		return err
	}
	lis := html.Find(".img_list li")
	if book.BookName == "" {
		book.BookName = html.Find(".page_tit .tit").Text()
	}
	for i := range lis.Nodes {
		li := lis.Eq(i)
		url, _ := li.Find("img").First().Attr("src")
		title := li.Find("label").Text()
		if url != "" {
			book.AddSection(title, GetHDImage(url))
		}
	}
	return nil
}

func GetHDImage(url string) string {
	// 预览图
	// https://pic.comicstatic.icu/img/cn/1570141/1.jpg
	// 高清图
	// https://img.comicstatic.icu/img/cn/1570141/1.jpg
	if strings.Contains(url, "pic.") {
		return strings.ReplaceAll(url, "pic.", "img.")
	}
	// 没有匹配到则用预览图
	return url
}

func GetComicID(page string) (string, error) {
	u, err := url.Parse(page)
	if err != nil {
		return "", err
	}
	return path.Base(u.Path), nil
}
