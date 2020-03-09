package model

import (
	"github.com/asdine/storm/v3"
	"github.com/labstack/gommon/log"
	"github.com/ystyle/kas/util/config"
	"github.com/ystyle/kas/util/env"
	"github.com/ystyle/kas/util/file"
	"path"
	"sync"
)

var store *storm.DB
var once sync.Once

func DB() *storm.DB {
	once.Do(func() {
		filename := path.Join(config.StoreDir, env.GetString("DB_NAME", "kaf.db"))
		file.CheckDir(path.Dir(filename))
		db, err := storm.Open(filename, storm.BoltOptions(config.Perm, nil))
		if err != nil {
			log.Fatal(err)
		}
		store = db
	})
	return store
}
