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
