package model

import (
	"github.com/asdine/storm/v3"
	"github.com/labstack/gommon/log"
	"github.com/ystyle/kas/util/config"
	"github.com/ystyle/kas/util/env"
	"path"
)

var store *storm.DB

func init() {
	filename := env.GetString("DB_NAME", "kaf.db")
	db, err := storm.Open(path.Join(config.StoreDir, filename))
	if err != nil {
		log.Fatal(err)
	}
	store = db
}

func DB() *storm.DB {
	return store
}
