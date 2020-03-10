package main

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ystyle/kas/core"
	"github.com/ystyle/kas/model"
	"github.com/ystyle/kas/services"
	"github.com/ystyle/kas/util/config"
	"github.com/ystyle/kas/util/file"
	"github.com/ystyle/kas/util/hcomic"
	"github.com/ystyle/kas/util/site"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"
)

var upgrader = websocket.Upgrader{}

func WS(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	wm := core.GetWsManager()
	wm.Add(core.NewWsClient(ws, c.Request()))
	return nil
}

func init() {
	createStoreDir()
}

func createStoreDir() {
	file.CheckDir(config.StoreDir)
}

func PrintStatistics() {
	wm := core.GetWsManager()
	timer := time.NewTimer(time.Second * 5)
	i := 0
	for {
		select {
		case <-timer.C:
			timer.Reset(time.Second * 5)
			clients := len(wm.GetClients())
			var drives []model.Drive
			err := model.DB().All(&drives)
			if err != nil {
				log.Error(err)
				continue
			}
			var count uint
			for _, drive := range drives {
				count += drive.Count
			}
			// 连接有变动时就打印
			if clients != i && err == nil {
				fmt.Printf("注册设备: %d, 总计转换次数为: %d 当前连接数为: %d\n", len(drives), count, clients)
			}
			i = clients
		}
	}
}

func main() {
	log.EnableColor()
	if os.Getenv("MODE") == "DEBUG" {
		log.SetLevel(log.DEBUG)
		log.Info("log level: Debug")
	}
	wm := core.GetWsManager()
	// hcomic
	wm.RegisterService("hcomic:submit", services.Submit)
	// text to epub / mobi
	wm.RegisterService("text:upload", services.TextUpload)
	wm.RegisterService("text:preview", services.TextPreView)
	wm.RegisterService("text:convert", services.TextConvert)
	wm.RegisterService("text:download", services.TextDownload)
	// aricle
	wm.RegisterService("article:submit", services.ArticleSubmit)
	// ping
	wm.RegisterService("ping", services.Ping)
	// 注册设备
	wm.RegisterService("regsiter", services.Register)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())

	box := rice.MustFindBox("public")
	assetHandler := http.FileServer(box.HTTPBox())
	site.Init(box)

	e.GET("/*", echo.WrapHandler(assetHandler))
	e.GET("/asset/*", echo.WrapHandler(assetHandler))
	e.Static("/download", "storage")
	e.GET("/ws", WS)
	e.GET("/ws#", WS)

	// 打印服务器负载
	go PrintStatistics()

	if runtime.GOOS == "windows" {
		dir := path.Dir(os.Args[0])
		hcomic.Run(dir, "cmd", "/c", "start", "http://127.0.0.1:1323")
	}
	e.Logger.Fatal(e.Start(":1323"))
}
