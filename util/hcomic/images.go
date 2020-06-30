package hcomic

import (
	"fmt"
	"github.com/ystyle/kas/model"
	"github.com/ystyle/kas/util/web"
	"net/url"
	"path"
	"regexp"
	"strings"
)

func GetAllImages(book *model.HcomicInfo) error {
	html, err := web.GetHtmlNode(book.Url)
	if err != nil {
		return err
	}
	reg, _ := regexp.Compile(`this.src='(.*)';`)
	meta := html.Find("meta[name=\"applicable-device\"]")
	if attr, has := meta.Attr("content"); has && attr == "pc,mobile" {
		if book.BookName == "" {
			book.BookName = html.Find("#info-block #info h1").Text()
		}
		imgs := html.Find(".container .thumb-container .gallerythumb img")
		for i := range imgs.Nodes {
			img := imgs.Eq(i)
			src, _ := img.Attr("data-src")
			jstext, _ := img.Attr("onerror")
			group := reg.FindStringSubmatch(jstext)
			backurl := ""
			if len(group) > 0 {
				backurl = group[1]
			}
			book.AddSection(fmt.Sprintf("#%d", i+1), GetHDImage(src), GetHDImage(backurl))
		}
	} else {
		lis := html.Find(".img_list li")
		if book.BookName == "" {
			book.BookName = html.Find(".page_tit .tit").Text()
		}
		for i := range lis.Nodes {
			li := lis.Eq(i)
			url, _ := li.Find("img").First().Attr("src")
			title := li.Find("label").Text()
			if url != "" {
				book.AddSection(title, GetHDImage(url), "")
			}
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
		return strings.ReplaceAll(url, "pic.comicstatic.icu", "img.comicstatic.icu")
	} else if strings.Contains(url, "mt.404cdn.com") {
		u := strings.ReplaceAll(url, "mt.404cdn.com", "mi.404cdn.com")
		return strings.ReplaceAll(u, "t.", ".")
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
