package model

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/ystyle/kas/util/array"
	"github.com/ystyle/kas/util/config"
	"os"
	"path"
)

type TextInfo struct {
	ID        string     // id
	BookName  string     // 书名
	Author    string     // 作者
	Match     string     // 章节匹配规则
	Indent    uint       // 首先缩进
	Align     string     // 标题对齐方式
	MaxLen    int        // 标题最大字数
	OnlyKF8   string     // 只生成KF8格式: 1是，0否
	Content   []byte     // 文件内容
	Sections  []*Section // 章节内容
	Format    []string   // 格式
	CacheEpub string     // epub 缓存目录
	CacheMobi string     // mobi 缓存目录
	CacheCSS  string     // CSS保存目录
	StoreEpub string     // epub保存目录
	StoreMobi string     // mobi保存目录
}

type Section struct {
	Title   string // 标题
	Content string // 正文
}

func (text *TextInfo) SetDefault() {
	text.ID = uuid.NewV4().String()
	if text.OnlyKF8 == "" {
		text.OnlyKF8 = "1"
	}
	if text.Author == "" {
		text.Author = "KAF"
	}
	if text.Match == "" {
		text.Match = "^.{0,8}(第.{1,20}(章|节)|(S|s)ection.{1,20}|(C|c)hapter.{1,20}|(P|p)age.{1,20}|引子|楔子)"
	}
	if text.Indent == 0 {
		text.Indent = 2
	}
	if text.MaxLen == 0 {
		text.MaxLen = 35
	}
	if text.Format == nil {
		text.Format = append(text.Format, "mobi")
	}
	Aligns := []string{"center", "left", "right"}
	if !array.IncludesString(Aligns, text.Align) {
		text.Align = "center"
	}
	text.CacheCSS = path.Join(config.CacheDir, "text", fmt.Sprintf("%s.css", text.ID))
	text.CacheEpub = path.Join(config.CacheDir, "text", fmt.Sprintf("%s.epub", text.ID))
	text.CacheMobi = path.Join(config.CacheDir, "text", fmt.Sprintf("%s.mobi", text.ID))
	text.StoreEpub = path.Join(config.StoreDir, "text", fmt.Sprintf("[epub]%s.zip", text.ID))
	text.StoreMobi = path.Join(config.StoreDir, "text", fmt.Sprintf("[mobi]%s.zip", text.ID))
}

func (text *TextInfo) AddSection(title, content string) {
	text.Sections = append(text.Sections, &Section{
		Title:   title,
		Content: content,
	})
}

func (text *TextInfo) ClearCache() {
	os.Remove(text.CacheEpub)
	os.Remove(text.CacheMobi)
	os.Remove(text.CacheCSS)
}
