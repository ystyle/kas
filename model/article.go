package model

type ArticleInfo struct {
	Title   string   // 书籍的标题
	UrlList []string //
}

type ArticleItem struct {
	Title   string // 标题
	Url     string // 网页地址
	Content string // 提取的内容
	images  map[string]string
}

type SiteInfo struct {
	Name    string   // 网站名称
	Url     string   // 地址规则
	Title   string   // 标题匹配规则
	Desc    string   // 详情匹配规则
	Include string   // 正文匹配规则
	Exclude []string // 过滤规则
}
