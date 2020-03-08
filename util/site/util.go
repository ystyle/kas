package site

import (
	"encoding/json"
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/gommon/log"
	"github.com/ystyle/kas/model"
	"github.com/ystyle/kas/util/array"
	"github.com/ystyle/kas/util/config"
	"github.com/ystyle/kas/util/file"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

var (
	sites        []model.SiteInfo
	website_list = "website_list.json"
)

func Init(box *rice.Box) {
	var buff []byte
	localeFile := path.Join(config.StoreDir, website_list)
	if has, _ := file.IsExists(localeFile); has {
		buff, _ = ioutil.ReadFile(localeFile)
	} else {
		bs, err := box.Bytes("asset/website_list.json")
		if err != nil {
			log.Error(err)
			return
		}
		ioutil.WriteFile(localeFile, bs, config.Perm)
		buff = bs
	}

	err := json.Unmarshal(buff, &sites)
	if err != nil {
		log.Error(err)
	}
}

func getConfig(url string) *model.SiteInfo {
	for _, site := range sites {
		reg, err := regexp.Compile(site.Url)
		if err != nil {
			log.Error(err)
			continue
		}
		if reg.MatchString(url) {
			return &site
		}
	}
	return nil
}

func ParseContent(node *goquery.Selection, item *model.ArticleItem) error {
	// 清除不必要的代码和样式
	node.Find("script").Remove()
	node.Find("style").Remove()
	node.Find("[class]").RemoveClass()
	node.Find("[style]").RemoveAttr("style")
	node.Find("div,p").SetAttr("style", "text-indent: 2rem")
	node.Find("button").Remove()
	// 取网站配置
	site := getConfig(item.Url)
	if site != nil {
		content := getHtml(node, site)
		excludeHtmlContent(content, site)
		collectImages(content, item)
		html, err := content.Html()
		if err != nil {
			return err
		}
		html = excludeTextContent(html, site)
		item.Content = clearUnspportTag(html)
	} else {
		collectImages(node, item)
		html, err := node.Html()
		if err != nil {
			return err
		}
		item.Content = clearUnspportTag(html)
	}
	return nil
}

func clearUnspportTag(html string) string {
	html = strings.ReplaceAll(html, "<li", "<div")
	html = strings.ReplaceAll(html, "</li>", "</div>")
	//html = strings.ReplaceAll(html, "style", "data-style")
	return html
}

func collectImages(node *goquery.Selection, item *model.ArticleItem) {
	item.Images = make(map[string]string)
	imgs := node.Find("img")
	supportImageTypes := []string{".git", ".jpg", ".jpeg", ".png"}
	for i := range imgs.Nodes {
		img := imgs.Eq(i)
		src, _ := img.Attr("src")
		if src == "" {
			src, _ = img.Attr("data-src")
		}
		// 取不到图片地址时跳过
		if src == "" {
			img.Remove()
			continue
		}
		filename := path.Base(src)
		if !array.IncludesFromString(supportImageTypes, filename) {
			img.Remove()
			continue
		}
		if strings.HasPrefix(src, "//") {
			if strings.HasPrefix(item.Url, "https") {
				src = "https:" + src
			} else {
				src = "http:" + src
			}
		}
		key := fmt.Sprintf(".img%d", i)
		item.Images[key] = src
		img.ReplaceWithHtml(fmt.Sprintf("{{ %s }}", key))
	}
}

func excludeTextContent(html string, site *model.SiteInfo) string {
	for _, s := range site.Exclude {
		if strings.HasPrefix(s, "text:") {
			text := strings.ReplaceAll(s, "text:", "")
			html = strings.ReplaceAll(html, text, "")
		}
	}
	return html
}

func excludeHtmlContent(node *goquery.Selection, site *model.SiteInfo) {
	for _, s := range site.Exclude {
		if strings.HasPrefix(s, "selector:") {
			selector := strings.ReplaceAll(s, "selector:", "")
			node.Find(selector).Remove()
		}
	}
}

func getHtml(node *goquery.Selection, site *model.SiteInfo) *goquery.Selection {
	if strings.HasPrefix(site.Include, "selector:") {
		selector := strings.ReplaceAll(site.Include, "selector:", "")
		return node.Find(selector)
	}
	return node
}
