package services

import (
	"bytes"
	"fmt"
	"github.com/bmaupin/go-epub"
	"github.com/labstack/gommon/log"
	"github.com/ystyle/kas/core"
	"github.com/ystyle/kas/model"
	"github.com/ystyle/kas/util/array"
	"github.com/ystyle/kas/util/character"
	"github.com/ystyle/kas/util/file"
	"github.com/ystyle/kas/util/kindlegen"
	"github.com/ystyle/kas/util/zlib"
	"io"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

const (
	htmlPStart     = `<p  style="text-indent: %dem">`
	htmlPEnd       = "</p>"
	htmlTitleStart = `<h3 style="text-align:%s">`
	htmlTitleEnd   = "</h3>"
)

func TextUpload(client *core.WsClient, message core.Message) {
	var bookinfo model.TextInfo
	err := message.JsonParse(&bookinfo)
	if err != nil {
		client.WsSend <- core.NewMessage("Error", "参数解析失败")
		return
	}
	bookinfo.SetDefault()
	reg, err := regexp.Compile(bookinfo.Match)
	if err != nil {
		errmsg := fmt.Sprintf("生成匹配规则出错: %s\n%s\n", bookinfo.Match, err.Error())
		client.WsSend <- core.NewMessage("Error", errmsg)
		return
	}

	out, err := zlib.Decode(bookinfo.Content)
	var buff bytes.Buffer
	buff.Write(out)
	var title string
	var content bytes.Buffer

	for {
		line, err := buff.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			client.WsSend <- core.NewMessage("Error", "参数解析失败")
			return
		}
		if !character.IsUtf8([]byte(line)) {
			line = character.ToUTF8(line)
		}
		line = strings.TrimSpace(line)
		// 空行直接跳过
		if len(line) == 0 {
			continue
		}
		// 处理标题
		if len(line) <= bookinfo.MaxLen && reg.MatchString(line) {
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
	err := e.Write(book.CacheMobi)
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
	err = kindlegen.Conver(book.CacheEpub, book.ID)
	if err != nil {
		client.WsSend <- core.NewMessage("Error", "生成mobi文件错误")
		return
	}
	// 复制到保存目录
	err = TextCompressZip(client, book, "mobi")
	if err != nil {
		client.WsSend <- core.NewMessage("Error", "生成epub文件错误")
		return
	}
	// 下载mobi文件
	bookDownload(client, book, "mobi")
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
