package web

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/gommon/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var client = &http.Client{}

func init() {
	proxy := os.Getenv("PROXY")
	if proxy != "" {
		if proxyUrl, err := url.Parse(proxy); err == nil {
			log.Info("设置代理服务器: ", proxy)
			client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
		}
	}
}

func getResponse(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Arch Linux kernel 4.6.5) AppleWebKit/537.36 (KHTML, like Gecko) Maxthon/4.0 Chrome/39.0.2146.0 Safari/537.36")
	req.Header.Set("Set-Cookie", "r18=ok")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("DNT", "1")
	return client.Do(req)
}

func GetHtmlNode(url string) (*goquery.Document, error) {
	res, err := getResponse(url)
	if err != nil {
		return nil, fmt.Errorf("目标网站无法连接: %w", err)
	}
	return goquery.NewDocumentFromReader(res.Body)
}

func GetContent(url string) (string, error) {
	res, err := getResponse(url)
	if err != nil {
		return "", fmt.Errorf("目标网站无法连接: %w", err)
	}
	buff, err := ioutil.ReadAll(res.Body)
	return string(buff), err
}

func Download(url string, dir string) error {
	res, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("目标网站无法连接: %w", err)
	}
	f, err := os.Create(dir)
	if err != nil {
		return fmt.Errorf("服务器错误: 无法创建文件")
	}
	_, _ = io.Copy(f, res.Body)
	return nil
}

func IsMobile(useragent string) bool {
	agents := []string{"Android", "iPhone", "SymbianOS", "Windows Phone", "iPad", "iPod"}
	for _, agent := range agents {
		if strings.Contains(useragent, agent) {
			return true
		}
	}
	return false

}
