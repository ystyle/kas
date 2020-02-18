package main

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ystyle/kas/core"
	"github.com/ystyle/kas/services"
	"github.com/ystyle/kas/util/config"
	"github.com/ystyle/kas/util/file"
	"github.com/ystyle/kas/util/hcomic"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"
)

func WS(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		wm := core.GetWsManager()
		wm.Add(core.NewWsClient(ws))
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func createStoreDir() {
	if ok, _ := file.IsExists(path.Join(config.StoreDir)); !ok {
		err := os.MkdirAll(path.Join(config.StoreDir), config.Perm)
		if err != nil {
			log.Fatal("服务启动失败: 没有写入权限")
			return
		}
	}
}

func main() {
	createStoreDir()

	log.EnableColor()
	if os.Getenv("MODE") == "DEBUG" {
		log.SetLevel(log.DEBUG)
		log.Info("log level: Debug")
	}

	wm := core.GetWsManager()
	// hcomic
	wm.RegisterService("download", services.Submit)
	// text to epub / mobi
	wm.RegisterService("text:upload", services.TextUpload)
	wm.RegisterService("text:preview", services.TextPreView)
	wm.RegisterService("text:convert", services.TextConvert)
	wm.RegisterService("text:download", services.TextDownload)
	// aricle
	wm.RegisterService("article:submit", services.ArticleSubmit)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())

	box := rice.MustFindBox("public")
	assetHandler := http.FileServer(box.HTTPBox())
	e.GET("/*", echo.WrapHandler(assetHandler))
	e.GET("/asset/*", echo.WrapHandler(assetHandler))
	e.Static("/download", "storage")
	e.GET("/ws", WS)

	timer := time.NewTimer(time.Second * 5)
	go func() {
		i := 0
		for {
			select {
			case <-timer.C:
				timer.Reset(time.Second * 60)
				clients := len(wm.GetClients())
				if clients != i {
					fmt.Println("连接数为: ", clients)
				}
				i = clients
			}
		}
	}()
	if runtime.GOOS == "windows" {
		dir := path.Dir(os.Args[0])
		hcomic.Run(dir, "cmd", "/c", "start", "http://127.0.0.1:1323")
	}
	e.Logger.Fatal(e.Start(":1323"))
}
