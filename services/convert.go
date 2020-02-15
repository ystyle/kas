package services

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/ystyle/hcc/core"
	"github.com/ystyle/hcc/util/hcomic"
	"github.com/ystyle/hcc/util/web"
	"github.com/ystyle/hcc/util/zip"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"
)

type ComicInfo struct {
	Url    string
	Title  string
	Author string
}

func Submit(client *core.WsClient, message core.Message) {
	// 解析参数
	var info ComicInfo
	err := message.JsonParse(&info)
	if err != nil {
		client.WsSend <- core.NewMessage("Error", "参数解析失败")
		return
	}

	// 计算各种文件名
	comicId, _ := hcomic.GetComicID(info.Url)
	workdir := path.Join("cache", comicId)
	imagesDir := path.Join(workdir, "html", "scaled-images")
	bookName := fmt.Sprintf("%s.mobi", comicId)
	bookFile := path.Join(workdir, bookName)
	zipFile := path.Join("storage", fmt.Sprintf("%s.zip", comicId))

	// zip文件存在时直接下载
	if ok, _ := zip.IsExists(zipFile); ok {
		client.WsSend <- core.NewMessage("info", "文件存在，从缓存读取...")
		DownloadZip(client, zipFile)
		return
	}

	// 解析漫画所有的图片地址
	images, title, err := hcomic.GetAllImages(info.Url)
	if err != nil {
		client.WsSend <- core.NewMessage("Error", err.Error())
		return
	}
	client.WsSend <- core.NewMessage("info", fmt.Sprintf("漫画图片链接解析完成, 共%d张", len(images)))

	// 生成工作目录
	client.WsSend <- core.NewMessage("info", "创建缓存目录...")
	err1 := os.RemoveAll(workdir)
	err = os.MkdirAll(imagesDir, 0666)
	err2 := os.MkdirAll(path.Join("storage"), 0666)
	if err != nil || err1 != nil || err2 != nil {
		client.WsSend <- core.NewMessage("Error", "服务错误: 没有写入失败")
		return
	}

	// 下载漫画图片
	client.WsSend <- core.NewMessage("info", "开始下载漫画图片...")
	var wg sync.WaitGroup
	wg.Add(len(images))
	for _, image := range images {
		go download(&wg, client, imagesDir, image)
	}
	wg.Wait()
	// 生成html文件
	if info.Title != "" {
		title = info.Title
		client.WsSend <- core.NewMessage("title", title)
	}
	client.WsSend <- core.NewMessage("info", "正在生成html文件...")
	err = hcomic.GenDoc(title, info.Author, info.Url, images)
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
		hcErr = hcomic.ConverToMobi(path.Join(workdir, "content.opf"), bookName)
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
	err = CompressZip(client, bookFile, zipFile)
	if err != nil {
		return
	}
	// 下载zip
	DownloadZip(client, zipFile)
}

func download(wg *sync.WaitGroup, client *core.WsClient, workdir, image string) {
	defer wg.Done()
	dir := path.Join(workdir, path.Base(image))
	_ = web.Download(hcomic.GetHDImage(image), dir)
	client.WsSend <- core.NewMessage("info", fmt.Sprintf("下载完成: %s", image))
}

func CompressZip(client *core.WsClient, filename, zipFile string) error {
	client.WsSend <- core.NewMessage("info", "正在压缩zip文件...")
	buff, err := zip.CompressZip(filename)
	if err != nil {
		log.Error(err)
		client.WsSend <- core.NewMessage("Error", "服务错误: 压缩mobi失败")
		return err
	}
	client.WsSend <- core.NewMessage("info", "压缩完成！")
	return ioutil.WriteFile(zipFile, buff, 0666)
}

func DownloadZip(client *core.WsClient, filename string) {
	buff, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err)
		client.WsSend <- core.NewMessage("Error", "服务错误: 读取zip失败")
		return
	}
	client.WsSend <- core.NewMessage("title", path.Base(filename))
	client.WsSend <- core.NewMessage("info", "正在下载...")
	client.WsSend <- core.NewMessage("download", buff)
}
