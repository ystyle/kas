package model

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/ystyle/kas/util/config"
	"path"
)

type HcomicInfo struct {
	ID              string           // id
	UUID            string           // 书籍唯一码
	BookName        string           // 书名
	Author          string           // 作者
	Url             string           // 地址
	Sections        []*HcomicSection // 章节内容
	WorkDir         string           // 缓存目录
	HtmlDir         string           // html 目录
	ScaledImagesDir string           // 图片缓存目录
	OpfFile         string           // opf 文件名
	NxcFile         string           // nxc 文件名
	CoverFile       string           // 封面文件名
	MobiName        string           // mobi文件名
	MobiFile        string           // mobi文件路径
	ZipFile         string           // zip位置
}

type HcomicSection struct {
	Index      int    // 顺序
	Title      string // 标题
	Url        string // 地址
	ImgFile    string // 图片缓存路径
	HtmlFile   string // html 缓存路径
	InnerHtml  string // html 内部目录
	InnerImage string // 图片 内部目录
}

func (hcomic *HcomicInfo) SetDefault() {
	if hcomic.Author == "" {
		hcomic.Author = "KAF"
	}
	hcomic.UUID = uuid.NewV4().String()
	hcomic.WorkDir = path.Join(config.CacheDir, "hcomic", hcomic.ID)
	hcomic.HtmlDir = path.Join(hcomic.WorkDir, "html")
	hcomic.ScaledImagesDir = path.Join(hcomic.HtmlDir, "scaled-images")
	hcomic.OpfFile = path.Join(hcomic.WorkDir, "content.opf")
	hcomic.NxcFile = path.Join(hcomic.WorkDir, "toc.ncx")
	hcomic.CoverFile = path.Join(hcomic.WorkDir, "cover-image.jpg")
	hcomic.MobiName = fmt.Sprintf("%s.mobi", hcomic.ID)
	hcomic.MobiFile = path.Join(hcomic.WorkDir, hcomic.MobiName)
	hcomic.ZipFile = path.Join(config.StoreDir, "hcomic", fmt.Sprintf("%s.zip", hcomic.ID))
}

func (hcomic *HcomicInfo) AddSection(title, url string) {
	index := len(hcomic.Sections)
	imageFilename := path.Base(url)
	hcomic.Sections = append(hcomic.Sections, &HcomicSection{
		Index:      index,
		Title:      title,
		Url:        url,
		ImgFile:    path.Join(hcomic.ScaledImagesDir, imageFilename),
		HtmlFile:   path.Join(hcomic.HtmlDir, fmt.Sprintf("Page-%d.html", index)),
		InnerHtml:  fmt.Sprintf("html/Page-%d.html", index),
		InnerImage: fmt.Sprintf("scaled-images/%s", imageFilename),
	})
}
