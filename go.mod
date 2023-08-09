module github.com/ystyle/kas

go 1.18

require (
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/asdine/storm/v3 v3.2.1
	github.com/bmaupin/go-epub v1.1.0
	github.com/gorilla/websocket v1.5.0
	github.com/labstack/echo/v4 v4.9.0
	github.com/labstack/gommon v0.3.1
	github.com/satori/go.uuid v1.2.0
	github.com/ystyle/google-analytics v0.0.0-20210425064301-a7f754dd0649
	github.com/ystyle/kaf-cli v1.3.5
	golang.org/x/net v0.7.0
	golang.org/x/text v0.7.0
)

require (
	github.com/766b/mobi v0.0.0-20200528201125-c87aa9e3c890 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/gabriel-vasile/mimetype v1.3.1 // indirect
	github.com/gofrs/uuid v3.1.0+incompatible // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/leotaku/mobi v0.0.0-20220405163106-82e29bde7964 // indirect
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	github.com/vincent-petithory/dataurl v0.0.0-20191104211930-d1553a71de50 // indirect
	go.etcd.io/bbolt v1.3.4 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
)

replace (
	golang.org/x/crypto v0.0.0-20220518034528-6f7dac969898 => github.com/golang/crypto v0.0.0-20220518034528-6f7dac969898
	golang.org/x/net v0.0.0-20220517181318-183a9ca12b87 => github.com/golang/net v0.0.0-20220517181318-183a9ca12b87
	golang.org/x/sys v0.0.0-20220519141025-dcacdad47464 => github.com/golang/sys v0.0.0-20220519141025-dcacdad47464
	golang.org/x/text v0.3.7 => github.com/golang/text v0.3.7
	golang.org/x/time v0.0.0-20220411224347-583f2d630306 => github.com/golang/time v0.0.0-20220411224347-583f2d630306
)
