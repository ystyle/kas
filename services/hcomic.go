package services

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/ystyle/kas/core"
	"github.com/ystyle/kas/model"
	"github.com/ystyle/kas/util/config"
	"github.com/ystyle/kas/util/file"
	"github.com/ystyle/kas/util/hcomic"
	"github.com/ystyle/kas/util/kindlegen"
	"github.com/ystyle/kas/util/web"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"
)

func Submit(client *core.WsClient, message core.Message) {
	// 解析参数
	var book model.HcomicInfo
	err := message.JsonParse(&book)
	if err != nil {
		client.WsSend <- core.NewMessage("Error", "参数解析失败")
		return
	}
	id, err := hcomic.GetComicID(book.Url)
	if err != nil {
		client.WsSend <- core.NewMessage("Error", "输入的url无效")
		return
	}
	book.ID = id
	book.SetDefault()

	// zip文件存在时直接下载
	if ok, _ := file.IsExists(book.ZipFile); ok {
		client.WsSend <- core.NewMessage("info", "文件存在，从缓存读取...")
		DownloadZip(client, book.ZipFile)
		return
	}

	// 解析漫画所有的图片地址
	err = hcomic.GetAllImages(&book)
	if err != nil {
		client.WsSend <- core.NewMessage("Error", err.Error())
		return
	}
	client.WsSend <- core.NewMessage("info", fmt.Sprintf("漫画图片链接解析完成, 共%d张", len(book.Sections)))

	// 生成工作目录
	client.WsSend <- core.NewMessage("info", "创建缓存目录...")
	err1 := os.RemoveAll(book.WorkDir)
	err2 := os.MkdirAll(book.ScaledImagesDir, config.Perm)
	if err1 != nil || err2 != nil {
		client.WsSend <- core.NewMessage("Error", "服务错误: 生成缓存目录失败!")
		return
	}

	// 下载漫画图片
	client.WsSend <- core.NewMessage("info", "开始下载漫画图片...")
	var wg sync.WaitGroup
	wg.Add(len(book.Sections))
	for _, section := range book.Sections {
		go download(&wg, client, section)
	}
	wg.Wait()
	// 生成html文件
	client.WsSend <- core.NewMessage("info", "正在生成html文件...")
	err = hcomic.GenDoc(book)
	if err != nil {
		log.Error(err)
		client.WsSend <- core.NewMessage("Error", "服务错误: 生成html失败")
		return
	}
	client.WsSend <- core.NewMessage("info", "转换成Kindle mobi格式...")

	// 转换成mobi， 尝试转10次
	var hcErr error
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 10)
		hcErr = kindlegen.Conver(book.OpfFile, book.MobiName)
		if hcErr == nil {
			break
		}
	}
	if hcErr != nil {
		log.Error(err)
		client.WsSend <- core.NewMessage("Error", "服务错误: 生成mobi失败")
		return
	}
	client.WsSend <- core.NewMessage("info", "生成mobi文件成功!")
	// 压缩成zip
	err = CompressZip(client, book)
	if err != nil {
		return
	}
	// 下载zip
	DownloadZip(client, book.ZipFile)
}

func download(wg *sync.WaitGroup, client *core.WsClient, section *model.HcomicSection) {
	defer wg.Done()
	_ = web.Download(section.Url, section.ImgFile)
	client.WsSend <- core.NewMessage("info", fmt.Sprintf("下载完成: %s", section.Url))
}

func CompressZip(client *core.WsClient, book model.HcomicInfo) error {
	client.WsSend <- core.NewMessage("info", "正在压缩zip文件...")
	dir := path.Dir(book.ZipFile)
	file.CheckDir(dir)
	err := file.CompressZipToFile(book.MobiFile, book.ZipFile)
	if err != nil {
		log.Error(err)
		client.WsSend <- core.NewMessage("Error", "服务错误: 压缩mobi失败")
		return err
	}
	client.WsSend <- core.NewMessage("info", "压缩完成！")
	return err
}

func DownloadZip(client *core.WsClient, filename string) {
	buff, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err)
		client.WsSend <- core.NewMessage("Error", "服务错误: 读取zip失败")
		return
	}
	client.WsSend <- core.NewMessage("title", path.Base(filename))
	client.WsSend <- core.NewMessage("info", fmt.Sprintf("文件大小: %s, 正在下载...", file.FormatBytesLength(len(buff))))
	if web.IsMobile(client.HttpRequest.UserAgent()) {
		client.WsSend <- core.NewMessage("downloadURL", fmt.Sprintf("/download/%s", path.Base(filename)))
	} else {
		client.WsSend <- core.NewMessage("download", buff)
	}
}
