package hcomic

import (
	"github.com/ystyle/hcc/util/web"
	"net/url"
	"path"
	"strings"
)

func GetAllImages(url string) ([]string, string, error) {
	html, err := web.GetContent(url)
	if err != nil {
		return nil, "", err
	}
	imgs := html.Find(".img_list li img")
	title := html.Find(".page_tit .tit").Text()
	var images []string
	for i := range imgs.Nodes {
		img := imgs.Eq(i)
		url, _ := img.Attr("src")
		images = append(images, url)
	}
	return images, title, nil
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
