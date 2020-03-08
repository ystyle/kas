### [Hcomic Creator](https://kas.ystyle.top)

漫画生成器. 输入地址直接生成kindle的mobi漫画格式

[![](public/asset/screenshot.png)](https://kas.ystyle.top)

### 使用方法
- [下载最新版本的kas.zip](https://github.com/ystyle/kas/releases/latest)
  - linux / osx 系统下载后放添加到PATH， 并自行安装kindlegen
- 解压, 双击`kas.exe`运行, 会自动打开浏览器, 手动打开是: [http://127.0.0.1:1323/](http://127.0.0.1:1323/)

### docker 方式
```shell script
docker run -d --name hcc \
 --restart always \
 -p 1323:1323 \
 -v /mnt/kas/storage:/app/storage \
  ystyle/kas
```

### KAF自定义服务器地址
- 默认服务器为: `ws://140.143.205.67:1323/ws`
  - 如果在自己服务器启动则把ip改为自己服务器ip
  - 如果在自己电脑启动则填自己内网地址
    - windows的话在连接的wifi上，点属性，查看ipv4地址