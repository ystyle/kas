package model

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/ystyle/kas/util/config"
	"path"
)

type ArticleInfo struct {
	ID       string
	Title    string // 书籍的标题
	UrlList  []*ArticleItem
	EpubFile string
	MobiFile string
	ZipFile  string
}

func (article *ArticleInfo) SetDefault() {
	article.ID = uuid.NewV4().String()
	article.EpubFile = path.Join(config.CacheDir, "article", fmt.Sprintf("%s.epub", article.ID))
	article.MobiFile = path.Join(config.CacheDir, "article", fmt.Sprintf("%s.mobi", article.ID))
	article.ZipFile = path.Join(config.StoreDir, "article", fmt.Sprintf("%s.zip", article.ID))
}

type ArticleItem struct {
	Title   string // 标题
	Url     string // 网页地址
	Content string // 提取的内容
	Images  map[string]string
}

type SiteInfo struct {
	Name    string   // 网站名称
	Url     string   // 地址规则
	Title   string   // 标题匹配规则
	Desc    string   // 详情匹配规则
	Include string   // 正文匹配规则
	Exclude []string // 过滤规则
}
