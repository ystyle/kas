package test

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ystyle/kas/util/config"
	"github.com/ystyle/kas/util/web"
	"net/url"
	"path"
	"regexp"
	"testing"
)

func TestUrl(t *testing.T) {
	u, _ := url.Parse("http://www.qdaily.com/articles/65038.html")
	fmt.Println(u.Hostname())
}

func TestReg(t *testing.T) {
	text := `$('.view_content').html()||$('#view_content').html()||$('#articleContent').html()||$('.forum-viewthread-article-box').html()`
	reg9, _ := regexp.Compile(`\$\('([#\w\d\s\.-]*)'\)\.(text|html)\(\)`)
	if reg9.MatchString(text) {
		group := reg9.FindAllStringSubmatch(text, -1)
		for _, items := range group {
			for _, item := range items {
				fmt.Println(item)
			}
		}
	}
}

func TestGoquery(t *testing.T) {
	var buff bytes.Buffer
	buff.WriteString("")
	node, err := goquery.NewDocumentFromReader(&buff)
	if err != nil {
		t.Error(err)
	}
	div := node.Find(".article-detail-bd")
	div.RemoveFiltered(".author-share")
	div.Find("script").Remove()
	div.Find("style").Remove()
	imgs := div.Find("img")
	images := map[string]string{}
	for i := range imgs.Nodes {
		img := imgs.Eq(i)
		src, _ := img.Attr("src")
		if src == "" {
			src, _ = img.Attr("data-src")
		}
		key := fmt.Sprintf(".img%d", i)
		images[key] = src
		img.ReplaceWithHtml(fmt.Sprintf("{{ %s }}", key))
	}
	fmt.Println(images)
	//div := node.Find("<div class='summary'>")
	fmt.Println(div.Html())
}

func TestHcomic(t *testing.T) {
	html, err := web.GetHtmlNode("https://bhmog.com/s/70622/")
	if err != nil {
		fmt.Println(err)
		return
	}
	meta := html.Find("meta[name=\"applicable-device\"]")
	if attr, has := meta.Attr("content"); has && attr == "pc,mobile" {
		imgs := html.Find(".container .gallery img")
		for i := range imgs.Nodes {
			img := imgs.Eq(i)
			src, _ := img.Attr("data-src")
			fmt.Println(src)
		}
	}
	fmt.Println()
	fmt.Println(html.Html())
}

func TestDownload(t *testing.T) {
	dir := path.Join(config.CacheDir, "19.jpg")
	fmt.Println(dir)
	// https://aa.hcomics.club/uploads/1615779/2.jpg
	err := web.Download("https://aa.hcomics.club/uploads/1615779/3.jpg", "2.jpg")
	if err != nil {
		fmt.Println(err)
	}
}
