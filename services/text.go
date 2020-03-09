package services

import (
	"bytes"
	"fmt"
	"github.com/bmaupin/go-epub"
	"github.com/labstack/gommon/log"
	"github.com/ystyle/kas/core"
	"github.com/ystyle/kas/model"
	"github.com/ystyle/kas/util/array"
	"github.com/ystyle/kas/util/file"
	"github.com/ystyle/kas/util/kindlegen"
	"github.com/ystyle/kas/util/zlib"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	htmlPStart     = `<p  style="text-indent: %dem">`
	htmlPEnd       = "</p>"
	htmlTitleStart = `<h3 style="text-align:%s">`
	htmlTitleEnd   = "</h3>"
	Tutorial       = `本书由KAF生成: <br/>
制作教程: <a href='https://ystyle.top/2019/12/31/txt-converto-epub-and-mobi/'>https://ystyle.top/2019/12/31/txt-converto-epub-and-mobi</a>
`
)

func TextUpload(client *core.WsClient, message core.Message) {
	var bookinfo model.TextInfo
	err := message.JsonParse(&bookinfo)
	if err != nil {
		client.WsSend <- core.NewMessage("Error", "参数解析失败")
		return
	}
	model.Statistics(message.DriveID)
	bookinfo.SetDefault()
	reg, err := regexp.Compile(bookinfo.Match)
	if err != nil {
		errmsg := fmt.Sprintf("生成匹配规则出错: %s\n%s\n", bookinfo.Match, err.Error())
		client.WsSend <- core.NewMessage("Error", errmsg)
		return
	}
	out, err := zlib.Decode(bookinfo.Content)
	var buff bytes.Buffer
	encodig, encodename, _ := charset.DetermineEncoding(out[:1024], "text/plain")
	if encodename != "utf-8" {
		bs, _, _ := transform.Bytes(encodig.NewDecoder(), out)
		buff.Write(bs)
	} else {
		buff.Write(out)
	}

	var title string
	var content bytes.Buffer

	for {
		line, err := buff.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				bookinfo.AddSection(title, content.String())
				break
			}
			client.WsSend <- core.NewMessage("Error", "参数解析失败")
			return
		}
		line = strings.TrimSpace(line)
		// 空行直接跳过
		if len(line) == 0 {
			continue
		}
		// 处理标题
		if utf8.RuneCountInString(line) <= bookinfo.MaxLen && reg.MatchString(line) {
			if title == "" {
				title = "说明"
			}
			bookinfo.AddSection(title, content.String())
			title = line
			content.Reset()
			content.WriteString(fmt.Sprintf(htmlTitleStart, bookinfo.Align))
			content.WriteString(title)
			content.WriteString(htmlTitleEnd)
			continue
		}
		if strings.HasSuffix(line, "==") ||
			strings.HasSuffix(line, "**") ||
			strings.HasSuffix(line, "--") ||
			strings.HasSuffix(line, "//") {
			content.WriteString(line)
			continue
		}
		content.WriteString(fmt.Sprintf(htmlPStart, bookinfo.Indent))
		content.WriteString(line)
		content.WriteString(htmlPEnd)
	}
	bookinfo.AddSection("制作说明", Tutorial)
	client.Caches[bookinfo.ID] = bookinfo
	client.WsSend <- core.NewMessage("info", "解析完成")
	client.WsSend <- core.NewMessage("text:uploaded", bookinfo.ID)
}

func TextPreView(client *core.WsClient, message core.Message) {
	id := message.GetString()
	if _, ok := client.Caches[id]; !ok {
		client.WsSend <- core.NewMessage("Error", "书籍信息不存在")
		return
	}
	book := client.Caches[id].(model.TextInfo)
	var titleList []string
	for i, section := range book.Sections {
		if i > 30 {
			break
		}
		titleList = append(titleList, section.Title)
	}
	client.WsSend <- core.NewMessage("text:titleList", titleList)
}

func TextConvert(client *core.WsClient, message core.Message) {
	id := message.GetString()
	if _, ok := client.Caches[id]; !ok {
		client.WsSend <- core.NewMessage("Error", "书籍信息不存在")
		return
	}
	book := client.Caches[id].(model.TextInfo)
	// 生成epub
	client.WsSend <- core.NewMessage("info", "正在生成生成epub文件...")
	e := epub.NewEpub(book.BookName)
	e.SetAuthor(book.Author)
	for _, section := range book.Sections {
		e.AddSection(section.Content, section.Title, "", "")
	}
	file.CheckDir(path.Dir(book.CacheEpub))
	err := e.Write(book.CacheEpub)
	if err != nil {
		client.WsSend <- core.NewMessage("Error", "生成epub文件错误")
		return
	}
	// 复制到保存目录
	err = TextCompressZip(client, book, "epub")
	if err != nil {
		client.WsSend <- core.NewMessage("Error", "生成epub文件错误")
		return
	}
	// 下载epub文件
	if array.IncludesString(book.Format, "epub") {
		bookDownload(client, book, "epub")
	}
	if !array.IncludesString(book.Format, "mobi") {
		return
	}
	// 转换mobi文件
	client.WsSend <- core.NewMessage("info", "正在生成生成mobi文件...")
	err = kindlegen.Conver(book.CacheEpub, path.Base(book.CacheMobi), book.OnlyKF8 == "1")
	if err != nil {
		client.WsSend <- core.NewMessage("Error", "生成mobi文件错误")
		return
	}
	// 复制到保存目录
	err = TextCompressZip(client, book, "mobi")
	if err != nil {
		client.WsSend <- core.NewMessage("Error", "生成mobi文件错误")
		return
	}
	// 下载mobi文件
	bookDownload(client, book, "mobi")
	os.Remove(book.CacheEpub)
	os.Remove(book.CacheMobi)
	os.Remove(book.StoreEpub)
	os.Remove(book.StoreMobi)
}

func TextCompressZip(client *core.WsClient, book model.TextInfo, format string) error {
	zipFile := book.StoreEpub
	ebookFile := book.CacheEpub
	if format == "mobi" {
		zipFile = book.StoreMobi
		ebookFile = book.CacheMobi
	}
	client.WsSend <- core.NewMessage("info", "正在压缩zip文件...")
	dir := path.Dir(book.StoreEpub)
	file.CheckDir(dir)
	err := file.CompressZipToFile(ebookFile, zipFile)
	if err != nil {
		log.Error(err)
		client.WsSend <- core.NewMessage("Error", "服务错误: 压缩文件失败")
		return err
	}
	client.WsSend <- core.NewMessage("info", "压缩完成！")
	return err
}

func TextDownload(client *core.WsClient, message core.Message) {
	id := message.GetString()
	if _, ok := client.Caches[id]; !ok {
		client.WsSend <- core.NewMessage("Error", "书籍信息不存在")
		return
	}
	book := client.Caches[id].(model.TextInfo)
	// 下载epub文件
	if array.IncludesString(book.Format, "epub") {
		bookDownload(client, book, "epub")
	}
	// 下载epub文件
	if array.IncludesString(book.Format, "mobi") {
		bookDownload(client, book, "mobi")
	}
}

func bookDownload(client *core.WsClient, book model.TextInfo, format string) {
	filename := book.StoreMobi
	if format == "epub" {
		filename = book.StoreEpub
	}
	buff, err := ioutil.ReadFile(filename)
	readErr := fmt.Sprintf("读取文件失败: %s", path.Base(filename))
	if err != nil {
		client.WsSend <- core.NewMessage("Error", readErr)
		return
	}
	client.WsSend <- core.NewMessage("info", fmt.Sprintf("正在下载: %s， 文件大小: %s", path.Base(filename), file.FormatBytesLength(len(buff))))
	client.WsSend <- core.NewMessage("text:download", buff)
}
